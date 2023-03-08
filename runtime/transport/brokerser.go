package transport

import (
	"context"
	"fmt"
	"m3game/broker"
	"m3game/config"
	"m3game/log"
	"m3game/proto/pb"
	"m3game/util"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// Process Nty,BroadCast,MutilCast use Broker

func BrokerSerTopic(s string) string {
	return fmt.Sprintf("BrokerSer-%s", s)
}

type grpchandlerFunc func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error)
type brokerhandlerFunc func(ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error)

func newBrokerSer(interceptor grpc.UnaryServerInterceptor) *brokerSer {
	return &brokerSer{
		serverinterceptor: interceptor,
		handlers:          make(map[string]brokerhandlerFunc),
	}
}

type brokerSer struct {
	handlers          map[string]brokerhandlerFunc
	serverinterceptor grpc.UnaryServerInterceptor
	broker            broker.Broker
}

var (
	_ grpc.ServiceRegistrar = (*brokerSer)(nil)
)

func (n *brokerSer) registerBroker(b broker.Broker) error {
	n.broker = b
	idstr := config.GetIDStr()
	if envid, worldid, funcid, _, err := util.AppStr2ID(idstr); err != nil {
		return err
	} else {
		if err := n.broker.Subscribe(broker.GenTopic(BrokerSerTopic(idstr)), n.recvbytes); err != nil {
			return errors.Wrapf(err, "Subscribe %s", broker.GenTopic(BrokerSerTopic(idstr)))
		}
		if err := n.broker.Subscribe(broker.GenTopic(BrokerSerTopic(util.SvcID2Str(envid, worldid, funcid))), n.recvbytes); err != nil {
			return errors.Wrapf(err, "Subscribe %s", broker.GenTopic(BrokerSerTopic(util.SvcID2Str(envid, worldid, funcid))))
		}
		if err := n.broker.Subscribe(broker.GenTopic(BrokerSerTopic(util.WorldID2Str(envid, worldid))), n.recvbytes); err != nil {
			return errors.Wrapf(err, "Subscribe %s", broker.GenTopic(BrokerSerTopic(util.WorldID2Str(envid, worldid))))
		}
		if err := n.broker.Subscribe(broker.GenTopic(BrokerSerTopic(util.EnvID2Str(envid))), n.recvbytes); err != nil {
			return errors.Wrapf(err, "Subscribe %s", broker.GenTopic(BrokerSerTopic(util.EnvID2Str(envid))))
		}
	}
	return nil
}

func (n *brokerSer) recvbytes(buff []byte) {
	bmsg := &pb.BrokerPkg{}
	if err := proto.Unmarshal(buff, bmsg); err != nil {
		log.Error("Unmarshal buff err %s", err.Error())
		return
	}
	if handlerfunc, ok := n.handlers[bmsg.FullMethod]; !ok {
		log.Error("not find method %s", bmsg.FullMethod)
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

func (n *brokerSer) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	for _, it := range sd.Methods {
		path := fmt.Sprintf("/%v/%v", sd.ServiceName, it.MethodName)
		n.handlers[path] = genbrokerhandlerFunc(ss, grpchandlerFunc(it.Handler))
		log.Info("Register BrokerMethod => %v", path)
	}
}

func (n *brokerSer) send(topic string, method string, msg proto.Message) error {
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

func genbrokerhandlerFunc(s interface{}, h grpchandlerFunc) brokerhandlerFunc {
	return func(ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		return h(s, ctx, dec, interceptor)
	}
}
