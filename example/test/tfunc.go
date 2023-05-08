package test

import (
	"log"
	"m3game/example/proto/pb"
	"m3game/plugins/gate/grpcgate"
)

func FWaitRecv(stream grpcgate.GGateSer_CSTransportClient) error {
	log.Println("WaitRecv 等待服务端主动推送...")
	if msg, err := stream.Recv(); err != nil {
		log.Println(err)
		return err
	} else {
		log.Println("Rsp:", string(msg.Content))
	}
	return nil
}

func FAuth(stream grpcgate.GGateSer_CSTransportClient, m map[string]string, token string) (string, error) {
	log.Println("Call.Auth 建立连接...")
	in := &pb.AuthReq{
		Token: token,
	}
	out := &pb.AuthRsp{}
	if err := CallGrpcGate(stream, "", m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return "", err
	}
	log.Println("RoleId:", out.PlayerId)
	return out.PlayerId, nil
}

func FRegister(stream grpcgate.GGateSer_CSTransportClient, m map[string]string, playerid string, name string) error {
	log.Println("Call.Register 注册...")
	in := &pb.Register_Req{
		Name:     name,
		PlayerID: playerid,
	}
	out := &pb.Register_Rsp{}
	if err := CallGrpcGate(stream, method_register, m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return err
	}
	log.Println("Rsp:", out)
	return nil
}

func FLogin(stream grpcgate.GGateSer_CSTransportClient, m map[string]string, playerid string) error {
	log.Println("Call.Login 登陆...")
	in := &pb.Login_Req{
		ActorID: playerid,
	}
	out := &pb.Login_Rsp{}
	if err := CallGrpcGate(stream, method_login, m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return err
	}
	log.Println("ActorDB:", out.ActorDB)
	return nil
}

func FPostChannel(stream grpcgate.GGateSer_CSTransportClient, m map[string]string, playerid string, content string) error {
	log.Println("Call.PostChannel 发送消息...")
	in := &pb.PostChannel_Req{
		ActorID: playerid,
		Content: content,
	}
	out := &pb.PostChannel_Rsp{}
	if err := CallGrpcGate(stream, method_postchannel, m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return err
	}
	log.Println("Rsp:", out)
	return nil
}

func FPullChannel(stream grpcgate.GGateSer_CSTransportClient, m map[string]string, playerid string) error {
	log.Println("Call.PullChannel 接受消息...")
	in := &pb.PullChannel_Req{
		ActorID: playerid,
	}
	out := &pb.PullChannel_Rsp{}
	if err := CallGrpcGate(stream, method_pullchannel, m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return err
	}
	log.Println("Rsp:", out)
	return nil
}

func FTraceHello(stream grpcgate.GGateSer_CSTransportClient, m map[string]string, content string) error {
	log.Println("Call.TraceHello 链路追踪...")
	in := &pb.TraceHello_Req{
		Req: content,
	}
	out := &pb.TraceHello_Rsp{}
	if err := CallGrpcGate(stream, method_tracehello, m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return err
	}
	log.Println("Rsp:", out)
	return nil
}

func FBreakHello(stream grpcgate.GGateSer_CSTransportClient, m map[string]string, content string) error {
	log.Println("Call.BreakHello 熔断...")
	in := &pb.BreakHello_Req{
		Req: content,
	}
	out := &pb.BreakHello_Rsp{}
	if err := CallGrpcGate(stream, method_breakhello, m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return err
	}
	log.Println("Rsp:", out)
	return nil
}

func FHello(stream grpcgate.GGateSer_CSTransportClient, m map[string]string, content string) error {
	log.Println("Call.Hello Hello ...")
	in := &pb.Hello_Req{
		Req: content,
	}
	out := &pb.Hello_Rsp{}
	if err := CallGrpcGate(stream, method_hello, m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return err
	}
	log.Println("Rsp:", out)
	return nil
}

func FGetInfo(stream grpcgate.GGateSer_CSTransportClient, m map[string]string, playerid string) error {
	log.Println("Call.GetInfo 获取角色信息...")
	in := &pb.GetInfo_Req{
		ActorID: playerid,
	}
	out := &pb.GetInfo_Rsp{}
	if err := CallGrpcGate(stream, method_getinfo, m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return err
	}
	log.Println("Rsp:", out)
	return nil
}

func FLvUp(stream grpcgate.GGateSer_CSTransportClient, m map[string]string, playerid string) error {
	log.Println("Call.LvUp 升级...")
	in := &pb.LvUp_Req{
		ActorID: playerid,
	}
	out := &pb.LvUp_Rsp{}
	if err := CallGrpcGate(stream, method_lvup, m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return err
	}
	log.Println("Rsp:", out)
	return nil
}

func FModifyName(stream grpcgate.GGateSer_CSTransportClient, m map[string]string, playerid string, name string) error {
	log.Println("Call.ModifyName 修改角色名...")
	in := &pb.ModifyName_Req{
		ActorID: playerid,
		NewName: name,
	}
	out := &pb.ModifyName_Rsp{}
	if err := CallGrpcGate(stream, method_modifyname, m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return err
	}
	log.Println("Rsp:", out)
	return nil
}

func FCreateEntity(stream grpcgate.GGateSer_CSTransportClient, m map[string]string, name string, pos *pb.Position) error {
	log.Println("Call.CreateEntity 创建对象...")
	in := &pb.CreateEntity_Req{
		Name:   name,
		SrcPos: pos,
	}
	out := &pb.CreateEntity_Rsp{}
	if err := CallGrpcGate(stream, pb.WorldSer_CreateEntity_FullMethodName, m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return err
	}
	log.Println("Rsp:", out)
	return nil
}

func FMoveEntity(stream grpcgate.GGateSer_CSTransportClient, m map[string]string, name string, pos *pb.Position) error {
	log.Println("Call.MoveEntity 移动对象...")
	in := &pb.MoveEntity_Req{
		Name:   name,
		DstPos: pos,
	}
	out := &pb.MoveEntity_Rsp{}
	if err := CallGrpcGate(stream, pb.WorldSer_MoveEntity_FullMethodName, m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return err
	}
	log.Println("Rsp:", out)
	return nil
}

func FViewPosition(stream grpcgate.GGateSer_CSTransportClient, m map[string]string, pos *pb.Position) error {
	log.Println("Call.ViewPosition 观察点位...")
	in := &pb.ViewPosition_Req{
		Pos: pos,
	}
	out := &pb.ViewPosition_Rsp{}
	if err := CallGrpcGate(stream, pb.WorldSer_ViewPosition_FullMethodName, m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return err
	}
	log.Println("Rsp:", out)
	return nil
}
