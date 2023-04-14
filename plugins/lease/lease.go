package lease

import (
	"context"
	"m3game/meta/errs"
	"m3game/plugins/log"
	"m3game/runtime/plugin"
)

var (
	_lease Lease
)

type Lease interface {
	plugin.PluginIns
	AllocLease(ctx context.Context, leaseid string) error             // 获取租约
	FreeLease(ctx context.Context, leaseid string) error              // 释放租约
	WhereIsLease(ctx context.Context, leaseid string) (string, error) // 获取租约位置
}

func New(lease Lease) (Lease, error) {
	if _lease != nil {
		log.Fatal("Lease Only One")
		return nil, errs.LeaseInsHasNewed.New("Lease is newed %s", _lease.Factory().Name())
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
