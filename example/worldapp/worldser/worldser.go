package worldser

import (
	"context"
	"m3game/example/proto/pb"
	"m3game/gameplay/world"
	"m3game/plugins/log"
	"m3game/runtime/rpc"
	"m3game/runtime/server/multi"

	"google.golang.org/grpc"
)

func init() {
	// 注册RPC信息到框架层
	if err := rpc.InjectionRPC(pb.File_world_proto.Services().Get(0)); err != nil {
		log.Fatal("InjectionRPC WorldSer %s", err.Error())
	}
}

func New() *WorldSer {
	worldser := &WorldSer{
		Server: multi.New("WorldSer"), // 以MultiSer为基础构建SimpleSer
	}
	worldser.wr = NewWorldReciver()
	world2d := world.NewWorld2d[*Cell](worldser.wr, 1000, 1000, 1000, 1000, 1, 10)
	worldser.wr.HookWorld(world2d)
	worldser.world = world2d
	return worldser
}

type WorldSer struct {
	*multi.Server
	pb.UnimplementedWorldSerServer
	world world.World[*Cell]
	wr    *WorldReciver
}

func (d *WorldSer) RunWorld(ctx context.Context) {
	go func() {
		world.Run(ctx, d.world)
	}()
}

func (d *WorldSer) CreateEntity(ctx context.Context, req *pb.CreateEntity_Req) (*pb.CreateEntity_Rsp, error) {
	rsp := &pb.CreateEntity_Rsp{}
	if e, err := d.wr.CreateEntity(req.GetName(), req.GetSrcPos()); err != nil {
		return rsp, err
	} else {
		rsp.Entity = e
		return rsp, nil
	}
}

func (d *WorldSer) MoveEntity(ctx context.Context, req *pb.MoveEntity_Req) (*pb.MoveEntity_Rsp, error) {
	rsp := &pb.MoveEntity_Rsp{}
	if e, err := d.wr.MoveEntity(req.GetName(), req.GetDstPos()); err != nil {
		return rsp, err
	} else {
		rsp.Entity = e
		return rsp, nil
	}
}

func (d *WorldSer) ViewEntity(context.Context, *pb.ViewEntity_Req) (*pb.ViewEntity_Rsp, error) {
	return nil, nil
}

func (d *WorldSer) ViewPosition(ctx context.Context, req *pb.ViewPosition_Req) (*pb.ViewPosition_Rsp, error) {
	rsp := &pb.ViewPosition_Rsp{}
	if e, err := d.wr.ViewPosition(req.GetPos()); err != nil {
		return rsp, err
	} else {
		rsp.Entitys = e
		return rsp, nil
	}
}

// 将SimpleSer注册到grpcser
func (s *WorldSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		pb.RegisterWorldSerServer(t, s)
		return nil
	}
}
