package test

import (
	"flag"
	"fmt"
	"log"
	"m3game/meta/metapb"
	"m3game/plugins/gate/grpcgate"
	"strings"

	"github.com/golang/protobuf/proto"
)

const (
	TUType = 1
	Token1 = "Player1"
	Token2 = "Player2"
	Token3 = "Player3"
)

type TestFunc struct {
	Help string
	F    func()
}

const (
	method_hello       = "/proto.MultiSer/Hello"
	method_tracehello  = "/proto.MultiSer/TraceHello"
	method_breakhello  = "/proto.MultiSer/BreakHello"
	method_register    = "/proto.ActorRegSer/Register"
	method_login       = "/proto.ActorSer/Login"
	method_getinfo     = "/proto.ActorSer/GetInfo"
	method_modifyname  = "/proto.ActorSer/ModifyName"
	method_lvup        = "/proto.ActorSer/LvUp"
	method_postchannel = "/proto.ActorSer/PostChannel"
	method_pullchannel = "/proto.ActorSer/PullChannel"
)

var (
	TestList = map[string]TestFunc{
		"Hello":          {Help: "Hello RPC, HelloWorld", F: THello},
		"Trace":          {Help: "TraceHello RPC, 测试链路追踪", F: TTrace},
		"Break":          {Help: "BreakHello RPC, 测试熔断限流", F: TBreak},
		"ActorCommon":    {Help: "Actor, 测试注册，登陆，改名，升级", F: TActorCommon},
		"ActorBroadCast": {Help: "Actor, 测试注册，登陆，广播", F: TActorBroadCast},
		"ActorMove":      {Help: "Actor, 测试两个ActorSer之间进行服务迁移", F: TActorMove},
	}
	_agenturl    string
	clientserial = 1
)

func Help() string {
	var h []string
	h = append(h, "Help:")
	for t, f := range TestList {
		h = append(h, fmt.Sprintf("Test [%s]", t))
		h = append(h, fmt.Sprintf("%s", f.Help))
	}
	return strings.Join(h, "\n")
}

func Start() {
	var testmode string
	flag.StringVar(&testmode, "testmode", "", Help())
	flag.StringVar(&_agenturl, "agenturl", "", "")
	flag.Parse()
	if f, ok := TestList[testmode]; !ok {
		log.Printf(Help())
		return
	} else {
		f.F()
	}
}

func CallGrpcGate(stream grpcgate.GateSer_CSTransportClient, method string, metas map[string]string, in proto.Message, out proto.Message) error {
	log.Printf("CallGrpcGate %s\n", method)
	inbyte, err := proto.Marshal(in)
	if err != nil {
		log.Println(err)
		return err
	}
	inmsg := &metapb.CSMsg{
		Method:  method,
		Content: inbyte,
	}
	for k, v := range metas {
		inmsg.Metas = append(inmsg.Metas, &metapb.Meta{Key: k, Value: v})
	}
	curserial := fmt.Sprintf("%d", clientserial)
	inmsg.Metas = append(inmsg.Metas, &metapb.Meta{Key: "m3clientserial", Value: curserial})
	clientserial += 1
	log.Println("inmsg", inmsg.Metas)
	stream.Send(inmsg)
	var outmsg *metapb.CSMsg
	for {
		var err error
		outmsg, err = stream.Recv()
		if err != nil {
			log.Println(err)
			return err
		}
		recvserial := ""
		log.Println(outmsg.Metas)
		for _, m := range outmsg.Metas {
			if m.Key == "m3clientserial" {
				recvserial = m.Value
				break
			}
		}
		if recvserial == curserial {
			break
		} else {
			log.Printf("Recv Serial %s But %s Content %s \n", recvserial, curserial, string(outmsg.Content))
		}
	}
	if err := proto.Unmarshal(outmsg.Content, out); err != nil {
		log.Println(err)
		log.Println(string(outmsg.Content))
		return err
	}
	return nil
}
