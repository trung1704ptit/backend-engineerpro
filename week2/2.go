package main

func stringCount(input string) map[rune]int {
	strMap := make(map[rune]int)

	for _, char := range input {
		strMap[char]++
	}

	return strMap
}

// func main() {
// 	var val string

// 	fmt.Print("Input the string: ")

// 	fmt.Scanln(&val)

// 	strMap := stringCount(val)

// 	for char, count := range strMap {
// 		fmt.Printf("%c: %d\n", char, count)
// 	}
// }
