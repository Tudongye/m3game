package transport

import (
	"context"
	"fmt"
	"m3game/broker"
	"m3game/util/log"
	"net"
	"regexp"
	"sync"

	"github.com/pkg/errors"

	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/proto"
)

var (
	_instance     *Transport
	_cfg          TransportCfg
	_tcpAddr      *net.TCPAddr
	_regexHealth  *regexp.Regexp
	_healthmethod = "/grpc.health.v1.Health/Check"
)

var _err_msgisnotm3pkg = errors.New("_err_msgisnotm3pkg")

type RuntimeReciver interface {
	RecvInterFunc(trecv *Reciver) (resp interface{}, err error)
	HealthCheck(idstr string) bool
}

func init() {
	var err error
	if _regexHealth, err = regexp.Compile("^Health/([^/]*)$"); err != nil {
		panic(fmt.Sprintf("Compile regexHealth err %v", err))
	}
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
	_instance.gser = grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc.UnaryServerInterceptor(recvInteror()),
		),
	)
	_instance.brokerser = newBrokerSer(grpc.UnaryServerInterceptor(recvInteror()))
	return nil
}

func Start(wg *sync.WaitGroup) error {
	log.Fatal("Transport.TcpAddr %s", _cfg.Addr)
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
	}()
	return nil
}

func ShutDown() {
	if _instance != nil {
		_instance.cancel()
	}
}

func SendInterFunc(s *Sender) error {
	return s.sendMsg()
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

func (t *Transport) recvInterFunc(rec *Reciver) (interface{}, error) {
	if rec.Info().FullMethod == _healthmethod {
		return rec.HandleMsg(rec.Ctx())
	}
	return _instance.runtime.RecvInterFunc(rec)
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

func recvInteror() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if rec, err := newReciver(ctx, req, info, handler); err != nil {
			return handler(ctx, req)
		} else {
			return _instance.recvInterFunc(rec)
		}
	}
}

func sendToBrokerSer(sender *Sender, topic string) error {
	return _instance.brokerser.send(topic, sender.Method(), sender.req.(proto.Message))
}
