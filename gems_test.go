package main

import (
	"reflect"
	"testing"
)

func TestNewGridOfBools(t *testing.T) {
	tests := []struct {
		name string
		r    int
		cols int
		want Grid[bool]
	}{
		{
			name: "3x3 grid",
			r:    3,
			cols: 3,
			want: Grid[bool]{rows: [][]bool{{false, false, false}, {false, false, false}, {false, false, false}}},
		},
		{
			name: "2x2 grid",
			r:    2,
			cols: 2,
			want: Grid[bool]{rows: [][]bool{{false, false}, {false, false}}},
		},
		{
			name: "0x0 grid",
			r:    0,
			cols: 0,
			want: Grid[bool]{rows: [][]bool{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGridOfBools(tt.r, tt.cols)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGridOfBools() = %v, want %v", got, tt.want)
			}
		})
	}
}
