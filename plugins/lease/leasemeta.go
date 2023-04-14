package lease

import (
	"context"
	"fmt"
	"m3game/meta/errs"
	"sync"
)

type LeaseMoveOutFunc func(context.Context) ([]byte, error) // 租约退出

func DefaultLeaseMoveOutFunc(context.Context) ([]byte, error) {
	return nil, nil
}

type SendKickLeaseFunc func(ctx context.Context, leaseid string, appid string) ([]byte, error) // 发送租约踢出请求

func DefaultSendKickLeaseFunc(ctx context.Context, leaseid string, appid string) ([]byte, error) {
	return nil, nil
}

type LeaseMeta struct { // 单类租约元数据
	prefix       string            // 租约前缀
	moveoutfuncs sync.Map          // 租约退出回调
	sendkickfunc SendKickLeaseFunc // 租约踢出请求
}

func NewLeaseMeta(prefix string, sendkickf SendKickLeaseFunc) *LeaseMeta {
	return &LeaseMeta{
		prefix:       prefix,
		sendkickfunc: sendkickf,
	}
}

func (lm *LeaseMeta) Prefix() string {
	return lm.prefix
}

func (lm *LeaseMeta) AllocLease(ctx context.Context, id string, f LeaseMoveOutFunc) error {
	if _lease == nil {
		return errs.LeaseInsIsNill.New("Lease is not initialized")
	}
	leaseid := fmt.Sprintf("%s-%s", lm.prefix, id)
	if _, ok := lm.moveoutfuncs.Load(leaseid); ok {
		return nil
	}
	if err := _lease.AllocLease(ctx, leaseid); err != nil {
		return err
	}
	lm.moveoutfuncs.Store(leaseid, f)
	return nil
}

func (lm *LeaseMeta) FreeLease(ctx context.Context, id string) ([]byte, error) {
	if _lease == nil {
		return nil, errs.LeaseInsIsNill.New("Lease is not initialized")
	}
	leaseid := fmt.Sprintf("%s-%s", lm.prefix, id)
	if v, ok := lm.moveoutfuncs.Load(leaseid); !ok {
		return nil, errs.LeaseMetaNotFindLease.New("")
	} else {
		if res, err := v.(LeaseMoveOutFunc)(ctx); err != nil {
			// 业务层租约踢出失败
			return nil, err
		} else {
			if err := _lease.FreeLease(ctx, leaseid); err != nil {
				return nil, err
			}
			lm.moveoutfuncs.Delete(leaseid)
			return res, nil
		}
	}
}

func (lm *LeaseMeta) KickLease(ctx context.Context, id string) ([]byte, error) {
	if _lease == nil {
		return nil, errs.LeaseInsIsNill.New("Lease is not initialized")
	}
	leaseid := fmt.Sprintf("%s-%s", lm.prefix, id)
	if appid, err := lm.WhereIsLease(ctx, id); err != nil {
		return nil, err
	} else {
		return lm.sendkickfunc(ctx, leaseid, appid)
	}
}

func (lm *LeaseMeta) RecvKickLease(ctx context.Context, leaseid string) ([]byte, error) {
	if _lease == nil {
		return nil, errs.LeaseInsIsNill.New("Lease is not initialized")
	}
	if v, ok := lm.moveoutfuncs.Load(leaseid); !ok {
		// 直接退约
		if err := _lease.FreeLease(ctx, leaseid); err != nil {
			return nil, err
		}
		return nil, nil
	} else {
		if res, err := v.(LeaseMoveOutFunc)(ctx); err != nil {
			// 业务层租约踢出失败
			return nil, err
		} else {
			if err := _lease.FreeLease(ctx, leaseid); err != nil {
				return nil, err
			}
			lm.moveoutfuncs.Delete(leaseid)
			return res, nil
		}
	}
}

func (lm *LeaseMeta) WhereIsLease(ctx context.Context, id string) (string, error) {
	if _lease == nil {
		return "", errs.LeaseInsIsNill.New("Lease is not initialized")
	}
	leaseid := fmt.Sprintf("%s-%s", lm.prefix, id)
	return _lease.WhereIsLease(ctx, leaseid)
}
