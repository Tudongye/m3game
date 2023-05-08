package world

import (
	"context"
)

type World[TC Cell] interface {
	GetGrid() Grid[TC]    // 获取当前帧网格
	GetBakGrid() Grid[TC] // 获取上帧网格
	SwapBakGrid()         // 网格交换
	CellNum() int
	WorkNum() int
	PageInter() int
	Reciver() WorldReciver[TC]
}

type Grid[TC Cell] interface {
	GetCell(idex int) TC // 获取单元格
}

type Cell interface {
}

type WorldReciver[TC Cell] interface {
	NewCell(idex int) TC
	CellFunc(ctx context.Context, p Page, idex int, w World[TC]) // 并行阶段
	WorldFunc(ctx context.Context, p Page, w World[TC])          // 串行阶段
}

type Page struct {
	PageStamp int64 // 当前帧时间戳
	PageInter int64 // 当前帧与上帧间隔
}

func NewWorld[TC Cell](reciver WorldReciver[TC], cellnum int, worknum int, pageinter int) World[TC] {
	w := &DefaultWorld[TC]{
		Bak:       0,
		cellnum:   cellnum,
		worknum:   worknum,
		reciver:   reciver,
		pageinter: pageinter,
	}
	for i := 0; i < len(w.Grids); i++ {
		g := &DefaultGrid[TC]{
			CellNum: cellnum,
			Cells:   make([]TC, cellnum),
		}
		for j := 0; j < cellnum; j++ {
			g.Cells[j] = reciver.NewCell(j)
		}
		w.Grids[i] = g
	}

	return w
}

// 世界
type DefaultWorld[TC Cell] struct {
	Grids     [2]Grid[TC]      // 网格
	Bak       int              // 备份标记
	cellnum   int              // 单元格数量
	worknum   int              // 并行工作数
	reciver   WorldReciver[TC] // 业务接口
	pageinter int              // 单帧时长
}

func (w *DefaultWorld[TC]) GetGrid() Grid[TC] {
	return w.Grids[1-w.Bak]
}

func (w *DefaultWorld[TC]) GetBakGrid() Grid[TC] {
	return w.Grids[w.Bak]
}

func (w *DefaultWorld[TC]) SwapBakGrid() {
	w.Bak = 1 - w.Bak
}

func (w *DefaultWorld[TC]) CellNum() int {
	return w.cellnum
}
func (w *DefaultWorld[TC]) WorkNum() int {
	return w.worknum
}
func (w *DefaultWorld[TC]) PageInter() int {
	return w.pageinter
}
func (w *DefaultWorld[TC]) Reciver() WorldReciver[TC] {
	return w.reciver
}

type DefaultGrid[TC Cell] struct {
	Cells   []TC
	CellNum int
}

func (g *DefaultGrid[TC]) GetCell(idex int) TC {
	var tc TC
	if idex >= g.CellNum {
		return tc
	}
	return g.Cells[idex]
}
