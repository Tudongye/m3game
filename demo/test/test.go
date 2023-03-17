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
)

type TestFunc struct {
	Help string
	F    func(string) error
}

const (
	method_RoleLogin      = "/proto.RoleSer/RoleLogin"
	method_RoleGetInfo    = "/proto.RoleSer/RoleGetInfo"
	method_RoleModifyName = "/proto.RoleSer/RoleModifyName"
	method_RolePowerUp    = "/proto.RoleSer/RolePowerUp"
)

var (
	TestList = map[string]TestFunc{
		"Test1":      {Help: "登录，改名，提升", F: TTest1},
		"MutilTest1": {Help: "登录，改名，提升", F: TMutilTest1},
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
		f.F("Token")
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
		log.Println(err, string(outmsg.Content))
		return err
	}
	return nil
}
