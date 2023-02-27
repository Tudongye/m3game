package transport

import (
	"context"
	"fmt"
	"log"
	"m3game/broker"
	"m3game/config"
	"m3game/proto/pb"
	"m3game/util"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// Process Nty,BroadCast,MutilCast use Broker

type grpchandlerFunc func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error)
type brokerhandlerFunc func(ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error)

type BrokerSer struct {
	mu                sync.Mutex
	handlers          map[string]brokerhandlerFunc
	serverinterceptor grpc.UnaryServerInterceptor
	broker            broker.Broker
}

func genhandlerFunc(s interface{}, h grpchandlerFunc) brokerhandlerFunc {
	return func(ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		return h(s, ctx, dec, interceptor)
	}
}

func CreateBrokerSer(interceptor grpc.UnaryServerInterceptor) *BrokerSer {
	return &BrokerSer{
		serverinterceptor: interceptor,
		handlers:          make(map[string]brokerhandlerFunc),
	}
}

func (n *BrokerSer) registerBroker(broker broker.Broker) {
	n.broker = broker
	idstr := config.GetIDStr()
	if envid, worldid, funcid, _, err := util.AppStr2ID(idstr); err != nil {
		log.Panicf(err.Error())
	} else {
		if err := broker.Subscribe(util.GenTopic(idstr), n.recvbytes); err != nil {
			log.Panicf(err.Error())
		}
		if err := broker.Subscribe(util.GenTopic(util.SvcID2Str(envid, worldid, funcid)), n.recvbytes); err != nil {
			log.Panicf(err.Error())
		}
		if err := broker.Subscribe(util.GenTopic(util.WorldID2Str(envid, worldid)), n.recvbytes); err != nil {
			log.Panicf(err.Error())
		}
		if err := broker.Subscribe(util.GenTopic(util.EnvID2Str(envid)), n.recvbytes); err != nil {
			log.Panicf(err.Error())
		}
	}
}

func (n *BrokerSer) recvbytes(buff []byte) {
	bmsg := &pb.BrokerPkg{}
	if err := proto.Unmarshal(buff, bmsg); err != nil {
		return
	}
	if handlerfunc, ok := n.handlers[bmsg.FullMethod]; !ok {
		return
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(_cfg.BroadTimeOut))
		defer cancel()
		handlerfunc(ctx,
			func(pkg interface{}) error {
				return proto.Unmarshal(bmsg.Content, pkg.(proto.Message))
			}, n.serverinterceptor)
	}
}

func (n *BrokerSer) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	for _, it := range sd.Methods {
		path := fmt.Sprintf("/%v/%v", sd.ServiceName, it.MethodName)
		n.handlers[path] = genhandlerFunc(ss, grpchandlerFunc(it.Handler))
		log.Printf("BrokerSer register path => %v", path)
	}
}

func (n *BrokerSer) SendToBroker(topic string, method string, msg proto.Message) error {
	bmsg := &pb.BrokerPkg{}
	bmsg.FullMethod = method
	var err error
	if bmsg.Content, err = proto.Marshal(msg); err != nil {
		return err
	}
	if buff, err := proto.Marshal(bmsg); err != nil {
		return err
	} else {
		return n.broker.Publish(topic, buff)
	}
}
