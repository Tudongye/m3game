package test

import (
	"context"
	"log"
	"m3game/example/proto/pb"
	"m3game/plugins/gate/grpcgate"

	"google.golang.org/grpc"
)

func TTrace() {
	log.Println("TTrace 链路追踪测试...")
	m := make(map[string]string)
	m["m3routetype"] = "RouteTypeRandom"
	conn, err := grpc.Dial(
		_agenturl,
		grpc.WithInsecure())
	if err != nil {
		panic(err.Error())
	}
	cli := grpcgate.NewGGateSerClient(conn)
	stream, err := cli.CSTransport(context.Background())
	if err != nil {
		panic(err.Error())
	}
	{
		log.Println("Call.Auth 建立连接...")
		in := &pb.AuthReq{
			Token: "PlayerTest",
		}
		out := &pb.AuthRsp{}
		if err := CallGrpcGate(stream, "", m, in, out); err != nil {
			log.Printf("CallGrpcGate Fail %s", err.Error())
		}
		log.Println("PlayerID:", out.PlayerId)
	}
	{
		log.Println("Call.TraceHello 链路追踪接口...")
		in := &pb.TraceHello_Req{
			Req: "I'm test.",
		}
		out := &pb.TraceHello_Rsp{}
		if err := CallGrpcGate(stream, method_tracehello, m, in, out); err != nil {
			log.Printf("CallGrpcGate Fail %s", err.Error())
		}
		log.Println("Rsp:", out.Rsp)
	}
	stream.CloseSend()
}

func TBreak() {
	log.Println("TBreak 限流熔断测试...")
	m := make(map[string]string)
	m["m3routetype"] = "RouteTypeRandom"
	conn, err := grpc.Dial(
		_agenturl,
		grpc.WithInsecure())
	if err != nil {
		panic(err.Error())
	}
	cli := grpcgate.NewGGateSerClient(conn)
	stream, err := cli.CSTransport(context.Background())
	if err != nil {
		panic(err.Error())
	}
	{
		log.Println("Call.Auth 建立连接...")
		in := &pb.AuthReq{
			Token: "PlayerTest",
		}
		out := &pb.AuthRsp{}
		if err := CallGrpcGate(stream, "", m, in, out); err != nil {
			log.Printf("CallGrpcGate Fail %s", err.Error())
		}
		log.Println("PlayerID:", out.PlayerId)
	}
	for i := 0; i < 5; i++ {
		log.Printf("Call.BreakHello 第%d次请求...", i)
		in := &pb.BreakHello_Req{
			Req: "I'm test.",
		}
		out := &pb.BreakHello_Rsp{}
		if err := CallGrpcGate(stream, method_breakhello, m, in, out); err != nil {
			log.Printf("CallGrpcGate Fail %s", err.Error())
		}
		log.Println("Rsp:", out.Rsp)
	}
	stream.CloseSend()
}

func THello() {
	log.Println("THello 测试...")
	m := make(map[string]string)
	m["m3routetype"] = "RouteTypeRandom"
	conn, err := grpc.Dial(
		_agenturl,
		grpc.WithInsecure())
	if err != nil {
		panic(err.Error())
	}
	cli := grpcgate.NewGGateSerClient(conn)
	stream, err := cli.CSTransport(context.Background())
	if err != nil {
		panic(err.Error())
	}
	{
		log.Println("Call.Auth 建立连接...")
		in := &pb.AuthReq{
			Token: "PlayerTest",
		}
		out := &pb.AuthRsp{}
		if err := CallGrpcGate(stream, "", m, in, out); err != nil {
			log.Printf("CallGrpcGate Fail %s", err.Error())
		}
		log.Println("PlayerID:", out.PlayerId)
	}
	{
		log.Println("Call.Hello 测试接口...")
		in := &pb.Hello_Req{
			Req: "I'm test.",
		}
		out := &pb.Hello_Rsp{}
		if err := CallGrpcGate(stream, method_hello, m, in, out); err != nil {
			log.Printf("CallGrpcGate Fail %s", err.Error())
		}
		log.Println("Rsp:", out.Rsp)
	}
	stream.CloseSend()
}
