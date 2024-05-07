package main

/*
	==================================

/* ==================================

	Grid methods

/* ==================================
/* ==================================
*/
type Grid[T comparable] struct {
	rows [][]T
}

type Point struct {
	x, y int
}

func NewGridOfBools(width, height int) Grid[bool] {
	r := make([][]bool, height)
	for i := range r {
		r[i] = make([]bool, width)
	}
	return Grid[bool]{rows: r}
}

func (g *Grid[T]) Copy() Grid[T] {
	clone := make([][]T, len(g.rows))
	for i := range g.rows {
		clone[i] = make([]T, len(g.rows[i]))
		copy(clone[i], g.rows[i])
	}
	return Grid[T]{rows: clone}
}

func (g *Grid[T]) deleteRow(y int) {
	g.rows = append(g.rows[:y], g.rows[y+1:]...)
}

func (g *Grid[T]) ForEach(f func(p Point, value T)) {
	for y := range g.rows {
		for x := range g.rows[y] {
			f(Point{x, y}, g.rows[y][x])
		}
	}
}

func (g *Grid[T]) IndexOf(p Point) int {
	return p.y*g.NumColumns() + p.x
}

func (g *Grid[T]) IndexToPoint(index int) Point {
	return Point{index % g.NumColumns(), index / g.NumColumns()}
}

func (g *Grid[T]) LastIndex() int {
	return g.NumColumns()*g.NumRows() - 1
}

func (g *Grid[T]) SetAll(value T) {
	for y := range g.rows {
		for x := range g.rows[y] {
			g.rows[y][x] = value
		}
	}
}

func (g *Grid[T]) GetPtr(p Point) *T {
	return &g.rows[p.y][p.x]
}

func (g *Grid[T]) Get(p Point) T {
	if p.y < 0 || p.y >= g.NumRows() || p.x < 0 || p.x >= g.NumColumns() {
		panic("Index out of bounds")
	}
	return g.rows[p.y][p.x]
}

func (g *Grid[T]) GetAtIndex(index int) T {
	return g.rows[index/g.WidthOf()][index%g.WidthOf()]
}

func (g *Grid[T]) NumColumns() int {
	return len(g.rows[0])
}

func (g *Grid[T]) NumRows() int {
	return len(g.rows)
}

func (g *Grid[T]) WidthOf() int {
	return len(g.rows[0])
}

func (g *Grid[T]) Set(p Point, value T) {
	g.rows[p.y][p.x] = value
}

func (g *Grid[T]) SetAtIndex(index int, value T) {
	g.rows[index/g.WidthOf()][index%g.WidthOf()] = value
}
