package world

import (
	"context"
	"sync"
	"time"
)

const (
	_parallelflag = "_parallelflag"
)

func ParallelCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, _parallelflag, 1)
}

func IsParallelCtx(ctx context.Context) bool {
	if ctx.Value(_parallelflag) == nil {
		return false
	}
	return true
}

func Run[TC Cell](ctx context.Context, w World[TC]) {
	lastPage := Page{
		PageInter: 0,
		PageStamp: time.Now().UnixMilli(),
	}
	cellChan := make(chan int, w.CellNum())
	ticker := time.NewTicker(time.Duration(w.PageInter()) * time.Millisecond)
	for {
		CurPage := Page{
			PageStamp: time.Now().UnixMilli(),
		}
		CurPage.PageInter = CurPage.PageStamp - lastPage.PageStamp
		// 并行阶段
		for i := 0; i < w.CellNum(); i++ {
			cellChan <- i
		}
		var wg sync.WaitGroup
		for i := 0; i < w.WorkNum(); i++ {
			wg.Add(1)
			go func(p Page) {
				defer wg.Done()
				end := false
				for {
					select {
					case idex := <-cellChan:
						pctx := ParallelCtx(ctx)
						w.Reciver().CellFunc(pctx, p, idex, w)
					default:
						end = true
						break
					}
					if end {
						break
					}
				}
			}(CurPage)
		}
		wg.Wait()
		// 串行阶段
		w.Reciver().WorldFunc(ctx, CurPage, w)
		w.SwapBakGrid()
		lastPage = CurPage
		<-ticker.C
	}
}
