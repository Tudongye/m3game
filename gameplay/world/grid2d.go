package world

import (
	"context"
)

type World2d[TC Cell2d] interface {
	World[TC]
	GetCellIndex(x, y int) int
}

type Grid2d[TC Cell2d] interface {
	Grid[TC]
}

type Cell2d interface {
	Cell
	RangeX() (int, int)
	RangeY() (int, int)
	Local(x, y int) bool
}

type WorldReciver2d[TC Cell2d] interface {
	NewCell2d(idex, x1, x2, y1, y2 int) TC
	CellFunc2d(ctx context.Context, p Page, idex int, w World2d[TC]) // 并行阶段
	WorldFunc2d(ctx context.Context, p Page, w World2d[TC])          // 串行阶段
}

func NewWorld2d[TC Cell2d](reciver WorldReciver2d[TC], rangex int, rangey int, cellx int, celly int, worknum int, pageinter int) World2d[TC] {
	m := NewMap2d(rangex, rangey, cellx, celly)
	w := &DefaultWorld2d[TC]{
		m: m,
	}
	worldreciver := &WorldReciver2dInter[TC]{WorldReciver2d: reciver, m: m}
	w.World = NewWorld[TC](worldreciver, m.cellnum, worknum, pageinter)
	return w
}

type DefaultWorld2d[TC Cell2d] struct {
	World[TC]
	m *Map2d
}

func (w *DefaultWorld2d[TC]) GetCellIndex(x, y int) int {
	return w.m.GetCellIndex(x, y)
}

type DefaultGrid2d[TC Cell2d] struct {
	DefaultGrid[TC]
}

func NewCell2d(idex, x1, x2, y1, y2 int) *DefaultCell2d {
	return &DefaultCell2d{idex: idex, x1: x1, x2: x2, y1: y1, y2: y2}
}

type DefaultCell2d struct {
	idex, x1, x2, y1, y2 int
}

func (c *DefaultCell2d) RangeX() (int, int) {
	return c.x1, c.x2
}

func (c *DefaultCell2d) RangeY() (int, int) {
	return c.y1, c.y2
}

func (c *DefaultCell2d) Local(x, y int) bool {
	if x >= c.x1 && x <= c.x2 && y >= c.y1 && y <= c.y2 {
		return true
	}
	return false
}

type WorldReciver2dInter[TC Cell2d] struct {
	WorldReciver2d[TC]
	m *Map2d
}

func (wr *WorldReciver2dInter[TC]) NewCell(idex int) TC {
	x1, x2, y1, y2 := wr.m.GetCellRange(idex)
	return wr.NewCell2d(idex, x1, x2, y1, y2)
}

func (wr *WorldReciver2dInter[TC]) CellFunc(ctx context.Context, p Page, idex int, w World[TC]) {
	wr.CellFunc2d(ctx, p, idex, w.(World2d[TC]))
}

func (wr *WorldReciver2dInter[TC]) WorldFunc(ctx context.Context, p Page, w World[TC]) {
	wr.WorldFunc2d(ctx, p, w.(World2d[TC]))
}

func NewMap2d(rangex, rangey, cellx, celly int) *Map2d {
	cellxnum := (rangex-1)/cellx + 1
	cellynum := (rangey-1)/celly + 1
	cellnum := cellxnum * cellynum
	m := &Map2d{
		rangex:   rangex,
		rangey:   rangey,
		cellx:    cellx,
		celly:    celly,
		cellxnum: cellxnum,
		cellynum: cellynum,
		cellnum:  cellnum,
	}
	return m
}

type Map2d struct {
	rangex, rangey     int // 地图边界
	cellx, celly       int // 单个格子大小
	cellxnum, cellynum int // 横纵格子数
	cellnum            int // 总格子数
}

func (m *Map2d) GetCellIndex(x, y int) int {
	if x < 0 || x >= m.rangex || y < 0 || y >= m.rangey {
		return -1
	}
	xindex := x / m.cellx
	yindex := y / m.celly
	return xindex*m.cellynum + yindex
}

func (m *Map2d) GetCellRange(index int) (int, int, int, int) {
	if index >= m.cellnum {
		return -1, -1, -1, -1
	}
	xindex := index / m.cellynum
	yindex := index - xindex*m.cellynum
	return xindex * m.cellx, min((xindex+1)*m.cellx-1, m.rangex-1), yindex * m.celly, min((yindex+1)*m.celly-1, m.rangey-1)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
