package test

import (
	"log"
	"m3game/demo/proto/pb"
	"m3game/plugins/gate/grpcgate"
)

func FWaitRecv(stream grpcgate.GateSer_CSTransportClient) error {
	log.Println("WaitRecv 等待服务端主动推送...")
	if msg, err := stream.Recv(); err != nil {
		log.Println(err)
		return err
	} else {
		log.Println("Rsp:", string(msg.Content))
	}
	return nil
}

func FAuth(stream grpcgate.GateSer_CSTransportClient, m map[string]string, token string) error {
	log.Println("Call.Auth 建立连接...")
	in := &pb.AuthReq{
		Token: token,
	}
	out := &pb.AuthRsp{}
	if err := CallGrpcGate(stream, "", m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return err
	}
	log.Println("RoleId:", out.RoleId)
	return nil
}

func FRoleLogin(stream grpcgate.GateSer_CSTransportClient, m map[string]string) error {
	log.Println("Call.RoleLogin 登陆...")
	in := &pb.RoleLogin_Req{}
	out := &pb.RoleLogin_Rsp{}
	if err := CallGrpcGate(stream, method_RoleLogin, m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return err
	}
	log.Println("RoleDB:", out.RoleDB)
	log.Println("ClubRoleDB:", out.ClubRoleDB)
	return nil
}

func FRoleGetInfo(stream grpcgate.GateSer_CSTransportClient, m map[string]string) error {
	log.Println("Call.RoleGetInfo 获取详情...")
	in := &pb.RoleGetInfo_Req{}
	out := &pb.RoleGetInfo_Rsp{}
	if err := CallGrpcGate(stream, method_RoleGetInfo, m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return err
	}
	log.Println("RoleDB:", out.RoleDB)
	log.Println("ClubRoleDB:", out.ClubRoleDB)
	return nil
}

func FRoleModifyName(stream grpcgate.GateSer_CSTransportClient, m map[string]string, name string) error {
	log.Println("Call.RoleModifyName 改名...", name)
	in := &pb.RoleModifyName_Req{
		NewName: name,
	}
	out := &pb.RoleModifyName_Rsp{}
	if err := CallGrpcGate(stream, method_RoleModifyName, m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return err
	}
	log.Println("RoleName:", out.Name)
	return nil
}

func FRolePowerUp(stream grpcgate.GateSer_CSTransportClient, m map[string]string, up int) error {
	log.Println("Call.RolePowerUp 火力提升...", up)
	in := &pb.RolePowerUp_Req{
		PowerUp: int32(up),
	}
	out := &pb.RolePowerUp_Rsp{}
	if err := CallGrpcGate(stream, method_RolePowerUp, m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return err
	}
	log.Println("RolePower:", out.Power)
	return nil
}
