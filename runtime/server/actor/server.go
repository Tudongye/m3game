package actor

import (
	"context"
	"fmt"
	"m3game/meta/errs"
	"m3game/plugins/log"
	"m3game/runtime/app"
	"m3game/runtime/server"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Config struct {
	ActiveTimeOut        int    `mapstructure:"ActiveTimeOut" validate:"gte=0"`
	SaveTimeInter        int    `mapstructure:"SaveTimeInter" validate:"gte=0"`
	TickTimeInter        int    `mapstructure:"TickTimeInter" validate:"gt=0"`
	MaxReqChanSize       int    `mapstructure:"MaxReqChanSize" validate:"gt=0"`
	MaxReqWaitTime       int    `mapstructure:"MaxReqWaitTime" validate:"gt=0"`
	AutoCreate           int    `mapstructure:"AutoCreate" validate:"gte=0,lte=1"`
	LeaseMode            int    `mapstructure:"LeaseMode" validate:"gte=0,lte=1"`
	LeasePrefix          string `mapstructure:"LeasePrefix" validate:"required"`
	AllocLeaseTimeOut    int    `mapstructure:"AllocLeaseTimeOut" validate:"gt=0"`
	WaitFreeLeaseTimeOut int    `mapstructure:"WaitFreeLeaseTimeOut" validate:"gt=0"`
}

func New(name string, creater ActorCreater, actoridmetakey string) *Server {
	if actoridmetakey != strings.ToLower(actoridmetakey) {
		log.Fatal("ActorMetaKey Must be Lower Str but %s", actoridmetakey)
		return nil
	}
	return &Server{
		name:           name,
		creater:        creater,
		actoridmetakey: actoridmetakey,
	}
}

type Server struct {
	cfg            Config
	name           string
	app            app.App
	actormgr       *ActorMgr
	creater        ActorCreater
	actoridmetakey string
}

var (
	_ server.Server = (*Server)(nil)
)

func (s *Server) Init(c map[string]interface{}, app app.App) error {
	s.app = app
	if err := mapstructure.Decode(c, &s.cfg); err != nil {
		return errs.ActorServerInitFail.Wrap(err, "Actor.Cfg Decode")
	}
	validate := validator.New()
	if err := validate.Struct(&s.cfg); err != nil {
		return errs.ActorServerInitFail.Wrap(err, "Config")
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
		if vlist, ok := md[s.actoridmetakey]; ok && len(vlist) > 0 {
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
