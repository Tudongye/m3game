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

func FRoleGetClubInfo(stream grpcgate.GateSer_CSTransportClient, m map[string]string, clubid int64) error {
	log.Println("Call.RoleGetClubInfo 获取社团信息...", clubid)
	in := &pb.RoleGetClubInfo_Req{
		ClubId: clubid,
	}
	out := &pb.RoleGetClubInfo_Rsp{}
	if err := CallGrpcGate(stream, method_RoleGetClubInfo, m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return err
	}
	log.Println("ClubDB:", out.ClubDB)
	return nil
}

func FRoleCreateClub(stream grpcgate.GateSer_CSTransportClient, m map[string]string) (int64, error) {
	log.Println("Call.RoleCreateClub 创建社团...")
	in := &pb.RoleCreateClub_Req{}
	out := &pb.RoleCreateClub_Rsp{}
	if err := CallGrpcGate(stream, method_RoleCreateClub, m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return 0, err
	}
	log.Println("ClubId:", out.ClubId)
	return out.ClubId, nil
}

func FRoleJoinClub(stream grpcgate.GateSer_CSTransportClient, m map[string]string, clubid int64) error {
	log.Println("Call.RoleJoinClub 加入社团...", clubid)
	in := &pb.RoleJoinClub_Req{
		ClubId: clubid,
	}
	out := &pb.RoleJoinClub_Rsp{}
	if err := CallGrpcGate(stream, method_RoleJoinClub, m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return err
	}
	log.Println("RolePower:")
	return nil
}

func FRoleExitClub(stream grpcgate.GateSer_CSTransportClient, m map[string]string) error {
	log.Println("Call.RoleExitClub 退出社团...")
	in := &pb.RoleExitClub_Req{}
	out := &pb.RoleExitClub_Rsp{}
	if err := CallGrpcGate(stream, method_RoleExitClub, m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return err
	}
	log.Println("RolePower:")
	return nil
}

func FRoleCancelClub(stream grpcgate.GateSer_CSTransportClient, m map[string]string) error {
	log.Println("Call.RoleCancelClub 解散社团...")
	in := &pb.RoleCancelClub_Req{}
	out := &pb.RoleCancelClub_Rsp{}
	if err := CallGrpcGate(stream, method_RoleCancelClub, m, in, out); err != nil {
		log.Printf("CallGrpcGate Fail %s", err.Error())
		return err
	}
	log.Println("RolePower:")
	return nil
}
