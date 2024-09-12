// https://dummy.restapiexample.com/api/v1/employees

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Employee struct {
	ID             string `json:"id"`
	EmployeeName   string `json:"employee_name"`
	EmployeeSalary int64  `json:"employee_salary"`
	EmployeeAge    int64  `json:"employee_age"`
	ProfileImage   string `json:"profile_image"`
}

// struct to hold full API response
type ApiResponse struct {
	Status  string     `json:"status"`
	Data    []Employee `json:"data"`
	Message string     `json:"message"`
}

func queryAPI() ([]Employee, error) {
	url := "https://dummy.restapiexample.com/api/v1/employees"

	// make request
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error Request: %v", err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error Reading body %v", err)
	}

	fmt.Println(body)

	// Unmarshal json into ApiResponse struct
	var apiResponse ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("error Unmarshal json %v", err)
	}

	return apiResponse.Data, nil
}

func worker(id int, jobs <-chan Employee, wg *sync.WaitGroup) {
	for employee := range jobs {
		time.Sleep(1 * time.Second)
		averageSalary := employee.EmployeeSalary / employee.EmployeeAge

		fmt.Printf("Worker %d with employee %s has average salary %d\n", id, employee.EmployeeName, averageSalary)
		wg.Done()
	}
}

func addJobWorkerPool(employees []Employee) {
	var wg sync.WaitGroup

	jobs := make(chan Employee, len(employees))

	for i := 1; i < 10; i++ {
		go worker(i, jobs, &wg)
	}

	for _, employee := range employees {
		wg.Add(1)
		jobs <- employee
	}

	close(jobs)

	wg.Wait()

	fmt.Println("We are done!")
}

func main() {
	employees, err := queryAPI()
	if err != nil {
		log.Fatalf("Failed to get employees: %v", err)
	}

	// Print out each employee's data
	for _, employee := range employees {
		fmt.Printf("ID: %s, Name: %s, Salary: %s, Age: %s\n", employee.ID, employee.EmployeeName, employee.EmployeeSalary, employee.EmployeeAge)
	}

	addJobWorkerPool(employees)
}
