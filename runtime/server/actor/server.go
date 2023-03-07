package actor

import (
	"fmt"
	"m3game/proto"
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/server"
	"m3game/runtime/transport"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func New(name string, creater ActorCreater) *Server {
	return &Server{
		name:     name,
		actormgr: newActorMgr(creater),
	}
}

type Server struct {
	name     string
	app      app.App
	actormgr *ActorMgr
}

type Config struct {
	ActiveTimeOut  int `mapstructure:"ActiveTimeOut"`
	SaveTimeInter  int `mapstructure:"SaveTimeInter"`
	TickTimeInter  int `mapstructure:"TickTimeInter"`
	MaxReqChanSize int `mapstructure:"MaxReqChanSize"`
	MaxReqWaitTime int `mapstructure:"MaxReqWaitTime"`
}

func (c *Config) CheckVaild() error {
	if c.ActiveTimeOut <= 0 {
		return fmt.Errorf("ActiveTimeOut %d invaild", c.ActiveTimeOut)
	}
	if c.SaveTimeInter <= 0 {
		return fmt.Errorf("SaveTimeInter %d invaild", c.SaveTimeInter)
	}
	if c.TickTimeInter <= 0 {
		return fmt.Errorf("TickTimeInter %d invaild", c.TickTimeInter)
	}
	if c.MaxReqChanSize <= 0 {
		return fmt.Errorf("MaxReqChanSize %d invaild", c.MaxReqChanSize)
	}
	if c.MaxReqWaitTime <= 0 {
		return fmt.Errorf("MaxReqWaitTime %d invaild", c.MaxReqWaitTime)
	}
	return nil
}

var (
	_cfg Config
	_    server.Server = (*Server)(nil)
)

func (s *Server) Init(c map[string]interface{}, app app.App) error {
	s.app = app
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return errors.Wrap(err, "Actor.Cfg Decode")
	}
	if err := _cfg.CheckVaild(); err != nil {
		return err
	}
	return nil
}

func (s *Server) Type() server.Type {
	return server.Actor
}

func (s *Server) Name() string {
	return fmt.Sprintf("%s.%s", server.Actor, s.name)
}

func (s *Server) Start(wg *sync.WaitGroup) error {
	return nil
}

func (s *Server) Stop() error {
	return nil
}

func (s *Server) Reload(map[string]interface{}) error {
	return nil
}

func (s *Server) RecvInterFunc(recv *transport.Reciver) (resp interface{}, err error) {
	if actorid, ok := recv.Metas().Get(proto.META_ACTORID); !ok {
		return nil, fmt.Errorf("no find actorid")
	} else {
		sctx := s.CreateContext(recv).(*Context)
		create := false
		if flag, ok := recv.Metas().Get(proto.META_CREATE_ACTOR); ok && flag == proto.META_FLAG_TRUE {
			create = true
		}
		return s.actormgr.recvInterFunc(actorid, create, sctx)
	}
}

func (s *Server) SendInterFunc(sctx *transport.Sender) error {
	return s.app.SendInterFunc(sctx, runtime.SendInterFunc)
}

func (s *Server) TransportRegister() func(grpc.ServiceRegistrar) error {
	return nil
}