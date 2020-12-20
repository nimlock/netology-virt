package main

import "testing"

func TestMetersToFt(t *testing.T) {
	var v float64
	v = MetersToFt(10)
	if v != 32.8084 {
		t.Error("Expected 32.8084, got ", v)
	}
}
