package worldser

import (
	"context"
	"errors"
	"m3game/example/proto/pb"
	"m3game/gameplay/world"
	"m3game/plugins/log"
	"sync"
	"time"
)

const (
	MAX_CELLCMDCHAN = 100 // Cell 队列长度
	MAX_CELLCMDLOOP = 10  // Cell 单帧处理数

	MAX_WORLDCMDCHAN = 100 // World 队列长度
	MAX_WORLDCMDLOOP = 10  // World 单帧处理数

	MAX_TIMEOUT = 10
)

const (
	CMD_VIEWPOSITION = 1 // 查看 Postion所在Cell
	CMD_CREATEENTITY = 2 // 创建Entity
	CMD_MOVEENTITY   = 3 // 移动Entity
)

var (
	err_ch_full     = errors.New("err_ch_full")
	err_cmd_timeout = errors.New("err_cmd_timeout")
)

type Cell struct {
	*world.DefaultCell2d
	inputs []*pb.Entity // 输入Entity
	mu     sync.Mutex
}

type CellShare struct {
	entitys map[string]*pb.Entity // Cell当前Entity
	cmdch   chan *WorldCmd        // View请求列表
}

func (c *Cell) MoveIn(e *pb.Entity) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.inputs = append(c.inputs, e)
}

func NewWorldReciver() *WorldReciver {
	wr := &WorldReciver{}
	wr.cmdch = make(chan *WorldCmd, MAX_WORLDCMDCHAN)
	wr.cellshares = make(map[int]*CellShare)
	wr.entitypool = sync.Pool{
		New: func() interface{} {
			return &pb.Entity{
				Name:     "",
				SrcPos:   &pb.Position{},
				DstPos:   &pb.Position{},
				CurPos:   &pb.Position{},
				MoveTime: 0,
			}
		},
	}
	return wr
}

type WorldReciver struct {
	entitypool sync.Pool
	entitys    sync.Map
	cellshares map[int]*CellShare
	mu         sync.Mutex
	cmdch      chan *WorldCmd
	world      world.World2d[*Cell]
}

func (wr *WorldReciver) HookWorld(w world.World2d[*Cell]) {
	wr.world = w
}

func (wr *WorldReciver) NewCell2d(idex, x1, x2, y1, y2 int) *Cell {
	wr.mu.Lock()
	defer wr.mu.Unlock()
	if _, ok := wr.cellshares[idex]; !ok {
		wr.cellshares[idex] = &CellShare{
			entitys: make(map[string]*pb.Entity),
			cmdch:   make(chan *WorldCmd, MAX_CELLCMDCHAN),
		}
	}
	cell := &Cell{
		DefaultCell2d: world.NewCell2d(idex, x1, x2, y1, y2),
	}
	return cell
}

func (wr *WorldReciver) CellFunc2d(ctx context.Context, p world.Page, idex int, w world.World2d[*Cell]) {
	oldc := w.GetBakGrid().GetCell(idex)
	share := wr.cellshares[idex]
	// 处理入队列（上帧收尾）
	for _, e := range oldc.inputs {
		share.entitys[e.Name] = e
	}
	oldc.inputs = oldc.inputs[:0]

	// 处理View请求
	if len(share.cmdch) > 0 {
		var viewlist []*pb.Entity
		for _, e := range share.entitys {
			viewlist = append(viewlist, e)
		}
		for i := 0; i < MAX_CELLCMDLOOP; i++ {
			select {
			case cmd := <-share.cmdch:
				cmd.ViewPositionCmd.RspCh <- viewlist
				continue
			default:
				break
			}
		}
	}

	// 处理本Cell
	var moveouts []string
	for _, e := range share.entitys {
		if e.MoveTime != 0 {
			// 计算移动距离
			move := CalMove(p.PageStamp-e.MoveTime) - e.Moved
			// 移动
			if MovePos(e, move) {
				e.MoveTime = 0
				e.Moved = 0
			} else {
				e.Moved += move
			}
			curidex := w.GetCellIndex(int(e.CurPos.X), int(e.CurPos.Y))
			if curidex != idex {
				// 移动到其他Cell
				w.GetGrid().GetCell(curidex).MoveIn(e)
				moveouts = append(moveouts, e.Name)
			}
		}
	}
	for _, name := range moveouts {
		delete(share.entitys, name)
	}
	return
}

func (wr *WorldReciver) WorldFunc2d(ctx context.Context, p world.Page, w world.World2d[*Cell]) {
	for i := 0; i < MAX_WORLDCMDLOOP; i++ {
		select {
		case cmd := <-wr.cmdch:
			wr.ProcessCmd(cmd, p, w)
			continue
		default:
			break
		}
	}
	return
}

func (wr *WorldReciver) PushCmd(cmd *WorldCmd) bool {
	select {
	case wr.cmdch <- cmd:
		return true
	default:
		return false
	}
}

func (wr *WorldReciver) PushCellCmd(cmd *WorldCmd, idex int) bool {
	select {
	case wr.cellshares[idex].cmdch <- cmd:
		return true
	default:
		return false
	}
}

func (wr *WorldReciver) ProcessCmd(cmd *WorldCmd, p world.Page, w world.World2d[*Cell]) {
	log.Debug("ProcessCmd %v", cmd.Cmd)
	switch cmd.Cmd {
	case CMD_CREATEENTITY: // 创建Entity
		e := wr.AllocEntity(cmd.CreateEntityCmd.Name)
		e.SrcPos.X = cmd.CreateEntityCmd.Pos.X
		e.SrcPos.Y = cmd.CreateEntityCmd.Pos.Y
		e.CurPos.X = cmd.CreateEntityCmd.Pos.X
		e.CurPos.Y = cmd.CreateEntityCmd.Pos.Y
		idex := w.GetCellIndex(int(e.CurPos.X), int(e.CurPos.Y))
		w.GetGrid().GetCell(idex).MoveIn(e)
		cmd.CreateEntityCmd.RspCh <- e
		return
	case CMD_MOVEENTITY: // 移动Move
		e := wr.GetEntity(cmd.MoveEntityCmd.Name)
		MoveEntity(e, cmd.MoveEntityCmd.DstPos, p.PageStamp)
		cmd.MoveEntityCmd.RspCh <- e
		return
	default:
		log.Error("UnKnow Cmd %v", cmd.Cmd)
		return
	}
}

