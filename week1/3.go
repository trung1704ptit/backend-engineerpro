package main

import (
	"fmt"
	"sort"
)

func getSum(slice []float64) float64 {
	sum := 0.0
	for _, val := range slice {
		sum += val
	}

	return sum
}

func getMax(slice []float64) float64 {
	max := slice[0]
	for _, val := range slice {
		if val > max {
			max = val
		}
	}
	return max
}

func getMin(slice []float64) float64 {
	min := slice[0]
	for _, val := range slice {
		if val < min {
			min = val
		}
	}
	return min
}

func getAverage(slice []float64) float64 {
	sum := getSum(slice)

	return sum / float64(len(slice))
}

func getSorted(slice []float64) []float64 {
	sort.Float64s(slice)
	return slice
}

func Week1Ex3() {
	var n int
	var slice []float64

	fmt.Println("Input number of slice: ")
	fmt.Scan(&n)

	for i := 0; i < n; i++ {
		var val float64
		fmt.Printf("Input element %d: ", i+1)
		fmt.Scan(&val)
		slice = append(slice, val)
	}

	fmt.Println("Sum: ", getSum(slice))
	fmt.Println("Max number: ", getMax(slice))
	fmt.Println("Min number: ", getMin(slice))
	fmt.Println("Average: ", getAverage(slice))
	fmt.Println("Sorted: ", getSorted(slice))
}
