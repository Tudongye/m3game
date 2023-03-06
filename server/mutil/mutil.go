package mutil

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
}

var (
	_ server.Server = (*Server)(nil)
)

func (s *Server) Init(cfg map[string]interface{}, app app.App) error {
	s.app = app
	return nil
}

func (s *Server) Type() server.Type {
	return server.Mutil
}

func (s *Server) Name() string {
	return fmt.Sprintf("%s.%s", server.Mutil, s.name)
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
	sctx := s.CreateContext(recv).(*Context)
	ctx := server.WithContext(sctx.Reciver().Ctx(), sctx)
	return sctx.Reciver().HandleMsg(ctx)
}

func (s *Server) SendInterFunc(sctx *transport.Sender) error {
	return s.app.SendInterFunc(sctx, runtime.SendInterFunc)
}

func (s *Server) TransportRegister() func(grpc.ServiceRegistrar) error {
	return nil
}
