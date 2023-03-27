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
	grpcoption *metapb.M3RpcOption // rpc_option

	hashkeyd protoreflect.FieldDescriptor
}

func (r *RPCMeta) RpcName() string {
	return r.rpcname
}

func (r *RPCMeta) GrpcOption() *metapb.M3RpcOption {
	return r.grpcoption
}

func (r *RPCMeta) HashKey(msg proto.Message) (string, error) {
	var hashkey string
	switch r.hashkeyd.Kind() {
	case protoreflect.StringKind:
		hashkey = msg.ProtoReflect().Get(r.hashkeyd).Interface().(string)
	case protoreflect.Int32Kind:
		hashkey = fmt.Sprintf("%d", msg.ProtoReflect().Get(r.hashkeyd).Interface().(int32))
	case protoreflect.Int64Kind:
		hashkey = fmt.Sprintf("%d", msg.ProtoReflect().Get(r.hashkeyd).Interface().(int64))
	case protoreflect.Uint32Kind:
		hashkey = fmt.Sprintf("%d", msg.ProtoReflect().Get(r.hashkeyd).Interface().(uint32))
	case protoreflect.Uint64Kind:
		hashkey = fmt.Sprintf("%d", msg.ProtoReflect().Get(r.hashkeyd).Interface().(uint64))
	default:
		return "", fmt.Errorf("Unknow HashKey Kind %v", r.hashkeyd.Kind())
	}
	return hashkey, nil
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
		} else if m3grpcopt, ok := v.(*metapb.M3RpcOption); !ok {
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
			switch fieldd.Kind() {
			case protoreflect.StringKind:
			case protoreflect.Int32Kind:
			case protoreflect.Int64Kind:
			case protoreflect.Uint32Kind:
			case protoreflect.Uint64Kind:
				meta.hashkeyd = fieldd
			default:
				return fmt.Errorf("Invaild HashKey Kind %v", fieldd.Kind())
			}

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
