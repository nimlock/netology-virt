package main

import "testing"

func TestFindMinInList(t *testing.T) {
	list := []int{1563, 212, 15, 88, 657, 387, 692, 1294}
	var v int
	v = FindMinInList(list)
	if v != 15 {
		t.Error("Expected 15, got ", v)
	}
}
