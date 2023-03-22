package actor

import (
	"context"
	"fmt"
	"m3game/meta"
	"m3game/runtime/app"
	"m3game/runtime/server"
	"m3game/util"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Config struct {
	ActiveTimeOut        int    `mapstructure:"ActiveTimeOut"`
	SaveTimeInter        int    `mapstructure:"SaveTimeInter"`
	TickTimeInter        int    `mapstructure:"TickTimeInter"`
	MaxReqChanSize       int    `mapstructure:"MaxReqChanSize"`
	MaxReqWaitTime       int    `mapstructure:"MaxReqWaitTime"`
	LeaseMode            int    `mapstructure:"LeaseMode"`
	LeasePrefix          string `mapstructure:"LeasePrefix"`
	AllocLeaseTimeOut    int    `mapstructure:"AllocLeaseTimeOut"`
	WaitFreeLeaseTimeOut int    `mapstructure:"WaitFreeLeaseTimeOut"`
}

func (c *Config) checkValid() error {
	if err := util.GreatInt(c.ActiveTimeOut, 0, "ActiveTimeOut"); err != nil {
		return err
	}
	if err := util.GreatInt(c.SaveTimeInter, 0, "SaveTimeInter"); err != nil {
		return err
	}
	if err := util.GreatInt(c.TickTimeInter, 0, "TickTimeInter"); err != nil {
		return err
	}
	if err := util.GreatInt(c.MaxReqChanSize, 0, "MaxReqChanSize"); err != nil {
		return err
	}
	if err := util.GreatInt(c.MaxReqWaitTime, 0, "MaxReqWaitTime"); err != nil {
		return err
	}

	if c.LeaseMode == 1 {
		if err := util.InEqualStr(c.LeasePrefix, "", "LeasePrefix"); err != nil {
			return err
		}
		if err := util.InEqualInt(c.AllocLeaseTimeOut, 0, "AllocLeaseTimeOut"); err != nil {
			return err
		}
		if err := util.InEqualInt(c.WaitFreeLeaseTimeOut, 0, "WaitFreeLeaseTimeOut"); err != nil {
			return err
		}
	}
	return nil
}

func New(name string, creater ActorCreater) *Server {
	return &Server{
		name:    name,
		creater: creater,
	}
}

type Server struct {
	cfg      Config
	name     string
	app      app.App
	actormgr *ActorMgr
	creater  ActorCreater
}

var (
	_ server.Server = (*Server)(nil)
)

func (s *Server) Init(c map[string]interface{}, app app.App) error {
	s.app = app
	if err := mapstructure.Decode(c, &s.cfg); err != nil {
		return errors.Wrap(err, "Actor.Cfg Decode")
	}
	if err := s.cfg.checkValid(); err != nil {
		return err
	}
	s.actormgr = newActorMgr(s.creater, &s.cfg)
	return nil
}

func (s *Server) Type() server.Type {
	return server.Actor
}

func (s *Server) Name() string {
	return fmt.Sprintf("%s.%s", server.Actor, s.name)
}

func (s *Server) Prepare(context.Context) error {
	return nil
}
func (s *Server) Start(context.Context) {

}

func (s *Server) Reload(map[string]interface{}) error {
	return nil
}

func (s *Server) ServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	var actorid string
	if md, ok := metadata.FromIncomingContext(ctx); ok && md != nil {
		if vlist, ok := md[string(meta.M3ActorActorID)]; ok && len(vlist) > 0 {
			actorid = vlist[0]
		}
	}
	return s.actormgr.CallFunc(actorid, ctx, req, info, handler)
}

func (s *Server) ClientInterceptor(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return invoker(ctx, method, req, resp, cc, opts...)
}

func (s *Server) TransportRegister() func(grpc.ServiceRegistrar) error {
	return nil
}

func (s *Server) ActorMgr() *ActorMgr {
	return s.actormgr
}
