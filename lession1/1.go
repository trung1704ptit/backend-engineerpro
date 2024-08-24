package main

import "fmt"

func perimeter(x, y float64) float64 {
	return (x + y) * 2
}

func area(x, y float64) float64 {
	return x * y
}

func Ex1() {
	var x, y float64

	fmt.Print("Enter value x: ")
	fmt.Scan(&x)

	fmt.Print("Enter value y: ")
	fmt.Scan(&y)

	fmt.Println("Rectangle perimeter: ", perimeter(x, y))
	fmt.Println("Rectangle area: ", area(x, y))
}
