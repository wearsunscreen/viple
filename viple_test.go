package main

import (
	"reflect"
	"testing"
)

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

func TestRemoveDuplicatesOf(t *testing.T) {
	tests := []struct {
		name  string
		in    []int
		value int
		want  []int
	}{
		{
			name:  "No duplicates",
			in:    []int{1, 2, 3, 4, 5},
			value: 4,
			want:  []int{1, 2, 3, 4, 5},
		},
		{
			name:  "No duplicates, value not in slice",
			in:    []int{1, 2, 3, 4, 5},
			value: 9,
			want:  []int{1, 2, 3, 4, 5},
		},
		{
			name:  "With duplicater removing all but the the first 4 elements",
			value: 4,
			in:    []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 5},
			want:  []int{1, 2, 2, 3, 3, 3, 4, 5, 5, 5, 5, 5},
		},
		{
			name:  "With duplicates/unsorted, removing all but the the first 4 elements",
			in:    []int{4, 4, 5, 2, 3, 1, 2, 3, 3, 4, 4, 5, 5, 5, 5},
			value: 4,
			want:  []int{4, 5, 2, 3, 1, 2, 3, 3, 5, 5, 5, 5},
		},
		{
			name:  "All duplicates",
			in:    []int{1, 1, 1, 1, 1},
			value: 1,
			want:  []int{1},
		},
		{
			name:  "Empty slice",
			in:    []int{},
			value: 1,
			want:  []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			removeDuplicatesOf(&tt.in, tt.value)
			if !reflect.DeepEqual(tt.in, tt.want) {
				t.Errorf("removeDuplicatesOf() = %v, want %v", tt.in, tt.want)
			}
		})
	}
}

func TestTail(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		n    int
		want []int
	}{
		{
			name: "Get last 3 elements",
			in:   []int{1, 2, 3, 4, 5},
			n:    3,
			want: []int{3, 4, 5},
		},
		{
			name: "Get more elements than exist",
			in:   []int{1, 2, 3, 4, 5},
			n:    10,
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "Get no elements",
			in:   []int{1, 2, 3, 4, 5},
			n:    0,
			want: []int{},
		},
		{
			name: "Nil slice",
			in:   nil,
			n:    3,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tail(tt.in, tt.n)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tail() = %v, want %v", got, tt.want)
			}
		})
	}
}
