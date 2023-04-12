package test

import (
	"context"
	"fmt"
	"log"
	"m3game/plugins/gate/grpcgate"
	"sync"
	"time"

	"google.golang.org/grpc"
)

func TMutilTest1(Token string) error {
	log.Println("TMutilTest1 Role并发登陆...")
	var wg sync.WaitGroup
	var mutx sync.Mutex
	failnum := 0
	j := time.Now().Unix() % 10000 * 10000
	for i := 0; i < 10000; i++ {
		time.Sleep(10 * time.Millisecond)
		wg.Add(1)
		go func(t string) {
			log.Printf("%s Start...\n", t)
			defer wg.Done()
			if err := TTest1(t); err != nil {
				mutx.Lock()
				defer mutx.Unlock()
				failnum += 1
			}

		}(fmt.Sprintf("Token%d", i+int(j)))
	}
	wg.Wait()
	log.Printf("FailNum %d\n", failnum)
	return nil
}

func TTest1(Token string) error {
	log.Println("TTest1 Role常规用例...")
	m := make(map[string]string)
	conn, err := grpc.Dial(
		_agenturl,
		grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return err
	}
	cli := grpcgate.NewGGateSerClient(conn)
	stream, err := cli.CSTransport(context.Background())
	if err != nil {
		log.Println(err)
		return err
	}
	// 鉴权
	if err := FAuth(stream, m, Token); err != nil {
		log.Println(err)
		return err
	}
	// 登陆
	if err := FRoleLogin(stream, m); err != nil {
		log.Println(err)
		return err
	}
	// 改名
	if err := FRoleModifyName(stream, m, "NewName"); err != nil {
		log.Println(err)
		return err
	}
	// 火力提升
	if err := FRolePowerUp(stream, m, 1000); err != nil {
		log.Println(err)
		return err
	}
	// 查询数据
	if err := FRoleGetInfo(stream, m, false); err != nil {
		log.Println(err)
		return err
	}
	// 查询详细数据
	if err := FRoleGetInfo(stream, m, true); err != nil {
		log.Println(err)
		return err
	}
	stream.CloseSend()
	return nil
}

func TTest2(Token string) error {
	log.Println("TTest2 Club常规用例...")
	m := make(map[string]string)
	conn, err := grpc.Dial(
		_agenturl,
		grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return err
	}
	cli := grpcgate.NewGGateSerClient(conn)
	stream, err := cli.CSTransport(context.Background())
	if err != nil {
		log.Println(err)
		return err
	}
	// 鉴权
	if err := FAuth(stream, m, Token); err != nil {
		log.Println(err)
		return err
	}
	// 登陆
	if err := FRoleLogin(stream, m); err != nil {
		log.Println(err)
		return err
	}
	// 查询数据
	if err := FRoleGetInfo(stream, m, false); err != nil {
		log.Println(err)
		return err
	}
	// 创建社团
	var clubid int64
	if c, err := FRoleCreateClub(stream, m); err != nil {
		log.Println(err)
		return err
	} else {
		clubid = c
	}
	// 查询数据
	if err := FRoleGetInfo(stream, m, false); err != nil {
		log.Println(err)
		return err
	}
	// 查询社团信息
	if err := FRoleGetClubInfo(stream, m, clubid); err != nil {
		log.Println(err)
		return err
	}
	// 解散社团
	if err := FRoleCancelClub(stream, m); err != nil {
		log.Println(err)
		return err
	}
	// 查询数据
	if err := FRoleGetInfo(stream, m, false); err != nil {
		log.Println(err)
		return err
	}
	stream.CloseSend()
	return nil
}
