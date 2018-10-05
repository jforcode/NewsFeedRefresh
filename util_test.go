package main

import "testing"

func TestMin(t *testing.T) {
	util := Util{}
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"normal", 5, 20, 5},
		{"normal diff order", 20, 5, 5},
		{"negative", 56, -20, -20},
		{"zero", 0, 23, 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := util.MinInt(test.a, test.b)
			if !(test.expected == actual) {
				t.FailNow()
			}
		})
	}
}
