package test

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"m3game/meta/metapb"
	"m3game/plugins/gate/grpcgate"
	"net/http"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
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
	}
	_agenturl string
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

func CallLogic(method string, req interface{}, rsp interface{}, m map[string]string) error {
	log.Printf("CallLogic %s,", method)
	client := &http.Client{}
	bytesData, _ := json.Marshal(req)
	if req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", _agenturl, method), bytes.NewReader(bytesData)); err != nil {
		return errors.Wrap(err, "NewRequest")
	} else {
		for k, v := range m {
			req.Header.Add(k, v)
		}
		if r, err := client.Do(req); err != nil {
			return errors.Wrap(err, "client.Do")
		} else {
			rc, _ := io.ReadAll(r.Body)
			if r.StatusCode != 200 {
				return fmt.Errorf("RspCode %d %s", r.StatusCode, rc)
			}
			if err := json.Unmarshal(rc, rsp); err != nil {
				return errors.Wrapf(err, "Decode [%s]", rc)
			}
			return nil
		}
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
	stream.Send(inmsg)
	outmsg, err := stream.Recv()
	if err != nil {
		log.Println(err)
		return err
	}
	if err := proto.Unmarshal(outmsg.Content, out); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
