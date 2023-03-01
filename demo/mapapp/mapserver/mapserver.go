package mapserver

import (
	"context"
	dpb "m3game/demo/proto/pb"
	"m3game/server/async"
	"time"

	"google.golang.org/grpc"
)

var (
	_map map[string]int32
)

func init() {
	_map = make(map[string]int32)
}

func New() *MapSer {
	return &MapSer{
		Server: async.New("MapSer"),
	}
}

type MapSer struct {
	*async.Server
	dpb.UnimplementedMapSerServer
}

func (d *MapSer) Move(ctx context.Context, in *dpb.Move_Req) (*dpb.Move_Rsp, error) {
	out := new(dpb.Move_Rsp)
	name := in.Name
	distance := in.Distance
	if _, ok := _map[name]; !ok {
		_map[name] = 0
	}
	time.Sleep(time.Duration(distance) * time.Second)
	_map[name] += int32(distance)
	out.Name = name
	out.Location = _map[name]
	return out, nil
}

func (s *MapSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		dpb.RegisterMapSerServer(t, s)
		return nil
	}
}
