package main

import "testing"

func TestIsPointInRect(t *testing.T) {
	tests := []struct {
		pointX, pointY, rectLeft, rectTop, rectWidth, rectHeight float32
		want                                                     bool
	}{
		{5, 5, 0, 0, 10, 10, true},    // point is inside the rectangle
		{15, 15, 0, 0, 10, 10, false}, // point is outside the rectangle
		{0, 0, 0, 0, 10, 10, false},   // point is on the top-left corner of the rectangle
		{10, 10, 0, 0, 10, 10, false}, // point is on the bottom-right corner of the rectangle
	}

	for _, tt := range tests {
		got := isPointInRect(tt.pointX, tt.pointY, tt.rectLeft, tt.rectTop, tt.rectWidth, tt.rectHeight)
		if got != tt.want {
			t.Errorf("isPointInRect(%v, %v, %v, %v, %v, %v) = %v; want %v", tt.pointX, tt.pointY, tt.rectLeft, tt.rectTop, tt.rectWidth, tt.rectHeight, got, tt.want)
		}
	}
}

func TestAnyIn2DSlice(t *testing.T) {
	tests := []struct {
		input [][]bool
		want  bool
	}{
		{[][]bool{{false, false}, {false, false}}, false}, // no true values
		{[][]bool{{true, false}, {false, false}}, true},   // one true value
		{[][]bool{{true, true}, {true, true}}, true},      // all true values
	}

	for _, tt := range tests {
		got := anyIn2DSlice(tt.input)
		if got != tt.want {
			t.Errorf("anyIn2DSlice(%v) = %v; want %v", tt.input, got, tt.want)
		}
	}
}

func TestAnyInSlice(t *testing.T) {
	tests := []struct {
		input []bool
		want  bool
	}{
		{[]bool{false, false, false, false}, false}, // no true values
		{[]bool{true, false, false, false}, true},   // one true value
		{[]bool{true, true, true, true}, true},      // all true values
	}

	for _, tt := range tests {
		got := anyInSlice(tt.input)
		if got != tt.want {
			t.Errorf("anyInSlice(%v) = %v; want %v", tt.input, got, tt.want)
		}
	}
}
