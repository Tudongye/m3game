package transport

import (
	"context"
	"fmt"
	"m3game/meta"
	"m3game/meta/errs"
	"m3game/meta/monitor"
	"m3game/plugins/broker"
	"m3game/plugins/metric"

	"m3game/plugins/log"
	"net"
	"regexp"

	validator "github.com/go-playground/validator/v10"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

const (
	_grpcHealthCheckMethod = "/grpc.health.v1.Health/Check" // Grpc健康检测默认health
	_healthPathPattern     = "^Health/([^/]*)$"             // 健康检测参数
)

var (
	_regexHealth *regexp.Regexp
	_transport   *Transport
)

// Runtime接收器
type RuntimeReciver interface {
	ServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error)
	HealthCheck(idstr string) bool
}

func init() {
	var err error
	if _regexHealth, err = regexp.Compile(_healthPathPattern); err != nil {
		log.Fatal("Compile Health Fail %s", err.Error())
	}
}

type TransportCfg struct {
	Host             string `mapstructure:"Host" validate:"required"`         // 监听地址
	Port             int    `mapstructure:"Port" validate:"gt=0"`             // 监听地址
	BroadcastTimeout int    `mapstructure:"BroadcastTimeout" validate:"gt=0"` // BrokerSer的Handler超时时间
	CloseBroker      int    `mapstructure:"CloseBroker"`                      // 是否关闭Broker
	Addr             string
}

func New(c map[string]interface{}, runtimeReciver RuntimeReciver) (*Transport, error) {
	if _transport != nil {
		return _transport, nil
	}
	if _regexHealth == nil {
		return nil, errs.TransportInitFail.New("_regexHealth is nil")
	}
	var cfg TransportCfg
	if err := mapstructure.Decode(c, &cfg); err != nil {
		return nil, errs.TransportInitFail.Wrap(err, "decode cfg")
	}
	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, err
	}
	cfg.Addr = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	_transport = &Transport{
		cfg:            cfg,
		runtimeReciver: runtimeReciver,
	}
	// 注册M3RPC的ServerInterceptor
	_transport.RegisterServerInterceptor(_transport.serverInterceptor)
	return _transport, nil
}

type Transport struct {
	cfg                TransportCfg
	server             *grpc.Server
	brokerser          *brokerSer
	runtimeReciver     RuntimeReciver
	serverInterceptors []grpc.UnaryServerInterceptor
}

func (t *Transport) Start(ctx context.Context) error {
	go func() {
		select {
		case <-ctx.Done():
			t.server.Stop()
		}
	}()
	log.Info("Transport Listen %s", t.cfg.Addr)
	var err error
	tcpAddr, err := net.ResolveTCPAddr("tcp", t.cfg.Addr)
	if err != nil {
		return errs.TransportInitFail.Wrap(err, "transport")
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return errs.TransportInitFail.Wrap(err, "transport.ListenTCP")
	}
	defer listener.Close()
	grpc_health_v1.RegisterHealthServer(t.server, t)
	return t.server.Serve(listener)
}

func (t *Transport) serverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == _grpcHealthCheckMethod {
		metric.Counter(monitor.HandleHealthRPCTotal).Inc()
		return handler(ctx, req)
	}
	metric.Counter(monitor.HandleRPCTotal).Inc()
	rsp, err := t.runtimeReciver.ServerInterceptor(ctx, req, info, handler)
	if err != nil {
		metric.Counter(monitor.HandleRPCFailTotal).Inc()
	}
	return rsp, err
}

func (t *Transport) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	if !_regexHealth.MatchString(req.Service) {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}
	groups := _regexHealth.FindStringSubmatch(req.Service)
	idstr := groups[1]
	if t.runtimeReciver.HealthCheck(idstr) {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_SERVING,
		}, nil
	} else {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}
}

func (t *Transport) Watch(req *grpc_health_v1.HealthCheckRequest, w grpc_health_v1.Health_WatchServer) error {
	return nil
}

func (t *Transport) Prepare(ctx context.Context) error {
	t.server = grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(t.serverInterceptors...),
		),
	)
	if t.cfg.CloseBroker == 0 {
		t.brokerser = newBrokerSer(
			t.cfg,
			grpc_middleware.ChainUnaryServer(t.serverInterceptors...),
		)
		brokerins := broker.Instance()
		if brokerins == nil {
			return errs.TransportInitFail.New("Broker-Plugin not find")
		}
		if err := t.brokerser.setBroker(brokerins); err != nil {
			return err
		}
	}
	return nil
}

// 注册Grpc的Handler
func (t *Transport) RegisterServer(f func(grpc.ServiceRegistrar) error) error {
	if err := f(t.server); err != nil {
		return errs.TransportRegisterSerFail.Wrap(err, "server.register")
	}
	if t.cfg.CloseBroker == 0 {
		if err := f(t.brokerser); err != nil {
			return errs.BrokerSerRegisterSerFail.Wrap(err, "brokerser.register")
		}
	}
	return nil
}

func (t *Transport) Reload(c map[string]interface{}) error {
	return nil
}

func (t *Transport) RegisterServerInterceptor(f grpc.UnaryServerInterceptor) {
	t.serverInterceptors = append(t.serverInterceptors, f)
}

func (t *Transport) ClientInterceptor(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	metric.Counter(monitor.CallPRCTotal).Inc()
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		if vlist, ok := md[string(meta.M3RouteType)]; ok && len(vlist) > 0 {
			if meta.RouteType(vlist[0]) == meta.RouteTypeBroad ||
				meta.RouteType(vlist[0]) == meta.RouteTypeMulti {
				if vlist, ok := md[string(meta.M3RouteTopic)]; ok && len(vlist) > 0 {
					if t.cfg.CloseBroker == 0 {
						err := t.brokerser.send(ctx, vlist[0], method, req, opts...)
						if err != nil {
							metric.Counter(monitor.CallRPCFailTotal).Inc()
						}
						return err
					} else {
						return errs.BrokerSerClose.New("")
					}
				} else {
					metric.Counter(monitor.CallRPCFailTotal).Inc()
					return errs.TransportCliCantFindTopic.New("RouteTypeBroad & RouteTypeMulti not find Topic")
				}

			}
		}
	}
	err := invoker(ctx, method, req, resp, cc, opts...)
	if err != nil {
		metric.Counter(monitor.CallRPCFailTotal).Inc()
	}
	return err
}

func (t *Transport) Addr() string {
	return t.cfg.Addr
}

func (t *Transport) Host() string {
	return t.cfg.Host
}

func (t *Transport) Port() int {
	return t.cfg.Port
}
