package main

import (
	"fmt"
)

const (
	UintSize = 32 << (^uint(0) >> 32 & 1)
	MaxInt   = 1<<(UintSize-1) - 1 // Определяем наибольшее значение в int
)

func main() {
	x := []int{48, 96, 86, 68, 57, 82, 63, 70, 37, 34, 83, 27, 19, 97, 9, 17, 555}
	fmt.Printf("Max value in array/split is %d\n", FindMinInList(x))
}

func FindMinInList(list []int) int {
	minValue := MaxInt
	for _, value := range list {
		if value < minValue {
			minValue = value
		}
	}
	return minValue
}
