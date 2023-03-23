package lease

import (
	"context"
	"errors"
	"fmt"
	"m3game/plugins/log"
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

func New(lease Lease) (Lease, error) {
	if _lease != nil {
		log.Fatal("Lease Only One")
		return nil, fmt.Errorf("Lease is newed %s", _lease.Factory().Name())
	}
	_lease = lease
	return _lease, nil
}

func Instance() Lease {
	if _lease == nil {
		log.Fatal("Lease not newd")
		return nil
	}
	return _lease
}

func SetReciver(l LeaseReciver) {
	_leasereciver = l
}

func GetReciver() LeaseReciver {
	return _leasereciver
}

func AllocLease(ctx context.Context, id string, f LeaseMoveOutFunc) error {
	if _lease == nil {
		return errors.New("Lease is not initialized")
	}
	return _lease.AllocLease(ctx, id, f)
}

func FreeLease(ctx context.Context, id string) error {
	if _lease == nil {
		return errors.New("Lease is not initialized")
	}
	return _lease.FreeLease(ctx, id)
}

func KickLease(ctx context.Context, id string) ([]byte, error) {
	if _lease == nil {
		return nil, errors.New("Lease is not initialized")
	}
	return _lease.KickLease(ctx, id)
}

func RecvKickLease(ctx context.Context, id string) ([]byte, error) {
	if _lease == nil {
		return nil, errors.New("Lease is not initialized")
	}
	return _lease.RecvKickLease(ctx, id)
}

func GetLease(ctx context.Context, id string) ([]byte, error) {
	if _lease == nil {
		return nil, errors.New("Lease is not initialized")
	}
	return _lease.GetLease(ctx, id)
}
