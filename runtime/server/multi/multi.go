package multi

import (
	"context"
	"fmt"
	"m3game/runtime/app"
	"m3game/runtime/server"

	"google.golang.org/grpc"
)

func New(name string) *Server {
	return &Server{
		name: name,
	}
}

type Server struct {
	name string
	app  app.App
}

var (
	_ server.Server = (*Server)(nil)
)

func (s *Server) Init(cfg map[string]interface{}, app app.App) error {
	s.app = app
	return nil
}

func (s *Server) Type() server.Type {
	return server.Multi
}

func (s *Server) Name() string {
	return fmt.Sprintf("%s.%s", server.Multi, s.name)
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
	return handler(ctx, req)
}

func (s *Server) ClientInterceptor(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return invoker(ctx, method, req, resp, cc, opts...)
}

func (s *Server) TransportRegister() func(grpc.ServiceRegistrar) error {
	return nil
}
