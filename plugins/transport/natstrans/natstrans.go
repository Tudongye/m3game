package natstrans

import (
	"context"
	"fmt"
	"m3game/config"
	"m3game/meta"
	"m3game/meta/errs"
	"m3game/meta/metapb"
	"m3game/meta/monitor"
	"m3game/plugins/gate"
	"m3game/plugins/log"
	"m3game/plugins/metric"
	"m3game/plugins/transport"
	"m3game/runtime/plugin"
	"math/rand"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/mitchellh/mapstructure"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

var (
	_          plugin.Factory   = (*Factory)(nil)
	_          plugin.PluginIns = (*NatsTrans)(nil)
	_natstrans *NatsTrans
	_factory   = &Factory{}
)

const (
	_name            = "trans_nats"
	_maxserial       = 9999999999
	_callbacktimeout = 30
)

var gatecodec *gate.GateCodec

func init() {
	plugin.RegisterFactory(_factory)
}

type NatsTransCfg struct {
	URL string `mapstructure:"URL" validate:"required"` // nats地址
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
	if _natstrans != nil {
		return _natstrans, nil
	}
	var cfg NatsTransCfg
	if err := mapstructure.Decode(c, &cfg); err != nil {
		return nil, errs.TransportSetupFail.Wrap(err, "RedisDB Decode Cfg")
	}
	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, errs.TransportSetupFail.Wrap(err, "")
	}
	rand.Seed(time.Now().UnixNano())
	_natstrans = &NatsTrans{
		cfg:       cfg,
		subs:      make(map[string]*nats.Subscription),
		curserial: int64(rand.Intn(10000000)),
	}
	if nc, err := nats.Connect(cfg.URL); err != nil {
		return nil, errs.NatsSetupFail.Wrap(err, "Nats.Conntect %s", cfg.URL)
	} else {
		_natstrans.nc = nc
		if js, err := nc.JetStream(nats.PublishAsyncMaxPending(256)); err != nil {
			return nil, errs.NatsSetupFail.Wrap(err, "nc.JetStream %s", cfg.URL)
		} else {
			_natstrans.js = js
		}
	}
	_natstrans.RegisterServerInterceptor(_natstrans.serverInterceptor)
	_natstrans.RegisterClientInterceptor(_natstrans.clientInterceptor)
	if _, err := transport.New(_natstrans); err != nil {
		return nil, err
	}
	return _natstrans, nil
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

type NatsTrans struct {
	cfg                NatsTransCfg
	nc                 *nats.Conn
	js                 nats.JetStreamContext
	subs               map[string]*nats.Subscription
	handlermap         sync.Map
	serverInterceptors []grpc.UnaryServerInterceptor
	clientInterceptors []grpc.UnaryClientInterceptor
	sserverInterceptor grpc.UnaryServerInterceptor
	curserial          int64
	mu                 sync.Mutex
	callbackmap        sync.Map
	balancermap        sync.Map
}

type CallBack struct {
	createtime int64
	respchan   chan *TransMsg
}

func NewCallBack() *CallBack {
	return &CallBack{
		createtime: time.Now().Unix(),
		respchan:   make(chan *TransMsg),
	}
}

func (t *NatsTrans) Factory() plugin.Factory {
	return _factory
}

func (t *NatsTrans) Host() string {
	return t.cfg.URL
}

func (t *NatsTrans) Port() int {
	return 0
}

func (t *NatsTrans) Prepare(ctx context.Context) error {
	t.sserverInterceptor = grpc_middleware.ChainUnaryServer(t.serverInterceptors...)
	return nil
}
func (t *NatsTrans) Start(ctx context.Context) error {
	go func() {
		select {
		case <-ctx.Done():
			t.nc.Close()
		}
	}()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// 清理过期callback
				timenow := time.Now().Unix()
				var timeouts []int64
				t.callbackmap.Range(func(key, value any) bool {
					cb := value.(*CallBack)
					if cb.createtime+_callbacktimeout < timenow {
						// 超时
						select {
						case cb.respchan <- nil:
							timeouts = append(timeouts, key.(int64))
						default:
							// 恰好收到回包
						}
					}
					return true
				})
				for _, serial := range timeouts {
					t.callbackmap.Delete(serial)
				}
			}
		}
	}()

	appTopic := appTopic(config.GetIDStr())
	log.Info("Sub %s", appTopic)
	if sub, err := t.nc.Subscribe(appTopic, func(m *nats.Msg) {
		go func() {
			if err := t.recvbytes(m.Data); err != nil {
				log.Error("broker subscribe %s handler err %s", appTopic, err.Error())
			}
		}()
	}); err != nil {
		return errs.BrokerSerSetBrokerFail.Wrap(err, "Subscribe %s", appTopic)
	} else {
		t.subs[appTopic] = sub
	}

	svcTopic := svcTopic(config.GetSvcID().String())
	log.Info("Sub %s", svcTopic)
	if sub, err := t.nc.Subscribe(svcTopic, func(m *nats.Msg) {
		go func() {
			if err := t.recvbytes(m.Data); err != nil {
				log.Error("broker subscribe %s handler err %s", svcTopic, err.Error())
			}
		}()
	}); err != nil {
		return errs.BrokerSerSetBrokerFail.Wrap(err, "Subscribe %s", svcTopic)
	} else {
		t.subs[svcTopic] = sub
	}
	return nil
}

