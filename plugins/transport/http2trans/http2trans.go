package http2trans

import (
	"context"
	"fmt"
	"m3game/meta"
	"m3game/meta/errs"
	"m3game/meta/monitor"
	"m3game/plugins/log"
	"m3game/plugins/metric"
	"m3game/plugins/transport"
	"m3game/runtime/mesh"
	"m3game/runtime/plugin"
	"net"
	"time"

	"github.com/go-playground/validator/v10"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	_           plugin.Factory   = (*Factory)(nil)
	_           plugin.PluginIns = (*Http2Trans)(nil)
	_http2trans *Http2Trans
	_factory    = &Factory{}
)

const (
	_name = "trans_http2"
)

func init() {
	plugin.RegisterFactory(_factory)
}

type Http2TransCfg struct {
	Host string `mapstructure:"Host" validate:"required"` // 监听地址
	Port int    `mapstructure:"Port" validate:"gt=0"`     // 监听端口
}

type Factory struct {
}

func (f *Factory) Type() plugin.Type {
	return plugin.Trans
}

func (f *Factory) Name() string {
	return _name
}

func (f *Factory) Setup(ctx context.Context, c map[string]interface{}) (plugin.PluginIns, error) {
	if _http2trans != nil {
		return _http2trans, nil
	}
	var cfg Http2TransCfg
	if err := mapstructure.Decode(c, &cfg); err != nil {
		return nil, errs.TransportSetupFail.Wrap(err, "RedisDB Decode Cfg")
	}
	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, errs.TransportSetupFail.Wrap(err, "")
	}
	_http2trans = &Http2Trans{
		cfg:  cfg,
		addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}
	_http2trans.RegisterServerInterceptor(_http2trans.serverInterceptor)
	_http2trans.RegisterClientInterceptor(_http2trans.clientInterceptor)
	if _, err := transport.New(_http2trans); err != nil {
		return nil, err
	}
	return _http2trans, nil
}

func (f *Factory) Destroy(plugin.PluginIns) error {
	return nil
}

func (f *Factory) Reload(plugin.PluginIns, map[string]interface{}) error {
	return nil
}

func (f *Factory) CanUnload(plugin.PluginIns) bool {
	return false
}

type Http2Trans struct {
	cfg                Http2TransCfg
	addr               string
	gser               *grpc.Server
	serverInterceptors []grpc.UnaryServerInterceptor
	clientInterceptors []grpc.UnaryClientInterceptor
}

func (t *Http2Trans) Factory() plugin.Factory {
	return _factory
}

func (t *Http2Trans) Host() string {
	return t.cfg.Host
}

func (t *Http2Trans) Port() int {
	return t.cfg.Port
}

func (t *Http2Trans) Prepare(ctx context.Context) error {
	t.gser = grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(t.serverInterceptors...),
		),
	)
	return nil
}
func (t *Http2Trans) Start(ctx context.Context) error {
	go func() {
		select {
		case <-ctx.Done():
			t.gser.Stop()
		}
	}()
	log.Info("Transport Listen %s", t.addr)
	var err error
	tcpAddr, err := net.ResolveTCPAddr("tcp", t.addr)
	if err != nil {
		return errs.TransportInitFail.Wrap(err, "transport")
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return errs.TransportInitFail.Wrap(err, "transport.ListenTCP")
	}
	defer listener.Close()
	return t.gser.Serve(listener)
}
func (t *Http2Trans) RegisterSer(f func(grpc.ServiceRegistrar) error) error {

	if err := f(t.gser); err != nil {
		return errs.TransportRegisterSerFail.Wrap(err, "server.register")
	}
	return nil
}

func (t *Http2Trans) RegisterServerInterceptor(f grpc.UnaryServerInterceptor) {
	t.serverInterceptors = append(t.serverInterceptors, f)

}

func (t *Http2Trans) RegisterClientInterceptor(f grpc.UnaryClientInterceptor) {
	t.clientInterceptors = append([]grpc.UnaryClientInterceptor{f}, t.clientInterceptors...)
}

func (t *Http2Trans) ClientInterceptors() []grpc.UnaryClientInterceptor {
	return t.clientInterceptors
}

func (t *Http2Trans) ClientConn(target string, opts ...grpc.DialOption) (grpc.ClientConnInterface, error) {
	ttarget := fmt.Sprintf("router://%s", target)
	opts = append(opts,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"Balance_m3g"}`),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(t.ClientInterceptors()...)),
		grpc.WithTimeout(time.Second*10),
	)
	if conn, err := grpc.Dial(ttarget, opts...); err != nil {
		return nil, errors.Wrapf(err, "Dial Target %s", ttarget)
	} else {
		return conn, err
	}
}

func (t *Http2Trans) serverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	metric.Counter(monitor.HandleRPCTotal).Inc()
	rsp, err := handler(ctx, req)
	if err != nil {
		metric.Counter(monitor.HandleRPCFailTotal).Inc()
	}
	return rsp, err
}

func (t *Http2Trans) clientInterceptor(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	metric.Counter(monitor.CallPRCTotal).Inc()
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		if vlist, ok := md[string(meta.M3RouteType)]; ok && len(vlist) > 0 {
			if mesh.RouteType(vlist[0]) == mesh.RouteTypeBroad {
				metric.Counter(monitor.CallRPCFailTotal).Inc()
				return errs.TransportCliCantFindTopic.New("RouteTypeBroad not find Topic")
			}
		}
	}
	err := invoker(ctx, method, req, resp, cc, opts...)
	if err != nil {
		metric.Counter(monitor.CallRPCFailTotal).Inc()
	}
	return err
}
