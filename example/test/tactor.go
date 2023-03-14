package test

import (
	"context"
	"log"
	"m3game/example/proto/pb"
	"m3game/meta/metapb"
	"m3game/plugins/gate/grpcgate"

	"google.golang.org/grpc"
)

func TActorCommon() {
	log.Println("TActorCommon Actor常规用例...")
	m := make(map[string]string)
	m["m3routetype"] = "RouteTypeRandom"
	conn, err := grpc.Dial(
		_agenturl,
		grpc.WithInsecure())
	if err != nil {
		panic(err.Error())
	}
	cli := grpcgate.NewGateSerClient(conn)
	stream, err := cli.CSTransport(context.Background())
	if err != nil {
		panic(err.Error())
	}

	{
		log.Println("Call.Auth 建立连接...")
		in := &metapb.AuthReq{
			Token: "PlayerTest",
		}
		out := &metapb.AuthRsp{}
		if err := CallGrpcGate(stream, "", m, in, out); err != nil {
			log.Printf("CallGrpcGate Fail %s", err.Error())
		}
		log.Println("PlayerID:", out.PlayerID)
	}
	var actorid string
	{
		log.Println("Call.Register 注册接口,调用ActorApp的ActorRegSer...")
		out := &pb.Register_Rsp{}
		in := &pb.Register_Req{
			Name: "June",
		}
		if err := CallGrpcGate(stream, method_register, m, in, out); err != nil {
			log.Printf("CallGrpcGate Fail %s", err.Error())
			return
		}
		log.Println("Rsp-ActorID:", out.ActorID)
		actorid = out.ActorID
	}
	m["m3actoractorid"] = actorid
	{
		log.Println("Call.Login 登陆接口...")
		out := &pb.Login_Rsp{}
		in := &pb.Login_Req{
			ActorID: actorid,
		}
		if err := CallGrpcGate(stream, method_login, m, in, out); err != nil {
			log.Printf("CallGrpcGate Fail %s", err.Error())
			return
		}
		log.Println("Rsp:", out)
	}
	{
		log.Println("Call.GetInfo 获取角色信息...")
		out := &pb.GetInfo_Rsp{}
		in := &pb.GetInfo_Req{
			ActorID: actorid,
		}
		if err := CallGrpcGate(stream, method_getinfo, m, in, out); err != nil {
			log.Printf("CallGrpcGate Fail %s", err.Error())
			return
		}
		log.Println("Rsp:", out)
	}
	{
		log.Println("Call.LvUp 升级接口,Lv查表获得Title...")
		out := &pb.LvUp_Rsp{}
		in := &pb.LvUp_Req{
			ActorID: actorid,
		}
		if err := CallGrpcGate(stream, method_lvup, m, in, out); err != nil {
			log.Printf("CallGrpcGate Fail %s", err.Error())
			return
		}
		log.Println("Rsp:", out)
	}
	{
		log.Println("Call.ModifyName 改名接口...")
		out := &pb.ModifyName_Rsp{}
		in := &pb.ModifyName_Req{
			ActorID: actorid,
			NewName: "Mike",
		}
		if err := CallGrpcGate(stream, method_modifyname, m, in, out); err != nil {
			log.Printf("CallGrpcGate Fail %s", err.Error())
			return
		}
		log.Println("Rsp:", out)
	}
	{
		log.Println("Call.GetInfo 获取角色信息...")
		out := &pb.GetInfo_Rsp{}
		in := &pb.GetInfo_Req{
			ActorID: actorid,
		}
		if err := CallGrpcGate(stream, method_getinfo, m, in, out); err != nil {
			log.Printf("CallGrpcGate Fail %s", err.Error())
			return
		}
		log.Println("Rsp:", out)
	}
	{
		log.Println("WaitRecv 等待服务端主动推送下线通知,大概1分钟...")
		if msg, err := stream.Recv(); err != nil {
			log.Println(err)
		} else {
			log.Println("Rsp:", string(msg.Content))
		}
	}
	stream.CloseSend()

}

func TActorBroadCast() {
	log.Println("TActorBroadCast Actor广播用例...")
	m := make(map[string]string)
	m["m3routetype"] = "RouteTypeRandom"
	conn, err := grpc.Dial(
		_agenturl,
		grpc.WithInsecure())
	if err != nil {
		panic(err.Error())
	}
	cli := grpcgate.NewGateSerClient(conn)
	stream, err := cli.CSTransport(context.Background())
	if err != nil {
		panic(err.Error())
	}

	{
		log.Println("Call.Auth 建立连接...")
		in := &metapb.AuthReq{
			Token: "PlayerTest",
		}
		out := &metapb.AuthRsp{}
		if err := CallGrpcGate(stream, "", m, in, out); err != nil {
			log.Printf("CallGrpcGate Fail %s", err.Error())
		}
		log.Println("PlayerID:", out.PlayerID)
	}
	var actorid string
	{
		log.Println("Call.Register 注册接口,调用ActorApp的ActorRegSer...")
		out := &pb.Register_Rsp{}
		in := &pb.Register_Req{
			Name: "June",
		}
		if err := CallGrpcGate(stream, method_register, m, in, out); err != nil {
			log.Printf("CallGrpcGate Fail %s", err.Error())
			return
		}
		log.Println("Rsp-ActorID:", out.ActorID)
		actorid = out.ActorID
	}
	m["m3actoractorid"] = actorid
	{
		log.Println("Call.Login 登陆接口...")
		out := &pb.Login_Rsp{}
		in := &pb.Login_Req{
			ActorID: actorid,
		}
		if err := CallGrpcGate(stream, method_login, m, in, out); err != nil {
			log.Printf("CallGrpcGate Fail %s", err.Error())
			return
		}
		log.Println("Rsp:", out)
	}
	{
		log.Println("Call.PostChannel 发送频道消息接口, ActorSer向AsyncApp发送广播...")
		out := &pb.PostChannel_Rsp{}
		in := &pb.PostChannel_Req{
			ActorID: actorid,
			Content: "Good Morning",
		}
		if err := CallGrpcGate(stream, method_postchannel, m, in, out); err != nil {
			log.Printf("CallGrpcGate Fail %s", err.Error())
			return
		}
		log.Println("Rsp:", out)
	}
	{
		log.Println("Call.PullChannel 拉取消息接口,服务端强制sleep 1s...")
		out := &pb.PullChannel_Rsp{}
		in := &pb.PullChannel_Req{
			ActorID: actorid,
		}
		if err := CallGrpcGate(stream, method_pullchannel, m, in, out); err != nil {
			log.Printf("CallGrpcGate Fail %s", err.Error())
			return
		}
		log.Println("Rsp:", out)
	}
	stream.CloseSend()
}
