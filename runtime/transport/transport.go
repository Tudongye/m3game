package transport

import (
	"context"
	"log"
	"m3game/broker"
	"net"
	"regexp"

	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/proto"
)

var (
	_instance      *Transport
	_cfg           TransportCfg
	_tcpAddr       *net.TCPAddr
	regexHealth, _ = regexp.Compile("^Health/([^/]*)$")
	healthmethod   = "/grpc.health.v1.Health/Check"
)

const (
	_cfgkey = "Runtime.Transport"
)

type RegistServerFunc func(*Transport) error

type RuntimeReciver interface {
	RecvInterFunc(trecv *Reciver) (resp interface{}, err error)
	HealthCheck(idstr string) bool
}

type TransportCfg struct {
	Addr         string `mapstructure:"Addr"`
	BroadTimeOut int    `mapstructure:"BroadTimeOut"`
}

func (t *TransportCfg) CheckVaild() error {
	if _, _, err := net.SplitHostPort(t.Addr); err != nil {
		return err
	}
	return nil
}

type Transport struct {
	gser      *grpc.Server
	brokerser *BrokerSer
	cancel    context.CancelFunc
	runtime   RuntimeReciver
}

func (t *Transport) GrpcSer() *grpc.Server {
	return t.gser
}

func (t *Transport) start(ctx context.Context, listener *net.TCPListener) {
	go func() {
		select {
		case <-ctx.Done():
			Stop()
			listener.Close()
		}
	}()
	grpc_health_v1.RegisterHealthServer(t.gser, t)
	t.gser.Serve(listener)
}

func (t *Transport) recvInterFunc(rec *Reciver) (interface{}, error) {
	if rec.Info().FullMethod == healthmethod {
		return rec.HandleMsg(rec.Ctx())
	}
	return _instance.runtime.RecvInterFunc(rec)
}

func (t *Transport) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	if !regexHealth.MatchString(req.Service) {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}
	groups := regexHealth.FindStringSubmatch(req.Service)
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

func Init(c map[string]interface{}, runtime RuntimeReciver) error {
	if _instance != nil {
		return nil
	}
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return err
	}
	if err := _cfg.CheckVaild(); err != nil {
		return err
	}
	_instance = &Transport{
		runtime: runtime,
	}
	_instance.gser = grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc.UnaryServerInterceptor(RecvInteror()),
		),
	)
	_instance.brokerser = CreateBrokerSer(grpc.UnaryServerInterceptor(RecvInteror()))
	return nil
}

func Start() error {
	log.Printf("Transport.TcpAddr %s\n", _cfg.Addr)
	var err error
	_tcpAddr, err = net.ResolveTCPAddr("tcp", _cfg.Addr)
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", _tcpAddr)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	_instance.cancel = cancel
	go _instance.start(ctx, listener)
	return nil
}

func ShutDown() error {
	if _instance != nil {
		_instance.cancel()
	}
	return nil
}

func Stop() error {
	_instance.gser.Stop()
	return nil
}

func SendInterFunc(s *Sender) error {
	return s.SendMsg()
}

func RecvInteror() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if rec, err := CreateReciver(ctx, req, info, handler); err != nil {
			return handler(ctx, req)
		} else {
			return _instance.recvInterFunc(rec)
		}
	}
}

func SendInteror(f func(*Sender) error) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		s := CreateSender(ctx, method, req, resp, cc, invoker, opts)
		return f(s)
	}
}

func Addr() string {
	return _cfg.Addr
}

func TCPAddr() *net.TCPAddr {
	return _tcpAddr
}

func RegistServer(f func(grpc.ServiceRegistrar) error) error {
	if err := f(_instance.gser); err != nil {
		return err
	}
	if err := f(_instance.brokerser); err != nil {
		return err
	}
	return nil
}

func RegisterBroker(broker broker.Broker) {
	_instance.brokerser.registerBroker(broker)
}

func SendToBroker(sender *Sender, topic string) error {
	return _instance.brokerser.SendToBroker(topic, sender.Method(), sender.req.(proto.Message))
}
