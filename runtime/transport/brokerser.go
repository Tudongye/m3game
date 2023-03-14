package transport

import (
	"context"
	"fmt"
	"m3game/config"
	"m3game/meta/metapb"
	"m3game/plugins/broker"
	"m3game/plugins/log"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

// Process Nty,BroadCast,MultiCast use Broker

func BrokerSerTopic(s string) string {
	return fmt.Sprintf("BrokerSer-%s", s)
}

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

	if err := n.broker.Subscribe(broker.GenTopic(BrokerSerTopic(string(config.GetAppID()))), n.recvbytes); err != nil {
		return errors.Wrapf(err, "Subscribe %s", broker.GenTopic(BrokerSerTopic(string(config.GetAppID()))))
	}
	if err := n.broker.Subscribe(broker.GenTopic(BrokerSerTopic(string(config.GetSvcID()))), n.recvbytes); err != nil {
		return errors.Wrapf(err, "Subscribe %s", broker.GenTopic(BrokerSerTopic(string(config.GetSvcID()))))
	}

	return nil
}

func (n *brokerSer) recvbytes(buff []byte) {
	bmsg := &metapb.BrokerMsg{}
	if err := proto.Unmarshal(buff, bmsg); err != nil {
		log.Error("Unmarshal buff err %s", err.Error())
		return
	}
	if handlerfunc, ok := n.handlers[bmsg.Method]; !ok {
		log.Error("not find method %s", bmsg.Method)
		return
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(_cfg.BroadTimeOut))
		defer cancel()
		metas := make(map[string]string)
		for _, meta := range bmsg.Metas {
			metas[meta.Key] = meta.Value
		}
		md := metadata.New(metas)
		ctx = metadata.NewIncomingContext(ctx, md)
		handlerfunc(ctx,
			func(pkg interface{}) error {
				return proto.Unmarshal(bmsg.Content, pkg.(proto.Message))
			},
			n.serverinterceptor)
	}
}

func (n *brokerSer) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	for _, it := range sd.Methods {
		path := fmt.Sprintf("/%v/%v", sd.ServiceName, it.MethodName)
		hander := it.Handler
		n.handlers[path] = func(ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
			return hander(ss, ctx, dec, interceptor)
		}
		log.Info("Register BrokerMethod => %v", path)
	}
}

func (n *brokerSer) send(ctx context.Context, topic string, method string, req interface{}, opts ...grpc.CallOption) error {
	bmsg := &metapb.BrokerMsg{}
	bmsg.Method = method
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		for k, vlist := range md {
			if len(vlist) > 0 {
				bmsg.Metas = append(bmsg.Metas, &metapb.Meta{Key: k, Value: vlist[0]})
			}
		}
	}
	var err error
	if reqmsg, ok := req.(proto.Message); !ok {
		return errors.New("")
	} else if bmsg.Content, err = proto.Marshal(reqmsg); err != nil {
		return err
	}
	if buff, err := proto.Marshal(bmsg); err != nil {
		return err
	} else {
		return n.broker.Publish(topic, buff)
	}
}
