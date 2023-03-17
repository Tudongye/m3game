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
	var wg sync.WaitGroup
	var mutx sync.Mutex
	failnum := 0
	j := time.Now().Unix() % 10000 * 100
	for i := 0; i < 100; i++ {
		time.Sleep(1 * time.Second)
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
	log.Println("TActorCommon Actor常规用例...")
	m := make(map[string]string)
	conn, err := grpc.Dial(
		_agenturl,
		grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return err
	}
	cli := grpcgate.NewGateSerClient(conn)
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
	if err := FRoleGetInfo(stream, m); err != nil {
		log.Println(err)
		return err
	}
	stream.CloseSend()
	return nil
}
