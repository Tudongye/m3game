package rpc

import (
	"fmt"
	"m3game/meta/metapb"
	"m3game/plugins/log"
	"sync"

	"google.golang.org/protobuf/proto"

	mapset "github.com/deckarep/golang-set/v2"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	RPCClientMethods = mapset.NewSet[string]()
	RPCMetas         sync.Map
)

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

func InjectionRPC(serviced protoreflect.ServiceDescriptor) error {
	for i := 0; i < serviced.Methods().Len(); i++ {
		methodd := serviced.Methods().Get(i)
		inputd := methodd.Input()
		inputname := inputd.FullName()
		rpcde := inputd.Parent()
		rpcname := rpcde.Name()
		servicefullname := serviced.FullName()
		methodname := methodd.Name()
		if _, ok := RPCMetas.Load(inputname); ok {
			continue
		}
		meta := &RPCMeta{
			rpcname:    string(rpcname),
			grpcoption: nil,
			hashkeyd:   nil,
		}
		// eache Rpc must have rpc_option
		if v := proto.GetExtension(rpcde.Options(), metapb.E_RpcOption); v == nil {
			return fmt.Errorf("RPC %s not have E_RpcOption", rpcname)
		} else if m3grpcopt, ok := v.(*metapb.M3GRPCOption); !ok {
			return fmt.Errorf("RPC %s E_RpcOption type err", rpcname)
		} else if m3grpcopt == nil {
			return fmt.Errorf("RPC %s E_RpcOption is nil", rpcname)
		} else {
			meta.grpcoption = m3grpcopt
		}
		if meta.grpcoption.Cs {
			RPCClientMethods.Add(fmt.Sprintf("/%s/%s", servicefullname, methodname))
		}
		if fieldd := inputd.Fields().ByName(protoreflect.Name(meta.grpcoption.RouteKey)); fieldd != nil {
			meta.hashkeyd = fieldd
		}
		if _, ok := RPCMetas.LoadOrStore(inputname, meta); ok {
			continue
		}
		log.Info("RPC Registor: Svc => %s, Method => %s, Input => %s, RPC => %s", servicefullname, methodname, inputname, rpcname)
	}
	return nil
}

func Meta(inputname protoreflect.FullName) *RPCMeta {
	if v, ok := RPCMetas.Load(inputname); !ok {
		return nil
	} else {
		return v.(*RPCMeta)
	}
}

func IsRPCClientMethod(method string) bool {
	return RPCClientMethods.Contains(method)
}
