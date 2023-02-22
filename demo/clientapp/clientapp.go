package clientapp

import (
	"context"
	"log"
	"m3game/app"
	"m3game/config"
	"m3game/demo/dirapp/dirclient"
	"m3game/demo/mapapp/mapclient"
	"m3game/demo/roleapp/roleclient"
	"m3game/mesh/router/consul"
	"m3game/proto"
	"m3game/runtime"
	"m3game/runtime/plugin"
	"m3game/server"
	"time"
)

const (
	AppFuncID = "client"
)

func CreateClientApp() *ClientApp {
	return &ClientApp{
		DefaultApp: app.CreateDefaultApp(AppFuncID),
	}
}

type ClientApp struct {
	*app.DefaultApp
}

func (d *ClientApp) Start() error {
	router := plugin.GetRouterPlugin()
	if router != nil {
		if err := router.Register(d); err != nil {
			return err
		}
	}
	if err := dirclient.Init(d.RouteIns(), func(c *dirclient.Client) { c.Client = proto.META_FLAG_TRUE }); err != nil {
		return err
	}
	if err := mapclient.Init(d.RouteIns(), func(c *mapclient.Client) { c.Client = proto.META_FLAG_TRUE }); err != nil {
		return err
	}
	if err := roleclient.Init(d.RouteIns(), func(c *roleclient.Client) { c.Client = proto.META_FLAG_TRUE }); err != nil {
		return err
	}
	testmode := config.GetEnv("testmode")
	if testmode == "dirapp" {
		log.Printf("Test:%s\n", testmode)
		go func() {
			time.Sleep(time.Second * 3)
			log.Println("Call Hello()")
			log.Printf("Req: good morning\n")
			if rsp, err := dirclient.DirClient().Hello(context.Background(), "good morning"); err != nil {
				log.Printf("Err: %s\n", err.Error())
			} else {
				log.Printf("Res: %s\n", rsp)
			}
		}()
	}
	if testmode == "mapapp" {
		log.Printf("Test:%s\n", testmode)
		go func() {
			time.Sleep(time.Second * 3)
			log.Println("G1 Call Move()")
			log.Printf("G1 Req: Mike 5")
			if n, l, err := mapclient.MapClient().Move(context.Background(), "Mike", 5); err != nil {
				log.Printf("G1 Err: %s\n", err.Error())
			} else {
				log.Printf("G1 Res: %s %d\n", n, l)
			}
		}()

		go func() {
			time.Sleep(time.Second * 3)
			log.Println("G2 Call Move()")
			log.Printf("G2 Req: June 10")
			if n, l, err := mapclient.MapClient().Move(context.Background(), "June", 10); err != nil {
				log.Printf("G2 Err: %s\n", err.Error())
			} else {
				log.Printf("G2 Res: %s %d\n", n, l)
			}
		}()
	}
	if testmode == "roleapp" {
		log.Printf("Test:%s\n", testmode)
		time.Sleep(time.Second * 3)
		actorid := "ABCDEFG1234567"
		name := "Mike"

		log.Println("Call Register()")
		log.Printf("Req: %s %s", actorid, name)
		if roleid, err := roleclient.RoleClient().Register(context.Background(), actorid, name); err != nil {
			log.Printf("Err: %s\n", err.Error())
		} else {
			log.Printf("Res: %s\n", roleid)
		}

		log.Println("Call Login()")
		log.Printf("Req: %s ", actorid)
		if name, tips, err := roleclient.RoleClient().Login(context.Background(), actorid); err != nil {
			log.Printf("Err: %s\n", err.Error())
		} else {
			log.Printf("Res: %s %s\n", name, tips)
		}

		log.Println("Call GetName()")
		log.Printf("Req: %s", actorid)
		if name, err := roleclient.RoleClient().GetName(context.Background(), actorid); err != nil {
			log.Printf("Err: %s\n", err.Error())
		} else {
			log.Printf("Res: %s\n", name)
		}

		newname := "June"
		log.Println("Call ModifyName()")
		log.Printf("Req: %s %s", actorid, newname)
		if name, err := roleclient.RoleClient().ModifyName(context.Background(), actorid, newname); err != nil {
			log.Printf("Err: %s\n", err.Error())
		} else {
			log.Printf("Res: %s\n", name)
		}

		log.Println("Call GetName()")
		log.Printf("Req: %s", actorid)
		if name, err := roleclient.RoleClient().GetName(context.Background(), actorid); err != nil {
			log.Printf("Err: %s\n", err.Error())
		} else {
			log.Printf("Res: %s\n", name)
		}
		go func() {
			log.Println("G1 Call Move()")
			log.Printf("G1 Req: %s %d", actorid, 15)
			if location, locatename, err := roleclient.RoleClient().MoveRole(context.Background(), actorid, 15); err != nil {
				log.Printf("G1 Err: %s\n", err.Error())
			} else {
				log.Printf("G1 Res: %d %s\n", location, locatename)
			}
		}()
		go func() {
			log.Println("G2 Call Move()")
			log.Printf("G2 Req: %s %d", actorid, 10)
			if location, locatename, err := roleclient.RoleClient().MoveRole(context.Background(), actorid, 10); err != nil {
				log.Printf("G2 Err: %s\n", err.Error())
			} else {
				log.Printf("G2 Res: %d %s\n", location, locatename)
			}
		}()
	}
	return nil
}

func (d *ClientApp) HealthCheck() bool {
	return true
}

func Run() error {
	plugin.RegisterPluginFactory(&consul.Factory{})
	runtime.Run(CreateClientApp(), []server.Server{})
	return nil
}
