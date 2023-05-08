package test

import (
	"context"
	"log"
	"m3game/example/proto/pb"
	"m3game/plugins/gate/grpcgate"
	"time"

	"google.golang.org/grpc"
)

func TWorld() error {
	log.Println("TWorld ...")
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

	Token := "PlayerTest"
	// 建立连接
	if _, err := FAuth(stream, m, Token); err != nil {
		log.Println(err)
		return err
	}

	// 创建实体
	if err := FCreateEntity(stream, m, "Guanyu", &pb.Position{X: 150, Y: 150}); err != nil {
		log.Println(err)
		return err
	}
	// 创建实体
	if err := FCreateEntity(stream, m, "Qinqiong", &pb.Position{X: 180, Y: 180}); err != nil {
		log.Println(err)
		return err
	}
	<-time.After(time.Second)
	// 查看点位
	if err := FViewPosition(stream, m, &pb.Position{X: 150, Y: 150}); err != nil {
		log.Println(err)
		return err
	}

	// 移动实体
	if err := FMoveEntity(stream, m, "Guanyu", &pb.Position{X: 180, Y: 180}); err != nil {
		log.Println(err)
		return err
	}
	// 移动实体
	if err := FMoveEntity(stream, m, "Qinqiong", &pb.Position{X: 150, Y: 150}); err != nil {
		log.Println(err)
		return err
	}
	for i := 0; i < 20; i++ {
		<-time.After(time.Second)
		// 查看点位
		if err := FViewPosition(stream, m, &pb.Position{X: 180, Y: 180}); err != nil {
			log.Println(err)
			return err
		}
	}
	stream.CloseSend()
	return nil
}
