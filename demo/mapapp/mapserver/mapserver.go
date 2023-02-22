package mapserver

import (
	"context"
	dpb "m3game/demo/proto/pb"
	"m3game/runtime/transport"
	"m3game/server/async"
	"time"
)

var (
	_map map[string]int32
)

func init() {
	_map = make(map[string]int32)
}

func CreateMapSer() *MapSer {
	return &MapSer{
		Server: async.CreateServer("MapSer"),
	}
}

type MapSer struct {
	*async.Server
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

func (s *MapSer) TransportRegister() func(*transport.Transport) error {
	return func(t *transport.Transport) error {
		dpb.RegisterMapSerServer(t.GrpcSer(), s)
		return nil
	}
}