func (t *NatsTrans) RegisterSer(f func(grpc.ServiceRegistrar) error) error {
	if err := f(t); err != nil {
		return errs.TransportRegisterSerFail.Wrap(err, "server.register")
	}
	return nil
}

func (t *NatsTrans) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	for _, it := range sd.Methods {
		path := fmt.Sprintf("/%v/%v", sd.ServiceName, it.MethodName)
		handler := it.Handler
		f := func(ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
			return handler(ss, ctx, dec, interceptor)
		}
		t.handlermap.Store(path, f)
		log.Info("Register BrokerMethod => %v", path)
	}
}

func (t *NatsTrans) RegisterServerInterceptor(f grpc.UnaryServerInterceptor) {
	t.serverInterceptors = append(t.serverInterceptors, f)
}

func (t *NatsTrans) RegisterClientInterceptor(f grpc.UnaryClientInterceptor) {
	t.clientInterceptors = append([]grpc.UnaryClientInterceptor{f}, t.clientInterceptors...)
}

func (t *NatsTrans) ClientInterceptors() []grpc.UnaryClientInterceptor {
	return t.clientInterceptors
}

func (t *NatsTrans) ClientConn(target string, opts ...grpc.DialOption) (grpc.ClientConnInterface, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	opts = append(opts,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(transport.Instance().ClientInterceptors()...)),
		grpc.WithTimeout(time.Second*10),
	)

	if conn, err := grpc.Dial(target, opts...); err != nil {
		return nil, errors.Wrapf(err, "Dial Target %s", target)
	} else {
		// 开启Balance
		if _, ok := t.balancermap.Load(target); !ok {
			b := NewBalancer(target)
			t.balancermap.Store(target, b)
		}
		return conn, err
	}
}

func (t *NatsTrans) serverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	metric.Counter(monitor.HandleRPCTotal).Inc()
	rsp, err := handler(ctx, req)
	if err != nil {
		metric.Counter(monitor.HandleRPCFailTotal).Inc()
	}
	return rsp, err
}

func (t *NatsTrans) clientInterceptor(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	metric.Counter(monitor.CallPRCTotal).Inc()
	if b, ok := t.balancermap.Load(cc.Target()); !ok {
		log.Fatal("Not Find Balancer %s", cc.Target())
		return errs.NatsTransportBalanceNotFind.New("")
	} else {
		if dstappid, broad, err := b.(*NatsBalancer).Pick(ctx); err != nil {
			log.Error(err.Error())
			return err
		} else {
			if broad {
				err := t.sendreq(ctx, svcTopic(cc.Target()), method, req, nil)
				if err != nil {
					metric.Counter(monitor.CallRPCFailTotal).Inc()
				}
				return err
			} else {
				err := t.sendreq(ctx, appTopic(dstappid), method, req, resp)
				if err != nil {
					metric.Counter(monitor.CallRPCFailTotal).Inc()
				}
				return err
			}
		}
	}
}

func (t *NatsTrans) dispatch(msg *TransMsg) error {
	if msg.Ack {
		// 是回包
		if v, ok := t.callbackmap.LoadAndDelete(msg.Serial); !ok {
			return errs.NatsTransportAckNotfindCB.New(fmt.Sprintf("%d", msg.Serial))
		} else {
			cb := v.(*CallBack)
			select {
			case cb.respchan <- msg:
				return nil
			default:
				return errs.NatsTransportAckCBTimeOut.New(fmt.Sprintf("%d", msg.Serial))
			}
		}
	}
	if value, ok := t.handlermap.Load(msg.Method); !ok {
		return errs.BrokerSerHandlerNotFind.New("not find method %s", msg.Method)
	} else {
		handlerfunc := value.(func(ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error))
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		metad := make(map[string]string)
		for _, meta := range msg.Metas {
			metad[meta.Key] = meta.Value
		}
		md := metadata.New(metad)
		ctx = metadata.NewIncomingContext(ctx, md)
		// 进入业务接口
		if resp, err := handlerfunc(ctx,
			func(pkg interface{}) error {
				return proto.Unmarshal(msg.Content, pkg.(proto.Message))
			},
			t.sserverInterceptor); err != nil {
			if !msg.Nty {
				return t.sendresp(msg, nil, err)
			}
			return nil
		} else {
			if !msg.Nty {
				return t.sendresp(msg, resp, nil)
			}
			return nil
		}
	}
}

