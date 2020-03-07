package main

import (
	"fmt"

	"demo.com/pkg/v1/math"
)

func main() {
	nums := []float64{4.62, 90.31, 18.4, 70, 498}
	avg := math.Average(nums)
	fmt.Println(avg)
}
