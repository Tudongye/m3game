package async

import (
	"fmt"
	"m3game/app"
	"m3game/runtime"
	"m3game/runtime/transport"
	"m3game/server"
	"sync"

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
	lock sync.Mutex
}

var (
	_ server.Server = (*Server)(nil)
)

func (s *Server) Init(cfg map[string]interface{}, app app.App) error {
	s.app = app
	return nil
}

func (s *Server) Type() server.Type {
	return server.Async
}

func (s *Server) Name() string {
	return fmt.Sprintf("%s.%s", server.Async, s.name)
}

func (s *Server) Start(wg *sync.WaitGroup) error {
	return nil
}

func (s *Server) Stop() error {
	return nil
}

func (s *Server) Reload() error {
	return nil
}
func (s *Server) RecvInterFunc(recv *transport.Reciver) (resp interface{}, err error) {
	sctx := s.CreateContext(recv).(*Context)
	ctx := server.WithContext(sctx.Reciver().Ctx(), sctx)
	s.lock.Lock()
	defer s.lock.Unlock()
	return sctx.Reciver().HandleMsg(ctx)
}

func (s *Server) SendInterFunc(sctx *transport.Sender) error {
	s.lock.Unlock()
	defer s.lock.Lock()
	return s.app.SendInterFunc(sctx, runtime.SendInterFunc)
}

func (s *Server) TransportRegister() func(grpc.ServiceRegistrar) error {
	return nil
}
