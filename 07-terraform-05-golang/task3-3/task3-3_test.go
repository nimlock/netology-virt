package main

import (
	"reflect"
	"testing"
)

func TestClearDivisionBy3(t *testing.T) {
	// var v float64
	v := ClearDivisionBy3(10, 20)
	rightAnswer := []int{12, 15, 18}
	if !(reflect.DeepEqual(v, rightAnswer)) {
		t.Error("Expected ", rightAnswer, " got ", v)
	}
}
