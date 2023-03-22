package transport

import (
	"context"
	"fmt"
	"m3game/config"
	"m3game/meta/metapb"
	"m3game/plugins/broker"
	"m3game/plugins/log"
	"m3game/util"
	"sync"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

// Process Nty,BroadCast,MultiCast use Broker

type brokerhandlerFunc func(ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error)

func newBrokerSer(cfg TransportCfg, interceptor grpc.UnaryServerInterceptor) *brokerSer {
	return &brokerSer{
		cfg:               cfg,
		serverinterceptor: interceptor,
	}
}

type brokerSer struct {
	cfg               TransportCfg
	handlermap        sync.Map
	serverinterceptor grpc.UnaryServerInterceptor
	broker            broker.Broker
}

var (
	_ grpc.ServiceRegistrar = (*brokerSer)(nil)
)

func (n *brokerSer) registerBroker(b broker.Broker) error {
	n.broker = b
	if err := n.broker.Subscribe(broker.GenTopic(util.BrokerSerTopic(string(config.GetAppID()))), n.recvbytes); err != nil {
		return errors.Wrapf(err, "Subscribe %s", broker.GenTopic(util.BrokerSerTopic(string(config.GetAppID()))))
	}
	if err := n.broker.Subscribe(broker.GenTopic(util.BrokerSerTopic(string(config.GetSvcID()))), n.recvbytes); err != nil {
		return errors.Wrapf(err, "Subscribe %s", broker.GenTopic(util.BrokerSerTopic(string(config.GetSvcID()))))
	}

	return nil
}

func (n *brokerSer) recvbytes(bytes []byte) {
	bmsg := &metapb.BrokerMsg{}
	if err := proto.Unmarshal(bytes, bmsg); err != nil {
		log.Error("Unmarshal bytes err %s", err.Error())
		return
	}
	if value, ok := n.handlermap.Load(bmsg.Method); !ok {
		log.Error("not find method %s", bmsg.Method)
		return
	} else {
		handlerfunc := value.(brokerhandlerFunc)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(n.cfg.BroadcastTimeout))
		defer cancel()
		metad := make(map[string]string)
		for _, meta := range bmsg.Metas {
			metad[meta.Key] = meta.Value
		}
		md := metadata.New(metad)
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
		handler := it.Handler
		f := func(ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
			return handler(ss, ctx, dec, interceptor)
		}
		n.handlermap.Store(path, f)
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
		return fmt.Errorf("SendMsg is not PBMsg, for topic %s", topic)
	} else if bmsg.Content, err = proto.Marshal(reqmsg); err != nil {
		return err
	}
	if bytes, err := proto.Marshal(bmsg); err != nil {
		return err
	} else {
		return n.broker.Publish(topic, bytes)
	}
}
