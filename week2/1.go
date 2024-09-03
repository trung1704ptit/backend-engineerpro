package main

import (
	"time"
)

type Person struct {
	Name        string
	Job         string
	YearOfBirth int64
}

func (p *Person) calculateAge() int {
	return time.Now().Year() - int(p.YearOfBirth)
}

func (p *Person) checkJobFit() bool {
	nameLength := len(p.Name)
	if nameLength == 0 {
		return false
	}
	return p.YearOfBirth%int64(len(p.Name)) == 0
}

// func main() {
// 	person := Person{Name: "Trung", Job: "Engineer", YearOfBirth: 1994}

// 	fmt.Printf("%s's age: %d\n", person.Name, person.calculateAge())
// 	fmt.Printf("Is %s's job fit: %v\n", person.Name, person.checkJobFit())
// }
