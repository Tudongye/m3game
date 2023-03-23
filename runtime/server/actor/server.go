package actor

import (
	"context"
	"fmt"
	"m3game/meta"
	"m3game/runtime/app"
	"m3game/runtime/server"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Config struct {
	ActiveTimeOut        int    `mapstructure:"ActiveTimeOut" validate:"gt=0"`
	SaveTimeInter        int    `mapstructure:"SaveTimeInter" validate:"gt=0"`
	TickTimeInter        int    `mapstructure:"TickTimeInter" validate:"gt=0"`
	MaxReqChanSize       int    `mapstructure:"MaxReqChanSize" validate:"gt=0"`
	MaxReqWaitTime       int    `mapstructure:"MaxReqWaitTime" validate:"gt=0"`
	LeaseMode            int    `mapstructure:"LeaseMode" validate:"gte=0,lte=1"`
	LeasePrefix          string `mapstructure:"LeasePrefix" validate:"gt=0"`
	AllocLeaseTimeOut    int    `mapstructure:"AllocLeaseTimeOut" validate:"required"`
	WaitFreeLeaseTimeOut int    `mapstructure:"WaitFreeLeaseTimeOut" validate:"gt=0"`
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
	validate := validator.New()
	if err := validate.Struct(&s.cfg); err != nil {
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
