package async

import (
	"m3game/runtime/server"
	"m3game/runtime/transport"
)

type Context struct {
	reciver *transport.Reciver
	server  *Server
}

func (c *Context) Server() server.Server {
	return c.server
}

func (c *Context) Reciver() *transport.Reciver {
	return c.reciver
}

func (s *Server) CreateContext(rec *transport.Reciver) server.Context {
	return &Context{
		reciver: rec,
		server:  s,
	}
}
