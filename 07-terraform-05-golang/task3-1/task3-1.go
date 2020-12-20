package main

import (
	"fmt"
	"math"
	"os"
)

func main() {
	fmt.Println("Enter a number of meters to convert into ft: ")
	var input float64
	fmt.Fscan(os.Stdin, &input)

	fmt.Printf("Ok, it equal to %.4f ft\n", MetersToFt(input))
}

func MetersToFt(meters float64) float64 {
	coeff := math.Pow(0.3048, -1)
	result := math.Round(meters*coeff*10000) / 10000
	return result
}