func (wr *WorldReciver) CreateEntity(name string, pos *pb.Position) (*pb.Entity, error) {
	cmd := &WorldCmd{
		Cmd: CMD_CREATEENTITY,
		CreateEntityCmd: &CreateEntityCmd{
			Name:  name,
			Pos:   pos,
			RspCh: make(chan *pb.Entity),
		},
	}
	if !wr.PushCmd(cmd) {
		return nil, err_ch_full
	}
	select {
	case e := <-cmd.CreateEntityCmd.RspCh:
		return e, nil
	case <-time.Tick(MAX_TIMEOUT * time.Second):
		return nil, err_cmd_timeout
	}
}

func (wr *WorldReciver) MoveEntity(name string, dstpos *pb.Position) (*pb.Entity, error) {
	cmd := &WorldCmd{
		Cmd: CMD_MOVEENTITY,
		MoveEntityCmd: &MoveEntityCmd{
			Name:   name,
			DstPos: dstpos,
			RspCh:  make(chan *pb.Entity),
		},
	}
	if !wr.PushCmd(cmd) {
		return nil, err_ch_full
	}
	select {
	case e := <-cmd.MoveEntityCmd.RspCh:
		return e, nil
	case <-time.Tick(MAX_TIMEOUT * time.Second):
		return nil, err_cmd_timeout
	}
}

func (wr *WorldReciver) ViewPosition(pos *pb.Position) ([]*pb.Entity, error) {
	idex := wr.world.GetCellIndex(int(pos.X), int(pos.Y))
	cmd := &WorldCmd{
		Cmd: CMD_VIEWPOSITION,
		ViewPositionCmd: &ViewPositionCmd{
			Pos:   pos,
			RspCh: make(chan []*pb.Entity),
		},
	}
	if !wr.PushCellCmd(cmd, idex) {
		return nil, err_ch_full
	}
	select {
	case e := <-cmd.ViewPositionCmd.RspCh:
		return e, nil
	case <-time.Tick(MAX_TIMEOUT * time.Second):
		return nil, err_cmd_timeout
	}
}

func (wr *WorldReciver) RecoverEntity(e *pb.Entity) {
	wr.entitys.Delete(e.Name)
	wr.FreeEntity(e)
}

func (wr *WorldReciver) GetEntity(name string) *pb.Entity {
	if v, ok := wr.entitys.Load(name); !ok {
		return nil
	} else {
		return v.(*pb.Entity)
	}
}

func (wr *WorldReciver) AllocEntity(name string) *pb.Entity {
	e := wr.entitypool.Get().(*pb.Entity)
	e.Name = name
	e.SrcPos.X = 0
	e.SrcPos.Y = 0
	e.DstPos.X = 0
	e.DstPos.Y = 0
	e.CurPos.X = 0
	e.CurPos.Y = 0
	e.MoveTime = 0
	e.Moved = 0
	wr.entitys.Store(e.Name, e)
	return e
}

func (wr *WorldReciver) FreeEntity(e *pb.Entity) {
	wr.entitypool.Put(e)
}

func MatchPos(a, b *pb.Position) bool {
	if a.X == b.X && a.Y == b.Y {
		return true
	}
	return false
}

func CalMove(movetime int64) int32 {
	move := int32(movetime / 1000)
	return move
}

func MoveEntity(e *pb.Entity, dstpos *pb.Position, movetime int64) {
	e.SrcPos.X = e.CurPos.X
	e.SrcPos.Y = e.CurPos.Y
	e.DstPos.X = dstpos.X
	e.DstPos.Y = dstpos.Y
	e.MoveTime = movetime
	e.Moved = 0
}

func MovePos(e *pb.Entity, move int32) bool {
	dst := e.DstPos
	cur := e.CurPos
	if move <= 0 {
		return false
	}
	if cur.X < dst.X {
		if dst.X-cur.X > move {
			cur.X += move
			return false
		}
		move -= dst.X - cur.X
		cur.X = dst.X
	} else if cur.X > dst.X {
		if cur.X-dst.X > move {
			cur.X -= move
			return false
		}
		move -= cur.X - dst.X
		cur.X = dst.X
	}

	if cur.Y < dst.Y {
		if dst.Y-cur.Y > move {
			cur.Y += move
			return false
		}
		move -= dst.Y - cur.Y
		cur.Y = dst.Y
	} else if cur.Y > dst.Y {
		if cur.Y-dst.Y > move {
			cur.Y -= move
			return false
		}
		move -= cur.Y - dst.Y
		cur.Y = dst.Y
	}
	return true
}

type WorldCmd struct {
	Cmd             int
	ViewPositionCmd *ViewPositionCmd
	CreateEntityCmd *CreateEntityCmd
	MoveEntityCmd   *MoveEntityCmd
}

type ViewPositionCmd struct {
	Pos   *pb.Position
	RspCh chan []*pb.Entity
}

type CreateEntityCmd struct {
	Name  string
	Pos   *pb.Position
	RspCh chan *pb.Entity
}

type MoveEntityCmd struct {
	Name   string
	DstPos *pb.Position
	RspCh  chan *pb.Entity
}
