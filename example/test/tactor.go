package test

import (
	"log"
	"m3game/example/proto"
	"m3game/example/proto/pb"
	"m3game/plugins/agent"
	mpb "m3game/proto/pb"
)

func TActorMode1() {
	log.Println("ActorMode1...")
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
	var actorid string
	{
		log.Println("Call Register...")
		var out pb.Register_Rsp
		in := &pb.Register_Req{
			RouteHead: GenRandomRouteHead(rsp.EnvID, rsp.WorldID, proto.ActorAppFuncID),
			Name:      "June",
		}
		in.RouteHead.Metas.Metas = append(in.RouteHead.Metas.Metas,
			&mpb.Meta{
				Key:   proto.RHMeta_ActorID,
				Value: actorid,
			},
		)
		if err := CallLogic(rsp.Uid, rsp.Session, method_register, in, &out); err != nil {
			log.Printf("Call Fail %v", err)
			return
		}
		log.Printf("Res: %s\n", out.ActorID)
		actorid = out.ActorID
	}
	{
		log.Println("Call Login...")
		var out pb.Login_Rsp
		in := &pb.Login_Req{
			RouteHead: GenRandomRouteHead(rsp.EnvID, rsp.WorldID, proto.ActorAppFuncID),
			ActorID:   actorid,
		}
		in.RouteHead.Metas.Metas = append(in.RouteHead.Metas.Metas,
			&mpb.Meta{
				Key:   proto.RHMeta_ActorID,
				Value: actorid,
			},
		)
		if err := CallLogic(rsp.Uid, rsp.Session, method_login, in, &out); err != nil {
			log.Printf("Call Fail %v", err)
			return
		}
		log.Printf("Res: %s\n", out.ActorDB)
	}
	{
		log.Println("Call GetInfo...")
		var out pb.GetInfo_Rsp
		in := &pb.GetInfo_Req{
			RouteHead: GenRandomRouteHead(rsp.EnvID, rsp.WorldID, proto.ActorAppFuncID),
			ActorID:   actorid,
		}
		in.RouteHead.Metas.Metas = append(in.RouteHead.Metas.Metas,
			&mpb.Meta{
				Key:   proto.RHMeta_ActorID,
				Value: actorid,
			},
		)
		if err := CallLogic(rsp.Uid, rsp.Session, method_getinfo, in, &out); err != nil {
			log.Printf("Call Fail %v", err)
			return
		}
		log.Printf("Res: %s %s\n", out.Name, out.Title)
	}
}

func TActorMode2() {
	log.Println("ActorMode2...")
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
	var actorid string
	{
		log.Println("Call Register...")
		var out pb.Register_Rsp
		in := &pb.Register_Req{
			RouteHead: GenRandomRouteHead(rsp.EnvID, rsp.WorldID, proto.ActorAppFuncID),
			Name:      "June",
		}
		in.RouteHead.Metas.Metas = append(in.RouteHead.Metas.Metas,
			&mpb.Meta{
				Key:   proto.RHMeta_ActorID,
				Value: actorid,
			},
		)
		if err := CallLogic(rsp.Uid, rsp.Session, method_register, in, &out); err != nil {
			log.Printf("Call Fail %v", err)
			return
		}
		log.Printf("Res: %s\n", out.ActorID)
		actorid = out.ActorID
	}
	{
		log.Println("Call Login...")
		var out pb.Login_Rsp
		in := &pb.Login_Req{
			RouteHead: GenRandomRouteHead(rsp.EnvID, rsp.WorldID, proto.ActorAppFuncID),
			ActorID:   actorid,
		}
		in.RouteHead.Metas.Metas = append(in.RouteHead.Metas.Metas,
			&mpb.Meta{
				Key:   proto.RHMeta_ActorID,
				Value: actorid,
			},
		)
		if err := CallLogic(rsp.Uid, rsp.Session, method_login, in, &out); err != nil {
			log.Printf("Call Fail %v", err)
			return
		}
		log.Printf("Res: %s\n", out.ActorDB)
	}
	{
		log.Println("Call ModifyName...")
		var out pb.ModifyName_Rsp
		in := &pb.ModifyName_Req{
			RouteHead: GenRandomRouteHead(rsp.EnvID, rsp.WorldID, proto.ActorAppFuncID),
			ActorID:   actorid,
			NewName:   "Mike",
		}
		in.RouteHead.Metas.Metas = append(in.RouteHead.Metas.Metas,
			&mpb.Meta{
				Key:   proto.RHMeta_ActorID,
				Value: actorid,
			},
		)
		if err := CallLogic(rsp.Uid, rsp.Session, method_modifyname, in, &out); err != nil {
			log.Printf("Call Fail %v", err)
			return
		}
		log.Printf("Res: %s\n", out.ActorName.Name)
	}
	{
		log.Println("Call GetInfo...")
		var out pb.GetInfo_Rsp
		in := &pb.GetInfo_Req{
			RouteHead: GenRandomRouteHead(rsp.EnvID, rsp.WorldID, proto.ActorAppFuncID),
			ActorID:   actorid,
		}
		in.RouteHead.Metas.Metas = append(in.RouteHead.Metas.Metas,
			&mpb.Meta{
				Key:   proto.RHMeta_ActorID,
				Value: actorid,
			},
		)
		if err := CallLogic(rsp.Uid, rsp.Session, method_getinfo, in, &out); err != nil {
			log.Printf("Call Fail %v", err)
			return
		}
		log.Printf("Res: %s %s\n", out.Name, out.Title)
	}
}

