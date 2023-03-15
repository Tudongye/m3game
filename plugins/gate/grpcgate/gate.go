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

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

var (
	_         plugin.Factory   = (*Factory)(nil)
	_         plugin.PluginIns = (*Gate)(nil)
	_instance *Gate
	_factory  = &Factory{}
	_cfg      grpcGateCfg
)

const (
	_factoryname = "gate_grpc"
)

func init() {
	plugin.RegisterFactory(_factory)
}

type grpcGateCfg struct {
	Addr string `mapstructure:"Addr"`
}

func (c *grpcGateCfg) CheckVaild() error {
	if c.Addr == "" {
		return errors.New("Addr cant be space")
	}
	return nil
}

type Factory struct {
}

func (f *Factory) Type() plugin.Type {
	return plugin.Gate
}
func (f *Factory) Name() string {
	return _factoryname
}

func (f *Factory) Setup(c map[string]interface{}) (plugin.PluginIns, error) {
	if _instance != nil {
		return _instance, nil
	}
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return nil, errors.Wrap(err, "Gate Decode Cfg")
	}
	if err := _cfg.CheckVaild(); err != nil {
		return nil, err
	}
	var err error
	tcpAddr, err := net.ResolveTCPAddr("tcp", _cfg.Addr)
	if err != nil {
		return nil, errors.Wrap(err, "transport")
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, errors.Wrap(err, "transport.ListenTCP")
	}
	_instance := &Gate{
		gser:  grpc.NewServer(),
		conns: make(map[string]*CSConn),
		no:    1,
	}
	RegisterGateSerServer(_instance.gser, _instance)
	go func() {
		log.Info("GrpcGate Listen %s", _cfg.Addr)
		if err := _instance.gser.Serve(listener); err != nil {
			log.Error("GrpcGate Err %s", err.Error())
		}
		_instance.isstoped = true
	}()
	gate.Set(_instance)
	return _instance, nil
}

func (f *Factory) Destroy(plugin.PluginIns) error {
	return nil
}

func (f *Factory) Reload(plugin.PluginIns, map[string]interface{}) error {
	return nil
}

func (f *Factory) CanDelete(p plugin.PluginIns) bool {
	g := p.(*Gate)
	return g.isstoped
}

type Gate struct {
	conns map[string]*CSConn
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
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	if c, ok := g.conns[playerid]; ok {
		return c
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
	{
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
			return err
		} else {
			playerid = res.PlayerID
			msg.Content, _ = proto.Marshal(res)
			srv.Send(msg)
		}
	}
	no := g.GenConnID()
	csconn := func(n int) *CSConn {
		g.mutex.Lock()
		defer g.mutex.Unlock()
		if c, ok := g.conns[playerid]; ok {
			c.Kick()
		}
		g.conns[playerid] = &CSConn{
			srv:      srv,
			no:       n,
			sendch:   make(chan *metapb.CSMsg, 10),
			exitch:   make(chan struct{}),
			isclosed: false,
			playerid: playerid,
		}
		return g.conns[playerid]
	}(no)
	go csconn.recvloop()
	go csconn.sendloop()
	<-csconn.exitch
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if c, ok := g.conns[playerid]; ok {
		if c.no == csconn.no {
			delete(g.conns, playerid)
		}
	}
	return nil
}

type CSConn struct {
	srv      GateSer_CSTransportServer
	sendch   chan *metapb.CSMsg
	exitch   chan struct{}
	isclosed bool
	no       int
	playerid string
	mutex    sync.Mutex
}

func (c *CSConn) Send(ctx context.Context, msg *metapb.CSMsg) error {
	if c.isclosed {
		return errors.New("CSConn closed")
	}
	select {
	case <-c.exitch:
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
	close(c.exitch)
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
		msg.Metas = append(msg.Metas, &metapb.Meta{Key: meta.M3PlayerID.String(), Value: c.playerid})
		msg.Metas = append(msg.Metas, &metapb.Meta{Key: meta.M3RouteSrcApp.String(), Value: config.GetAppID().String()})
		res, err := gate.LogicCall(msg)
		if err != nil {
			log.Error("Call Logic %s %s", msg.Method, err.Error())
			break
		} else {
			log.Debug("Send %s %v", msg.Method, c.Send(context.Background(), res))

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
