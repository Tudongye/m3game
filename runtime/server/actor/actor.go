package actor

import (
	"context"
)

const (
	_actorkey = "_actorkey"
)

type Actor interface {
	ID() string
	Exit()                 // 退出
	OnInit() error         // 加载后触发
	OnTick() error         // 触发定时任务时
	OnExit() error         // 退出时触发
	OnSave() error         // 写回
	OnMoveIn([]byte) error // 数据自动迁出
	OnMoveOut() []byte     // 数据自动迁入
	ExitCh() chan struct{} // 退出标记
}

func WithActor(ctx context.Context, actor Actor) context.Context {
	return context.WithValue(ctx, _actorkey, actor)
}

func ParseActor(ctx context.Context) Actor {
	var a Actor
	if t := ctx.Value(_actorkey); t == nil {
		return a
	} else {
		return t.(Actor)
	}
}

func ActorBaseCreator(actorid string) *ActorBase {
	return &ActorBase{
		actorid: actorid,
		exitch:  make(chan struct{}, 1),
	}
}

type ActorBase struct {
	actorid string
	exitch  chan struct{}
}

func (a *ActorBase) ID() string {
	return a.actorid
}

func (a *ActorBase) OnInit() error {
	return nil
}

func (a *ActorBase) OnTick() error {
	return nil
}

func (a *ActorBase) OnExit() error {
	return nil
}

func (a *ActorBase) OnMoveIn([]byte) error {
	return nil
}

func (a *ActorBase) OnMoveOut() []byte {
	return nil
}

func (a *ActorBase) Save() error {
	return nil
}

func (a *ActorBase) Exit() {
	select {
	case a.exitch <- struct{}{}:
		return
	default:
		return
	}
}

func (a *ActorBase) ExitCh() chan struct{} {
	return a.exitch
}
