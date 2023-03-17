package lease

import (
	"context"
	"m3game/runtime/plugin"
)

var (
	_lease        Lease
	_leasereciver LeaseReciver
)

type LeaseMoveOutFunc func(context.Context) ([]byte, error) // 租约退出

func DefaultLeaseMoveOutFunc(context.Context) ([]byte, error) {
	return nil, nil
}

type Lease interface {
	plugin.PluginIns
	AllocLease(ctx context.Context, id string, f LeaseMoveOutFunc) error // 获取租约
	FreeLease(ctx context.Context, id string) error                      // 释放租约
	KickLease(ctx context.Context, id string) ([]byte, error)            // 要求释放租约
	RecvKickLease(ctx context.Context, id string) ([]byte, error)        // 接受释放租约消息
	GetLease(ctx context.Context, id string) ([]byte, error)
}

type LeaseReciver interface {
	SendKickLease(ctx context.Context, id string, app string) ([]byte, error) // 发送释放租约消息
}

func Set(l Lease) {
	_lease = l
}

func SetReciver(l LeaseReciver) {
	_leasereciver = l
}

func Get() Lease {
	return _lease
}

func GetReciver() LeaseReciver {
	return _leasereciver
}

func AllocLease(ctx context.Context, id string, f LeaseMoveOutFunc) error {
	return Get().AllocLease(ctx, id, f)
}

func FreeLease(ctx context.Context, id string) error {
	return Get().FreeLease(ctx, id)
}

func KickLease(ctx context.Context, id string) ([]byte, error) {
	return Get().KickLease(ctx, id)
}

func RecvKickLease(ctx context.Context, id string) ([]byte, error) {
	return Get().RecvKickLease(ctx, id)
}
func GetLease(ctx context.Context, id string) ([]byte, error) {
	return Get().GetLease(ctx, id)
}
