package rpc

import (
	"fmt"
	"m3game/log"
	"m3game/proto/pb"

	"google.golang.org/protobuf/proto"

	"google.golang.org/protobuf/reflect/protoreflect"
)

const routeheadname = "RouteHead"

var RPCMetas = make(map[protoreflect.FullName]*RPCMeta)

type RPCMeta struct {
	grpcoption *pb.M3GRPCOption             // rpc_option
	routeheadd protoreflect.FieldDescriptor // RouteHead,for Assignment
	hashkeyd   protoreflect.FieldDescriptor
}

func (r *RPCMeta) GrpcOption() *pb.M3GRPCOption {
	return r.grpcoption
}
func (r *RPCMeta) RouteHead() protoreflect.FieldDescriptor {
	return r.routeheadd
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
			grpcoption: nil,
			routeheadd: nil,
			hashkeyd:   nil,
		}
		// eache Rpc must have rpc_option
		if v := proto.GetExtension(rpcde.Options(), pb.E_RpcOption); v == nil {
			panic(fmt.Sprintf("RPC %s not have E_RpcOption", rpcde.Name()))
		} else if m3grpcopt, ok := v.(*pb.M3GRPCOption); !ok {
			panic(fmt.Sprintf("RPC %s E_RpcOption type err", rpcde.Name()))
		} else if m3grpcopt == nil {
			panic(fmt.Sprintf("RPC %s E_RpcOption is nil", rpcde.Name()))
		} else {
			RPCMetas[inputname].grpcoption = m3grpcopt
		}

		if fieldd := inputd.Fields().ByName(routeheadname); fieldd == nil {
			panic(fmt.Sprintf("RPC %s input %s not have RouteHead", rpcde.Name(), inputname))
		} else {
			RPCMetas[inputname].routeheadd = fieldd
		}

		if fieldd := inputd.Fields().ByName(protoreflect.Name(RPCMetas[inputname].grpcoption.RouteKey)); fieldd != nil {
			RPCMetas[inputname].hashkeyd = fieldd
		}
		log.Info("RPC Registor: Svc => %s, Method => %s, Input => %s, RPC => %s", serviced.Name(), methodd.Name(), inputname, rpcde.Name())
	}
	return nil
}

func Method(inputname protoreflect.FullName) *RPCMeta {
	return RPCMetas[inputname]
}
