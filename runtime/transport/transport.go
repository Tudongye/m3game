package transport

import (
	"context"
	"m3game/meta"
	"m3game/plugins/broker"

	"m3game/plugins/log"
	"net"
	"regexp"

	"github.com/pkg/errors"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

const (
	_grpcHealthCheckMethod = "/grpc.health.v1.Health/Check"
	_healthPathPattern     = "^Health/([^/]*)$"
)

var (
	_regexHealth *regexp.Regexp
)

type RuntimeReciver interface {
	ServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error)
	HealthCheck(idstr string) bool
}

func init() {
	var err error
	if _regexHealth, err = regexp.Compile(_healthPathPattern); err != nil {
		log.Error("Compile Health Fail %s", err.Error())
		_regexHealth = nil
	}
}

func New(c map[string]interface{}, runtimeReciver RuntimeReciver) (*Transport, error) {
	if _regexHealth == nil {
		return nil, errors.New("_regexHealth is nil")
	}
	var cfg TransportCfg
	if err := mapstructure.Decode(c, &cfg); err != nil {
		return nil, errors.Wrap(err, "decode cfg")
	}
	if err := cfg.checkValid(); err != nil {
		return nil, err
	}
	transport := &Transport{
		cfg:            cfg,
		runtimeReciver: runtimeReciver,
	}
	transport.RegisterServerInterceptor(transport.serverInterceptor)
	return transport, nil
}

type TransportCfg struct {
	Addr             string `mapstructure:"Addr"`
	BroadcastTimeout int    `mapstructure:"BroadcastTimeout"`
}

func (t *TransportCfg) checkValid() error {
	if _, _, err := net.SplitHostPort(t.Addr); err != nil {
		return errors.Wrap(err, "TransportCfg.Addr")
	}
	return nil
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
		return errors.Wrap(err, "transport")
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return errors.Wrap(err, "transport.ListenTCP")
	}
	defer listener.Close()
	grpc_health_v1.RegisterHealthServer(t.server, t)
	return t.server.Serve(listener)
}

func (t *Transport) serverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == _grpcHealthCheckMethod {
		return handler(ctx, req)
	}
	return t.runtimeReciver.ServerInterceptor(ctx, req, info, handler)
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
	t.brokerser = newBrokerSer(
		t.cfg,
		grpc_middleware.ChainUnaryServer(t.serverInterceptors...),
	)
	brokerins := broker.Get()
	if brokerins == nil {
		return errors.New("Broker-Plugin not find")
	}
	if err := t.brokerser.registerBroker(brokerins); err != nil {
		return err
	}
	return nil
}

func (t *Transport) RegisterServer(f func(grpc.ServiceRegistrar) error) error {
	if err := f(t.server); err != nil {
		return errors.Wrap(err, "server.register")
	}
	if err := f(t.brokerser); err != nil {
		return errors.Wrap(err, "brokerser.register")
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
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		if vlist, ok := md[string(meta.M3RouteType)]; ok && len(vlist) > 0 {
			if meta.RouteType(vlist[0]) == meta.RouteTypeBroad ||
				meta.RouteType(vlist[0]) == meta.RouteTypeMulti {
				if vlist, ok := md[string(meta.M3RouteTopic)]; ok && len(vlist) > 0 {
					return t.brokerser.send(ctx, vlist[0], method, req, opts...)
				} else {
					return errors.New("RouteTypeBroad & RouteTypeMulti not find Topic")
				}

			}
		}
	}
	return invoker(ctx, method, req, resp, cc, opts...)
}

func (t *Transport) Addr() string {
	return t.cfg.Addr
}
