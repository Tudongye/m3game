package rpc

import (
	"fmt"
	"m3game/meta/metapb"
	"m3game/plugins/log"

	"google.golang.org/protobuf/proto"

	"google.golang.org/protobuf/reflect/protoreflect"
)

var RPCMetas = make(map[protoreflect.FullName]*RPCMeta)
var RPCCS = make(map[string]string)

type RPCMeta struct {
	rpcname    string
	grpcoption *metapb.M3GRPCOption // rpc_option

	hashkeyd protoreflect.FieldDescriptor
}

func (r *RPCMeta) RpcName() string {
	return r.rpcname
}

func (r *RPCMeta) GrpcOption() *metapb.M3GRPCOption {
	return r.grpcoption
}

func (r *RPCMeta) HashKeyd() protoreflect.FieldDescriptor {
	return r.hashkeyd
}

func RegisterRPCSvc(serviced protoreflect.ServiceDescriptor) error {
	for i := 0; i < serviced.Methods().Len(); i++ {
		methodd := serviced.Methods().Get(i)
		inputd := methodd.Input()
		inputname := inputd.FullName()
		rpcde := inputd.Parent()
		RPCMetas[inputname] = &RPCMeta{
			rpcname:    string(rpcde.Name()),
			grpcoption: nil,
			hashkeyd:   nil,
		}
		// eache Rpc must have rpc_option
		if v := proto.GetExtension(rpcde.Options(), metapb.E_RpcOption); v == nil {
			panic(fmt.Sprintf("RPC %s not have E_RpcOption", rpcde.Name()))
		} else if m3grpcopt, ok := v.(*metapb.M3GRPCOption); !ok {
			panic(fmt.Sprintf("RPC %s E_RpcOption type err", rpcde.Name()))
		} else if m3grpcopt == nil {
			panic(fmt.Sprintf("RPC %s E_RpcOption is nil", rpcde.Name()))
		} else {
			RPCMetas[inputname].grpcoption = m3grpcopt
		}
		if RPCMetas[inputname].grpcoption.Cs {
			RPCCS[fmt.Sprintf("/%s/%s", serviced.FullName(), methodd.Name())] = ""
		}
		if fieldd := inputd.Fields().ByName(protoreflect.Name(RPCMetas[inputname].grpcoption.RouteKey)); fieldd != nil {
			RPCMetas[inputname].hashkeyd = fieldd
		}
		log.Info("RPC Registor: Svc => %s, Method => %s, Input => %s, RPC => %s", serviced.FullName(), methodd.Name(), inputname, rpcde.Name())
	}
	return nil
}

func Method(inputname protoreflect.FullName) *RPCMeta {
	return RPCMetas[inputname]
}

func IsCSFullMethod(method string) bool {
	_, ok := RPCCS[method]
	return ok
}
