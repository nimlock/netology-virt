package main

import (
	"fmt"
)

func main() {
	start, end := 0, 100
	fmt.Printf("Let`s find all numbers from %d and %d that clear division by 3\n", start, end)
	fmt.Println(ClearDivisionBy3(start, end))
}

func ClearDivisionBy3(start, end int) []int {
	result := make([]int, 0, 0)
	for i := start; i < end; i++ {
		if i%3 == 0 {
			result = append(result, i)
		}
	}
	return result
}
