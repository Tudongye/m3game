package client

import (
	"fmt"
	"m3game/proto/pb"
	"m3game/util/log"

	pbproto "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// each rpc message must contain RouteHead
const routeheadname = "RouteHead"

// record rpc method option
type m3Method struct {
	grpcoption *pb.M3GRPCOption             // rpc_option
	routeheadd protoreflect.FieldDescriptor // RouteHead,for Assignment
	hashkeyd   protoreflect.FieldDescriptor
}

type Meta interface {
	SrcIns() *pb.RouteIns
	DstSvc() *pb.RouteSvc
	Method(n protoreflect.FullName) *m3Method
}

func NewMeta(serviced protoreflect.ServiceDescriptor, srcins *pb.RouteIns, dstsvc *pb.RouteSvc) Meta {
	m := &defaultMeta{
		srcins:  srcins,
		dstsvc:  dstsvc,
		methods: make(map[protoreflect.FullName]*m3Method),
	}
	for i := 0; i < serviced.Methods().Len(); i++ {
		methodd := serviced.Methods().Get(i)
		inputd := methodd.Input()
		inputname := inputd.FullName()
		rpcde := inputd.Parent()
		m.methods[inputname] = &m3Method{
			grpcoption: nil,
			routeheadd: nil,
			hashkeyd:   nil,
		}
		// eache Rpc must have rpc_option
		if v := pbproto.GetExtension(rpcde.Options(), pb.E_RpcOption); v == nil {
			panic(fmt.Sprintf("RPC %s not have E_RpcOption", rpcde.Name()))
		} else if m3grpcopt, ok := v.(*pb.M3GRPCOption); !ok {
			panic(fmt.Sprintf("RPC %s E_RpcOption type err", rpcde.Name()))
		} else {
			m.methods[inputname].grpcoption = m3grpcopt
		}

		if fieldd := inputd.Fields().ByName(routeheadname); fieldd == nil {
			panic(fmt.Sprintf("RPC %s input %s not have RouteHead", rpcde.Name(), inputname))
		} else {
			m.methods[inputname].routeheadd = fieldd
		}

		if fieldd := inputd.Fields().ByName(protoreflect.Name(m.methods[inputname].grpcoption.RouteKey)); fieldd != nil {
			m.methods[inputname].hashkeyd = fieldd
		}
		log.Fatal("RPC Registor: Svc => %s, Method => %s, Input => %s, RPC => %s", serviced.Name(), methodd.Name(), inputname, rpcde.Name())
	}
	return m
}

type defaultMeta struct {
	srcins  *pb.RouteIns
	dstsvc  *pb.RouteSvc
	methods map[protoreflect.FullName]*m3Method
}

func (m *defaultMeta) SrcIns() *pb.RouteIns {
	return m.srcins
}
func (m *defaultMeta) DstSvc() *pb.RouteSvc {
	return m.dstsvc
}

func (m *defaultMeta) Method(n protoreflect.FullName) *m3Method {
	return m.methods[n]
}
