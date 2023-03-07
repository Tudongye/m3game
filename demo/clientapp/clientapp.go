package clientapp

import (
	"context"
	_ "m3game/broker/nats"
	"m3game/config"
	"m3game/demo/dirapp/dirclient"
	"m3game/demo/mapapp/mapclient"
	dproto "m3game/demo/proto"
	"m3game/demo/roleapp/roleclient"
	"m3game/log"
	_ "m3game/mesh/router/consul"
	"m3game/proto"
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/client"
	"m3game/runtime/plugin"
	"m3game/runtime/server"
	_ "m3game/shape/sentinel"
	"sync"
	"time"
)

func newApp() *ClientApp {
	return &ClientApp{
		App: app.New(dproto.ClientAppFuncID),
	}
}

type ClientApp struct {
	app.App
}

func (d *ClientApp) Start(wg *sync.WaitGroup) error {
	router := plugin.GetRouterPlugin()
	if router != nil {
		if err := router.Register(d); err != nil {
			return err
		}
	}
	if err := dirclient.Init(d.RouteIns(), client.GenMetaClientOption(proto.META_FLAG_TRUE)); err != nil {
		return err
	}
	if err := mapclient.Init(d.RouteIns(), client.GenMetaClientOption(proto.META_FLAG_TRUE)); err != nil {
		return err
	}
	if err := roleclient.Init(d.RouteIns(), client.GenMetaClientOption(proto.META_FLAG_TRUE)); err != nil {
		return err
	}
	testmode := config.GetEnv("testmode")
	log.Info("TestMode:%s", testmode)
	if testmode == "dirapp" {
		log.Info("Test:%s", testmode)
		go func() {
			time.Sleep(time.Second * 3)
			log.Info("Call Hello()")
			log.Debug("Req: good morning")
			if rsp, err := dirclient.Hello(context.Background(), "good morning"); err != nil {
				log.Error("Err: %s", err.Error())
			} else {
				log.Debug("Res: %s", rsp)
			}
		}()

		for i := 0; i < 100; i++ {
			log.Info("Call Hello() %d", i)
			log.Debug("Req: good morning")
			if rsp, err := dirclient.Hello(context.Background(), "good morning"); err != nil {
				log.Error("Err: %s", err.Error())
			} else {
				log.Debug("Res: %s", rsp)
			}
		}
	}
	if testmode == "mapapp" {
		log.Info("Test:%s", testmode)
		go func() {
			time.Sleep(time.Second * 3)
			log.Info("G1 Call Move()")
			log.Debug("G1 Req: Mike 5")
			if n, l, err := mapclient.Move(context.Background(), "Mike", 5); err != nil {
				log.Error("G1 Err: %s", err.Error())
			} else {
				log.Debug("G1 Res: %s %d", n, l)
			}
		}()

		go func() {
			time.Sleep(time.Second * 3)
			log.Info("G2 Call Move()")
			log.Debug("G2 Req: June 10")
			if n, l, err := mapclient.Move(context.Background(), "June", 10); err != nil {
				log.Error("G2 Err: %s", err.Error())
			} else {
				log.Debug("G2 Res: %s %d", n, l)
			}
		}()
	}
	if testmode == "broad" {
		log.Info("Test:%s", testmode)
		time.Sleep(time.Second * 3)
		actorid := "ABCDEFG1234567"
		name := "Mike"
		log.Info("Call Register()")
		log.Debug("Req: %s %s", actorid, name)
		if roleid, err := roleclient.Register(context.Background(), actorid, name); err != nil {
			log.Error("Err: %s", err.Error())
		} else {
			log.Debug("Res: %s", roleid)
		}

		log.Info("Call Login()")
		log.Debug("Req: %s ", actorid)
		if name, tips, err := roleclient.Login(context.Background(), actorid); err != nil {
			log.Error("Err: %s", err.Error())
		} else {
			log.Debug("Res: %s %s", name, tips)
		}

		log.Info("Call PostChannel()")
		log.Debug("Req: %s %s", actorid, "Hello World")
		if err := roleclient.PostChannel(context.Background(), actorid, "Hello World"); err != nil {
			log.Error("Err: %s", err.Error())
		} else {
			log.Debug("Res: ")
		}

		log.Info("Call PostChannel()")
		log.Debug("Req: %s %s", actorid, "Hello World2")
		if err := roleclient.PostChannel(context.Background(), actorid, "Hello World2"); err != nil {
			log.Error("Err: %s", err.Error())
		} else {
			log.Debug("Res: ")
		}

		time.Sleep(time.Second * 3)
		log.Info("Call PullChannel()")
		log.Debug("Req: %s", actorid)
		if msgs, err := roleclient.PullChannel(context.Background(), actorid); err != nil {
			log.Error("Err: %s", err.Error())
		} else {
			log.Debug("Res: %v", msgs)
		}
	}
	if testmode == "roleapp" {
		log.Info("Test:%s", testmode)
		time.Sleep(time.Second * 3)
		actorid := "ABCDEFG1234567"
		name := "Mike"

		log.Info("Call Register()")
		log.Debug("Req: %s %s", actorid, name)
		if roleid, err := roleclient.Register(context.Background(), actorid, name); err != nil {
			log.Error("Err: %s", err.Error())
		} else {
			log.Debug("Res: %s", roleid)
		}

		log.Info("Call Login()")
		log.Debug("Req: %s ", actorid)
		if name, tips, err := roleclient.Login(context.Background(), actorid); err != nil {
			log.Error("Err: %s", err.Error())
		} else {
			log.Debug("Res: %s %s", name, tips)
		}

		log.Info("Call GetName()")
		log.Debug("Req: %s", actorid)
		if name, err := roleclient.GetName(context.Background(), actorid); err != nil {
			log.Error("Err: %s", err.Error())
		} else {
			log.Debug("Res: %s", name)
		}

		newname := "June"
		log.Info("Call ModifyName()")
		log.Debug("Req: %s %s", actorid, newname)
		if name, err := roleclient.ModifyName(context.Background(), actorid, newname); err != nil {
			log.Error("Err: %s", err.Error())
		} else {
			log.Debug("Res: %s", name)
		}

		log.Info("Call GetName()")
		log.Debug("Req: %s", actorid)
		if name, err := roleclient.GetName(context.Background(), actorid); err != nil {
			log.Error("Err: %s", err.Error())
		} else {
			log.Debug("Res: %s", name)
		}

		time.Sleep(time.Second * 20)
		go func() {
			log.Info("G1 Call Move()")
			log.Debug("G1 Req: %s %d", actorid, 10)
			if location, locatename, err := roleclient.MoveRole(context.Background(), actorid, 15); err != nil {
				log.Error("G1 Err: %s", err.Error())
			} else {
				log.Debug("G1 Res: %d %s", location, locatename)
			}
		}()
		go func() {
			log.Info("G2 Call Move()")
			log.Debug("G2 Req: %s %d", actorid, 5)
			if location, locatename, err := roleclient.MoveRole(context.Background(), actorid, 10); err != nil {
				log.Error("G2 Err: %s", err.Error())
			} else {
				log.Debug("G2 Res: %d %s", location, locatename)
			}
		}()
	}
	return nil
}

func (d *ClientApp) HealthCheck() bool {
	return true
}

func Run() error {
	runtime.Run(newApp(), []server.Server{})
	return nil
}
