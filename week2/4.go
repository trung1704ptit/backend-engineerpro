package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type PersonInfo struct {
	Name        string
	Job         string
	YearOfBirth int64
}

func ReadFile(filePath string) ([]PersonInfo, error) {
	var personList []PersonInfo
	f, err := os.Open(filePath)
	defer f.Close()

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		textLine := scanner.Text()
		splittedText := strings.Split(textLine, "|")

		yearOfBirth, err := strconv.ParseInt(splittedText[2], 10, 64)

		if err != nil {
			log.Printf("Error converting YearOfBirth for %s: %v", splittedText[0], err)
			continue
		}

		p := PersonInfo{
			Name:        strings.ToUpper(splittedText[0]),
			Job:         strings.ToLower(splittedText[1]),
			YearOfBirth: yearOfBirth,
		}

		personList = append(personList, p)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return personList, nil
}

func main() {
	path := "week2/a.txt"
	list, e := ReadFile(path)

	if e != nil {
		log.Fatal(e)
	}

	for _, val := range list {
		fmt.Println(val)
	}
}
