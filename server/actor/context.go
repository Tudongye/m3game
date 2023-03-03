package actor

import (
	"context"
	"m3game/runtime/transport"
	"m3game/server"
)

func ParseActor(ctx context.Context) Actor {
	if c := server.ParseContext(ctx); c == nil {
		return nil
	} else if sctx := c.(*Context); sctx == nil {
		return nil
	} else {
		return sctx.Actor()
	}
}

type Context struct {
	reciver *transport.Reciver
	server  *Server
	actor   Actor
}

func (c *Context) Server() server.Server {
	return c.server
}

func (c *Context) Reciver() *transport.Reciver {
	return c.reciver
}

func (c *Context) Actor() Actor {
	return c.actor
}

func (c *Context) SetActor(a Actor) {
	c.actor = a
}

func (s *Server) CreateContext(rec *transport.Reciver) server.Context {
	return &Context{
		reciver: rec,
		server:  s,
	}
}

type actorReq struct {
	ctx     *Context
	rspchan chan *actorRsp
}

type actorRsp struct {
	rsp interface{}
	err error
}

func newActorReq(ctx *Context) *actorReq {
	return &actorReq{
		ctx:     ctx,
		rspchan: make(chan *actorRsp),
	}
}
