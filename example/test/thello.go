package test

import (
	"log"
	"m3game/example/proto"
	"m3game/example/proto/pb"
	"m3game/plugins/agent"
)

func THello() {
	log.Println("THello...")
	log.Printf("Auth %s...\n", Token1)
	req := agent.AuthPara{
		Type:  agent.UType(1),
		Token: Token1,
	}
	rsp, err := CallAuth(&req)
	if err != nil {
		log.Printf("CallAuth Fail %v", err)
		return
	}
	log.Printf("Res: %s %s\n", rsp.Uid, rsp.Session)
	{
		log.Println("Call Hello...")
		var out pb.Hello_Rsp
		in := &pb.Hello_Req{
			RouteHead: GenRandomRouteHead(rsp.EnvID, rsp.WorldID, proto.MutilAppFuncID),
			Req:       "I'm test.",
		}
		if err := CallLogic(rsp.Uid, rsp.Session, method_hello, in, &out); err != nil {
			log.Printf("Call Fail %v", err)
			return
		}
		log.Printf("Res: %s\n", out.Rsp)
	}
}

func TTrace() {
	log.Println("TTrace...")
	log.Printf("Auth %s...\n", Token1)
	req := agent.AuthPara{
		Type:  agent.UType(1),
		Token: Token1,
	}
	rsp, err := CallAuth(&req)
	if err != nil {
		log.Printf("CallAuth Fail %v", err)
		return
	}
	log.Printf("Res: %s %s\n", rsp.Uid, rsp.Session)
	{
		log.Println("Call TraceHello...")
		var out pb.TraceHello_Rsp
		in := &pb.TraceHello_Req{
			RouteHead: GenRandomRouteHead(rsp.EnvID, rsp.WorldID, proto.MutilAppFuncID),
			Req:       "I'm test.",
		}
		if err := CallLogic(rsp.Uid, rsp.Session, method_tracehello, in, &out); err != nil {
			log.Printf("Call Fail %v", err)
			return
		}
		log.Printf("Res: %s\n", out.Rsp)
	}
}

func TBreak() {
	log.Println("TBreak...")
	log.Printf("Auth %s...\n", Token1)
	req := agent.AuthPara{
		Type:  agent.UType(1),
		Token: Token1,
	}
	rsp, err := CallAuth(&req)
	if err != nil {
		log.Printf("CallAuth Fail %v", err)
		return
	}
	log.Printf("Res: %s %s\n", rsp.Uid, rsp.Session)
	for i := 0; i < 10; i++ {
		log.Printf("Call BreakHello...%d\n", i)
		var out pb.BreakHello_Rsp
		in := &pb.BreakHello_Req{
			RouteHead: GenRandomRouteHead(rsp.EnvID, rsp.WorldID, proto.MutilAppFuncID),
			Req:       "I'm test.",
		}
		if err := CallLogic(rsp.Uid, rsp.Session, method_breakhello, in, &out); err != nil {
			log.Printf("Call Fail %v", err)
			return
		}
		log.Printf("Res: %s\n", out.Rsp)
	}
}
