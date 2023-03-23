package grpcgate

import (
	"context"
	"io"
	"m3game/config"
	"m3game/meta"
	"m3game/meta/metapb"
	"m3game/plugins/gate"
	"m3game/plugins/log"
	"m3game/runtime/plugin"
	"net"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

var (
	_         plugin.Factory   = (*Factory)(nil)
	_         plugin.PluginIns = (*Gate)(nil)
	_grpcgate *Gate
	_factory  = &Factory{}
)

const (
	_name = "gate_grpc"
)

func init() {
	plugin.RegisterFactory(_factory)
}

type grpcGateCfg struct {
	Addr string `mapstructure:"Addr" validate:"required,tcp4_addr"`
}

type Factory struct {
}

func (f *Factory) Type() plugin.Type {
	return plugin.Gate
}
func (f *Factory) Name() string {
	return _name
}

func (f *Factory) Setup(c map[string]interface{}) (plugin.PluginIns, error) {
	if _grpcgate != nil {
		return _grpcgate, nil
	}
	var cfg grpcGateCfg
	if err := mapstructure.Decode(c, &cfg); err != nil {
		return nil, errors.Wrap(err, "Gate Decode Cfg")
	}
	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, err
	}
	var err error
	tcpAddr, err := net.ResolveTCPAddr("tcp", cfg.Addr)
	if err != nil {
		return nil, errors.Wrap(err, "transport")
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, errors.Wrap(err, "transport.ListenTCP")
	}
	_grpcgate = &Gate{
		gser: grpc.NewServer(),
		no:   1,
	}
	RegisterGateSerServer(_grpcgate.gser, _grpcgate)
	go func() {
		log.Info("GrpcGate Listen %s", cfg.Addr)
		if err := _grpcgate.gser.Serve(listener); err != nil {
			log.Error("GrpcGate Err %s", err.Error())
		}
		_grpcgate.isstoped = true
	}()
	if _, err := gate.New(_grpcgate); err != nil {
		return nil, err
	}
	return _grpcgate, nil
}

func (f *Factory) Destroy(plugin.PluginIns) error {
	return nil
}

func (f *Factory) Reload(plugin.PluginIns, map[string]interface{}) error {
	return nil
}

func (f *Factory) CanUnload(p plugin.PluginIns) bool {
	g := p.(*Gate)
	return g.isstoped
}

type Gate struct {
	conns sync.Map
	gser  *grpc.Server
	UnimplementedGateSerServer
	mutex    sync.RWMutex
	isstoped bool
	no       int
}

func (g *Gate) Factory() plugin.Factory {
	return _factory
}
func (g *Gate) GetConn(playerid string) gate.CSConn {
	if c, ok := g.conns.Load(playerid); ok {
		return c.(gate.CSConn)
	}
	return nil
}
func (g *Gate) GenConnID() int {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.no += 1
	return g.no
}

func (g *Gate) CSTransport(srv GateSer_CSTransportServer) error {
	log.Debug("Recv CSTransport")
	var playerid string

	// 连接鉴权
	msg, err := srv.Recv()
	if err != nil {
		return err
	}
	authreq := &metapb.AuthReq{}
	if err := proto.Unmarshal(msg.Content, authreq); err != nil {
		return err
	}
	if res, err := gate.AuthCall(authreq); err != nil {
		log.Error("%v", err)
		return err
	} else {
		playerid = res.PlayerID
		msg.Content, _ = proto.Marshal(res)
		srv.Send(msg)
	}

	no := g.GenConnID()
	csconn := func(n int) *CSConn {
		g.mutex.Lock()
		defer g.mutex.Unlock()
		if c, ok := g.conns.Load(playerid); ok {
			c.(gate.CSConn).Kick()
			g.conns.Delete(playerid)
		}
		ctx, cancel := context.WithCancel(context.Background())
		csconn := &CSConn{
			srv:      srv,
			no:       n,
			sendch:   make(chan *metapb.CSMsg, 10),
			isclosed: false,
			playerid: playerid,
			ctx:      ctx,
			cancel:   cancel,
		}
		g.conns.Store(playerid, csconn)
		return csconn
	}(no)
	go csconn.recvloop()
	go csconn.sendloop()

	<-csconn.ctx.Done()

	g.mutex.Lock()
	defer g.mutex.Unlock()
	if c, ok := g.conns.Load(playerid); ok {
		if c.(*CSConn).no == csconn.no {
			g.conns.Delete(playerid)
		}
	}
	return nil
}

type CSConn struct {
	srv      GateSer_CSTransportServer
	sendch   chan *metapb.CSMsg
	isclosed bool
	no       int
	playerid string
	mutex    sync.Mutex
	ctx      context.Context
	cancel   context.CancelFunc
}

func (c *CSConn) Send(ctx context.Context, msg *metapb.CSMsg) error {
	if c.isclosed {
		return errors.New("CSConn closed")
	}
	select {
	case <-c.ctx.Done():
		return errors.New("CSConn closed")
	case c.sendch <- msg:
		return nil
	case <-ctx.Done():
		return errors.New("ctx done")
	default:
		return errors.New("chan full")
	}
}

func (c *CSConn) Kick() {
	c.safeclose()
}

func (c *CSConn) safeclose() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.isclosed {
		return
	}
	c.isclosed = true
	c.cancel()
}

func (c *CSConn) recvloop() {
	for {
		msg, err := c.srv.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Error("srv Recv %s", err.Error())
			break
		}
		log.Debug("Recv %s", msg.Method)
		msg.Metas = append(msg.Metas, &metapb.Meta{Key: meta.M3ActorActorID.String(), Value: c.playerid})
		msg.Metas = append(msg.Metas, &metapb.Meta{Key: meta.M3RouteSrcApp.String(), Value: config.GetAppID().String()})
		res, err := gate.LogicCall(c.playerid, msg)
		if err != nil {
			log.Error("Call Logic %s %s", msg.Method, err.Error())
			break
		} else {
			if err := c.Send(context.Background(), res); err != nil {
				log.Error("Send %s %s", msg.Method, err.Error())
			}

		}
	}
	c.safeclose()
}

func (c *CSConn) sendloop() {
	for msg := range c.sendch {
		if err := c.srv.Send(msg); err != nil {
			break
		}
	}
	c.safeclose()
}
