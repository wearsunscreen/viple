package main

import "testing"

func TestLimitToRange(t *testing.T) {
	tests := []struct {
		input, low, high, want int
	}{
		{5, 1, 10, 5},   // input is within the range
		{0, 1, 10, 1},   // input is below the range
		{15, 1, 10, 10}, // input is above the range
	}

	for _, tt := range tests {
		got := limitToRange(tt.input, tt.low, tt.high)
		if got != tt.want {
			t.Errorf("limitToRange(%v, %v, %v) = %v; want %v", tt.input, tt.low, tt.high, got, tt.want)
		}
	}
}
