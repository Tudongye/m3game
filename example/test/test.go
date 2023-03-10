package test

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"m3game/plugins/agent"
	"m3game/proto/pb"
	"m3game/util"
	"net/http"
	"strings"

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
	method_hello       = "/v1/example/mutilser/hello"
	method_tracehello  = "/v1/example/mutilser/tracehello"
	method_breakhello  = "/v1/example/mutilser/breakhello"
	method_register    = "/v1/example/actorser/register"
	method_login       = "/v1/example/actorser/login"
	method_getinfo     = "/v1/example/actorser/getinfo"
	method_modifyname  = "/v1/example/actorser/modifyname"
	method_lvup        = "/v1/example/actorser/lvup"
	method_postchannel = "/v1/example/actorser/postchannel"
	method_pullchannel = "/v1/example/actorser/pullchannel"
)

var (
	TestList = map[string]TestFunc{
		"Hello":      {Help: "Hello RPC, HelloWorld", F: THello},
		"Trace":      {Help: "TraceHello RPC, 测试链路追踪", F: TTrace},
		"Break":      {Help: "BreakHello RPC, 测试熔断限流", F: TBreak},
		"ActorMode1": {Help: "Actor RPC, 测试Actor注册，登陆 ", F: TActorMode1},
		"ActorMode2": {Help: "Actor, 测试注册，登陆，改名", F: TActorMode2},
		"ActorMode3": {Help: "Actor, 测试注册，登陆，升级", F: TActorMode3},
		"ActorMode4": {Help: "Actor, 测试注册，登陆，广播", F: TActorMode4},
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
		Help()
		return
	} else {
		f.F()
	}
}

func CallAuth(req *agent.AuthPara) (*agent.AuthRsp, error) {
	log.Printf("CallAuth %d, %s", req.Type, req.Token)
	client := &http.Client{}
	bytesData, _ := json.Marshal(req)
	if req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", _agenturl, "/example/auth"), bytes.NewReader(bytesData)); err != nil {
		return nil, errors.Wrap(err, "NewRequest")
	} else if r, err := client.Do(req); err != nil {
		return nil, errors.Wrap(err, "client.Do")
	} else {
		var rsp agent.AuthRsp
		if err := json.NewDecoder(r.Body).Decode(&rsp); err != nil {
			return nil, errors.Wrap(err, "Decode")
		}
		log.Print(rsp)
		return &rsp, nil
	}
}

func CallLogic(uid string, session string, method string, req interface{}, rsp interface{}) error {
	log.Printf("CallLogic %s, %s, %s", uid, session, method)
	client := &http.Client{}
	bytesData, _ := json.Marshal(req)
	if req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", _agenturl, "/example/logic"), bytes.NewReader(bytesData)); err != nil {
		return errors.Wrap(err, "NewRequest")
	} else {
		req.Header.Add("uid", uid)
		req.Header.Add("session", session)
		req.Header.Add("method", method)
		if r, err := client.Do(req); err != nil {
			return errors.Wrap(err, "client.Do")
		} else {
			if r.StatusCode != 200 {
				return fmt.Errorf("RspCode %d", r.StatusCode)
			}
			rc, _ := io.ReadAll(r.Body)
			log.Println(string(rc))
			if err := json.Unmarshal(rc, rsp); err != nil {
				return errors.Wrapf(err, "Decode [%s]", rc)
			}
			return nil
		}
	}
}

func GenRandomRouteHead(envid string, worldid string, funcid string) *pb.RouteHead {
	return &pb.RouteHead{
		DstSvc: &pb.RouteSvc{
			EnvID:   envid,
			WorldID: worldid,
			FuncID:  funcid,
			IDStr:   util.SvcID2Str(envid, worldid, funcid),
		},
		RouteType: pb.RouteType_RT_RAND,
		RoutePara: &pb.RoutePara{
			Para: &pb.RoutePara_RouteRandHead{
				RouteRandHead: &pb.RouteRandHead{
					Pass: "",
				},
			},
		},
		Metas: &pb.Metas{
			Metas: []*pb.Meta{
				{Key: "_rhmeta_client", Value: "1"},
			},
		},
	}
}
