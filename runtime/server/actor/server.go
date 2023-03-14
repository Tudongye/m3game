package actor

import (
	"context"
	"fmt"
	"m3game/meta"
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/server"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

func (s *Server) ServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	var actorid string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vlist, ok := md[string(meta.M3ActorActorID)]; ok && len(vlist) > 0 {
			actorid = vlist[0]
		}
	}
	sctx := server.GenContext(s)
	ctx = server.WithContext(ctx, sctx)
	return s.actormgr.CallFunc(actorid, ctx, req, info, handler)

}
func (s *Server) ClientInterceptor(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return runtime.ClientInterceptor(ctx, method, req, resp, cc, invoker, opts...)
}

func (s *Server) TransportRegister() func(grpc.ServiceRegistrar) error {
	return nil
}
