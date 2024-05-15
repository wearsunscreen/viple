// Package grid provides a generic type that represents a two-dimensional slice.
package main

// Grid is a generic type that represents a two-dimensional slice.
// It is parameterized by the type of elements it contains, denoted by T.
// Elements can be retrieved and set using x,y coordinates or by index into a one-dimensional array.
type Grid[T comparable] [][]T

// type Grid[T comparable] struct {
// 	rows [][]T
// }

// Coord represents a coordinate in a grid.
type Coord struct {
	x, y int
}

// To do:
// Add generic function to create a new grid of any type
// by passing a

// To do:
// redefine as a 2D slice
// type Grid[T comparable] [][]T

// NewGridOfBools creates a new Grid of booleans with the specified width and height.
func NewGridOfBools(width, height int) Grid[bool] {
	r := make([][]bool, height)
	for i := range r {
		r[i] = make([]bool, width)
	}
	return r
}

// Copy creates a deep copy of the Grid.
func (g Grid[T]) Copy() Grid[T] {
	clone := make([][]T, len(g))
	for i := range g {
		clone[i] = make([]T, len(g[i]))
		copy(clone[i], g[i])
	}
	return clone
}

// return a copy of the grid with a row removed.
func (g Grid[T]) DeleteRow(y int) Grid[T] {
	return append(g[:y], g[y+1:]...)
}

// ForEach applies the given function to each element in the Grid.
func (g Grid[T]) ForEach(f func(p Coord, value T)) {
	for y := range g {
		for x := range g[y] {
			f(Coord{x, y}, g[y][x])
		}
	}
}

// IndexOf returns the index of the element at the specified Coord in the Grid.
func (g Grid[T]) IndexOf(p Coord) int {
	return p.y*g.NumColumns() + p.x
}

// IndexToCoord converts the given index to a Coord in the Grid.
func (g Grid[T]) IndexToCoord(index int) Coord {
	return Coord{index % g.NumColumns(), index / g.NumColumns()}
}

// LastIndex returns the index of the last element in the Grid.
func (g Grid[T]) LastIndex() int {
	return g.NumColumns()*g.NumRows() - 1
}

// SetAll sets all elements in the Grid to the specified value.
func (g Grid[T]) SetAll(value T) {
	for y := range g {
		for x := range g[y] {
			g[y][x] = value
		}
	}
}

// GetPtr returns a pointer to the element at the specified Coord in the Grid.
func (g Grid[T]) GetPtr(p Coord) *T {
	return &g[p.y][p.x]
}

// Get returns the element at the specified Point in the Grid.
// It panics if the Point is out of bounds.
func (g Grid[T]) Get(p Coord) T {
	if p.y < 0 || p.y >= g.NumRows() || p.x < 0 || p.x >= g.NumColumns() {
		panic("Index out of bounds")
	}
	return g[p.y][p.x]
}

// GetAtIndex returns the element at the specified index in the Grid.
func (g Grid[T]) GetAtIndex(index int) T {
	return g[index/g.NumColumns()][index%g.NumColumns()]
}

// NumColumns returns the number of columns in the Grid.
func (g Grid[T]) NumColumns() int {
	return len(g[0])
}

// NumRows returns the number of rows in the Grid.
func (g Grid[T]) NumRows() int {
	return len(g)
}

// Set sets the element at the specified Coord in the Grid to the specified value.
func (g Grid[T]) Set(p Coord, value T) {
	g[p.y][p.x] = value
}

// SetAtIndex sets the element at the specified index in the Grid to the specified value.
func (g Grid[T]) SetAtIndex(index int, value T) {
	g[index/g.NumColumns()][index%g.NumColumns()] = value
}
