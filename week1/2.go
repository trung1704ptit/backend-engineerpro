package main

import "fmt"

func checkLength(s string) bool {
	return len(s)%2 == 0
}

func Week1Ex2() {
	var s string

	fmt.Print("Enter string: ")
	fmt.Scan(&s)

	fmt.Println(checkLength(s))
}