func (t *NatsTrans) recvbytes(bytes []byte) error {
	log.Debug("recvbytes %v", len(bytes))
	bmsg := &TransMsg{}
	if err := proto.Unmarshal(bytes, bmsg); err != nil {
		return errs.MsgUnmarshFail.New("Unmarshal bytes err %s", err.Error())
	}
	return t.dispatch(bmsg)
}

func (t *NatsTrans) callback(msg *TransMsg, resp interface{}, isgate bool) error {
	if msg.ErrCode != 0 {
		// 返回错误
		return errs.New(int(msg.ErrCode), msg.ErrContent)
	}
	// 解析
	if isgate {
		return gatecodec.Unmarshal(msg.Content, resp)

	} else {
		return proto.Unmarshal(msg.Content, resp.(proto.Message))

	}
}

func (t *NatsTrans) sendreq(ctx context.Context, topic string, method string, req interface{}, resp interface{}) error {
	bmsg := &TransMsg{}
	bmsg.Method = method
	bmsg.Serial = t.allocSerial()
	bmsg.Ack = false
	bmsg.SrcApp = config.GetIDStr()
	isgate := false
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		for k, vlist := range md {
			if len(vlist) > 0 {
				bmsg.Metas = append(bmsg.Metas, &metapb.Meta{Key: k, Value: vlist[0]})
				if k == meta.M3GateMsg.String() && vlist[0] == "1" {
					isgate = true
				}
			}
		}
	}
	var err error
	if isgate {
		if bmsg.Content, err = gatecodec.Marshal(req); err != nil {
			return errs.MsgMarshFail.Wrap(err, "")
		}
	} else {
		if bmsg.Content, err = proto.Marshal(req.(proto.Message)); err != nil {
			return errs.MsgMarshFail.Wrap(err, "")
		}

	}
	if resp == nil {
		// Nty 发出后返回
		bmsg.Nty = true
		return t.sendbytes(topic, bmsg)
	} else {
		cb := NewCallBack()
		t.callbackmap.Store(bmsg.Serial, cb)
		defer t.callbackmap.Delete(bmsg.Serial)
		// 发包
		if err := t.sendbytes(topic, bmsg); err != nil {
			return err
		}
		// 等待回包
		select {
		case respmsg := <-cb.respchan:
			if respmsg == nil {
				// 超时
				return errs.NatsTransportReqTimeOut.New("")
			}
			// 正常回包
			return t.callback(respmsg, resp, isgate)
		case <-ctx.Done():
			// ctx取消
			return errs.NatsTransportReqCtxDone.New("")
		}
	}
}

func (t *NatsTrans) sendresp(msg *TransMsg, resp interface{}, resperr error) error {
	bmsg := &TransMsg{}
	bmsg.Method = msg.Method
	bmsg.SrcApp = config.GetIDStr()
	bmsg.Ack = true
	bmsg.Serial = msg.Serial
	if resp != nil {
		var err error
		if bmsg.Content, err = proto.Marshal(resp.(proto.Message)); err != nil {
			return errs.MsgMarshFail.Wrap(err, "")
		}
	}
	if resperr != nil {
		bmsg.ErrCode = 1
		bmsg.ErrContent = resperr.Error()
	}
	return t.sendbytes(appTopic(msg.SrcApp), bmsg)
}

func (t *NatsTrans) sendbytes(topic string, msg *TransMsg) error {
	if bytes, err := proto.Marshal(msg); err != nil {
		return errs.MsgMarshFail.Wrap(err, "")
	} else {
		log.Debug("sendbytes %s %v", topic, len(bytes))
		return t.nc.Publish(topic, bytes)
	}
}

func (t *NatsTrans) allocSerial() int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.curserial++
	t.curserial = _maxserial % _maxserial
	return t.curserial
}

func appTopic(app string) string {
	return fmt.Sprintf("m3app-%s", app)
}

func svcTopic(svc string) string {
	return fmt.Sprintf("m3svc-%s", svc)
}
