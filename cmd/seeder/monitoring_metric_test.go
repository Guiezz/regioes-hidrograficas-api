package main

import "testing"

func TestHasMetric(t *testing.T) {
	cases := map[string]bool{
		"":    false,
		"   ": false,
		"0":   true,
		"0,5": true,
		"1":   true,
		"50%": true,
	}
	for in, want := range cases {
		if got := hasMetric(in); got != want {
			t.Errorf("hasMetric(%q) = %v, want %v", in, got, want)
		}
	}
}