func TActorMode3() {
	log.Println("ActorMode2...")
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
	var actorid string
	{
		log.Println("Call Register...")
		var out pb.Register_Rsp
		in := &pb.Register_Req{
			RouteHead: GenRandomRouteHead(rsp.EnvID, rsp.WorldID, proto.ActorAppFuncID),
			Name:      "June",
		}
		in.RouteHead.Metas.Metas = append(in.RouteHead.Metas.Metas,
			&mpb.Meta{
				Key:   proto.RHMeta_ActorID,
				Value: actorid,
			},
		)
		if err := CallLogic(rsp.Uid, rsp.Session, method_register, in, &out); err != nil {
			log.Printf("Call Fail %v", err)
			return
		}
		log.Printf("Res: %s\n", out.ActorID)
		actorid = out.ActorID
	}
	{
		log.Println("Call Login...")
		var out pb.Login_Rsp
		in := &pb.Login_Req{
			RouteHead: GenRandomRouteHead(rsp.EnvID, rsp.WorldID, proto.ActorAppFuncID),
			ActorID:   actorid,
		}
		in.RouteHead.Metas.Metas = append(in.RouteHead.Metas.Metas,
			&mpb.Meta{
				Key:   proto.RHMeta_ActorID,
				Value: actorid,
			},
		)
		if err := CallLogic(rsp.Uid, rsp.Session, method_login, in, &out); err != nil {
			log.Printf("Call Fail %v", err)
			return
		}
		log.Printf("Res: %s\n", out.ActorDB)
	}
	{
		log.Println("Call LvUp...")
		var out pb.LvUp_Rsp
		in := &pb.LvUp_Req{
			RouteHead: GenRandomRouteHead(rsp.EnvID, rsp.WorldID, proto.ActorAppFuncID),
			ActorID:   actorid,
		}
		in.RouteHead.Metas.Metas = append(in.RouteHead.Metas.Metas,
			&mpb.Meta{
				Key:   proto.RHMeta_ActorID,
				Value: actorid,
			},
		)
		if err := CallLogic(rsp.Uid, rsp.Session, method_lvup, in, &out); err != nil {
			log.Printf("Call Fail %v", err)
			return
		}
		log.Printf("Res: %d\n", out.ActorInfo.Level)
	}
	{
		log.Println("Call GetInfo...")
		var out pb.GetInfo_Rsp
		in := &pb.GetInfo_Req{
			RouteHead: GenRandomRouteHead(rsp.EnvID, rsp.WorldID, proto.ActorAppFuncID),
			ActorID:   actorid,
		}
		in.RouteHead.Metas.Metas = append(in.RouteHead.Metas.Metas,
			&mpb.Meta{
				Key:   proto.RHMeta_ActorID,
				Value: actorid,
			},
		)
		if err := CallLogic(rsp.Uid, rsp.Session, method_getinfo, in, &out); err != nil {
			log.Printf("Call Fail %v", err)
			return
		}
		log.Printf("Res: %s %s\n", out.Name, out.Title)
	}
}

func TActorMode4() {
	log.Println("ActorMode2...")
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
	var actorid string
	{
		log.Println("Call Register...")
		var out pb.Register_Rsp
		in := &pb.Register_Req{
			RouteHead: GenRandomRouteHead(rsp.EnvID, rsp.WorldID, proto.ActorAppFuncID),
			Name:      "June",
		}
		in.RouteHead.Metas.Metas = append(in.RouteHead.Metas.Metas,
			&mpb.Meta{
				Key:   proto.RHMeta_ActorID,
				Value: actorid,
			},
		)
		if err := CallLogic(rsp.Uid, rsp.Session, method_register, in, &out); err != nil {
			log.Printf("Call Fail %v", err)
			return
		}
		log.Printf("Res: %s\n", out.ActorID)
		actorid = out.ActorID
	}
	{
		log.Println("Call Login...")
		var out pb.Login_Rsp
		in := &pb.Login_Req{
			RouteHead: GenRandomRouteHead(rsp.EnvID, rsp.WorldID, proto.ActorAppFuncID),
			ActorID:   actorid,
		}
		in.RouteHead.Metas.Metas = append(in.RouteHead.Metas.Metas,
			&mpb.Meta{
				Key:   proto.RHMeta_ActorID,
				Value: actorid,
			},
		)
		if err := CallLogic(rsp.Uid, rsp.Session, method_login, in, &out); err != nil {
			log.Printf("Call Fail %v", err)
			return
		}
		log.Printf("Res: %s\n", out.ActorDB)
	}
	{
		log.Println("Call PostChannel...")
		var out pb.PostChannel_Rsp
		in := &pb.PostChannel_Req{
			RouteHead: GenRandomRouteHead(rsp.EnvID, rsp.WorldID, proto.ActorAppFuncID),
			ActorID:   actorid,
			Content:   "Good Morning",
		}
		in.RouteHead.Metas.Metas = append(in.RouteHead.Metas.Metas,
			&mpb.Meta{
				Key:   proto.RHMeta_ActorID,
				Value: actorid,
			},
		)
		if err := CallLogic(rsp.Uid, rsp.Session, method_postchannel, in, &out); err != nil {
			log.Printf("Call Fail %v", err)
			return
		}
		log.Printf("Res: \n")
	}
	{
		log.Println("Call PullChannel...")
		var out pb.PullChannel_Rsp
		in := &pb.PullChannel_Req{
			RouteHead: GenRandomRouteHead(rsp.EnvID, rsp.WorldID, proto.ActorAppFuncID),
			ActorID:   actorid,
		}
		in.RouteHead.Metas.Metas = append(in.RouteHead.Metas.Metas,
			&mpb.Meta{
				Key:   proto.RHMeta_ActorID,
				Value: actorid,
			},
		)
		if err := CallLogic(rsp.Uid, rsp.Session, method_pullchannel, in, &out); err != nil {
			log.Printf("Call Fail %v", err)
			return
		}
		log.Printf("Res: %v\n", out.Msgs)
	}
}
