package transport

import (
	"context"
	"fmt"
	"m3game/plugins/broker"

	"m3game/meta"
	"m3game/plugins/log"
	"net"
	"regexp"
	"sync"

	"github.com/pkg/errors"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

var (
	_instance     *Transport
	_cfg          TransportCfg
	_tcpAddr      *net.TCPAddr
	_regexHealth  *regexp.Regexp
	_healthmethod = "/grpc.health.v1.Health/Check"
)
var (
	_serverInterceptors []grpc.UnaryServerInterceptor
	_clientInterceptors []grpc.UnaryClientInterceptor
)

var _err_msgisnotm3pkg = errors.New("_err_msgisnotm3pkg")

type RuntimeReciver interface {
	ServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error)
	HealthCheck(idstr string) bool
}

func init() {
	var err error
	if _regexHealth, err = regexp.Compile("^Health/([^/]*)$"); err != nil {
		panic(fmt.Sprintf("Compile regexHealth err %v", err))
	}
	RegisterServerInterceptor(ServerInterceptor())
}

func Init(c map[string]interface{}, runtime RuntimeReciver) error {
	if _instance != nil {
		return nil
	}
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return errors.Wrap(err, "decode cfg")
	}
	if err := _cfg.checkvaild(); err != nil {
		return err
	}
	_instance = &Transport{
		runtime: runtime,
	}
	return nil
}

func CreateSer() error {
	_instance.gser = grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(_serverInterceptors...),
		),
	)
	_instance.brokerser = newBrokerSer(
		grpc_middleware.ChainUnaryServer(_serverInterceptors...),
	)
	b := broker.Get()
	if b == nil {
		return errors.New("Broker-Plugin not find")
	}
	if err := _instance.brokerser.registerBroker(b); err != nil {
		return err
	}
	return nil
}

func RegisterClientInterceptor(f grpc.UnaryClientInterceptor) {
	_clientInterceptors = append(_clientInterceptors, f)
}

func ClientInterceptors() []grpc.UnaryClientInterceptor {
	return _clientInterceptors
}

func RegisterServerInterceptor(f grpc.UnaryServerInterceptor) {
	_serverInterceptors = append(_serverInterceptors, f)
}

func Start(wg *sync.WaitGroup) error {
	log.Info("Transport Listen %s", _cfg.Addr)
	var err error
	_tcpAddr, err = net.ResolveTCPAddr("tcp", _cfg.Addr)
	if err != nil {
		return errors.Wrap(err, "transport")
	}
	listener, err := net.ListenTCP("tcp", _tcpAddr)
	if err != nil {
		return errors.Wrap(err, "transport.ListenTCP")
	}
	ctx, cancel := context.WithCancel(context.Background())
	_instance.cancel = cancel
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer listener.Close()
		_instance.start(ctx, listener)
		log.Info("Transport.Stoped...")
	}()
	return nil
}

func ShutDown() {
	if _instance != nil {
		log.Info("Transport.Stoping...")
		_instance.cancel()
	}
}

func Reload(c map[string]interface{}) error {
	return nil
}

func ClientInterceptor(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		if vlist, ok := md[string(meta.M3RouteType)]; ok && len(vlist) > 0 {
			if meta.RouteType(vlist[0]) == meta.RouteTypeBroad ||
				meta.RouteType(vlist[0]) == meta.RouteTypeMulti {
				if vlist, ok := md[string(meta.M3RouteTopic)]; ok && len(vlist) > 0 {
					return sendToBrokerSer(ctx, vlist[0], method, req, opts...)
				} else {
					return errors.New("RouteTypeBroad & RouteTypeMulti not find Topic")
				}

			}
		}
	}
	return invoker(ctx, method, req, resp, cc, opts...)
}

func Addr() string {
	return _cfg.Addr
}

func RegisterServer(f func(grpc.ServiceRegistrar) error) error {
	if err := f(_instance.gser); err != nil {
		return errors.Wrap(err, "gser.register")
	}
	if err := f(_instance.brokerser); err != nil {
		return errors.Wrap(err, "brokerser.register")
	}
	return nil
}

func RegisterBroker(broker broker.Broker) error {
	return _instance.brokerser.registerBroker(broker)
}

type TransportCfg struct {
	Addr         string `mapstructure:"Addr"`
	BroadTimeOut int    `mapstructure:"BroadTimeOut"`
}

func (t *TransportCfg) checkvaild() error {
	if _, _, err := net.SplitHostPort(t.Addr); err != nil {
		return errors.Wrap(err, "TransportCfg.Addr")
	}
	return nil
}

type Transport struct {
	gser      *grpc.Server
	brokerser *brokerSer
	cancel    context.CancelFunc
	runtime   RuntimeReciver
}

func (t *Transport) start(ctx context.Context, listener *net.TCPListener) error {
	go func() {
		select {
		case <-ctx.Done():
			t.gser.Stop()
		}
	}()
	grpc_health_v1.RegisterHealthServer(t.gser, t)
	return t.gser.Serve(listener)
}

func (t *Transport) serverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == _healthmethod {
		return handler(ctx, req)
	}
	return _instance.runtime.ServerInterceptor(ctx, req, info, handler)
}

func (t *Transport) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	if !_regexHealth.MatchString(req.Service) {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}
	groups := _regexHealth.FindStringSubmatch(req.Service)
	idstr := groups[1]
	if t.runtime.HealthCheck(idstr) {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_SERVING,
		}, nil
	} else {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}
}

func (r *Transport) Watch(req *grpc_health_v1.HealthCheckRequest, w grpc_health_v1.Health_WatchServer) error {
	return nil
}

func ServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		return _instance.serverInterceptor(ctx, req, info, handler)
	}
}

func sendToBrokerSer(ctx context.Context, topic string, method string, req interface{}, opts ...grpc.CallOption) error {
	return _instance.brokerser.send(ctx, topic, method, req, opts...)
}
