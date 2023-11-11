package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"testing"
	"time"
)

// Task 1: Concurrency and Goroutines

func factorial(n int, resultChan chan int) {
	result := 1
	for i := 1; i <= n; i++ {
		result *= i
	}
	resultChan <- result
}

// Task 2: Interfaces and Polymorphism

type Shape interface {
	CalculateArea() float64
	CalculatePerimeter() float64
}

type Rectangle struct {
	width  float64
	height float64
}

type Circle struct {
	radius float64
}

func (rec Rectangle) CalculatePerimeter() float64 {
	return 2 * (rec.width + rec.height)
}

func (rec Rectangle) CalculateArea() float64 {
	return rec.width * rec.height
}

func (circ Circle) CalculateArea() float64 {
	return math.Pi * circ.radius * circ.radius
}

func (circ Circle) CalculatePerimeter() float64 {
	return 2 * math.Pi * circ.radius
}

func calculateShapeInfo(s Shape) {
	area := s.CalculateArea()
	perimeter := s.CalculatePerimeter()

	fmt.Printf("Area: %f, Perimeter: %f\n", area, perimeter)
}

// Task 3: Error Handling with Custom Errors

type FileNotFoundError struct {
	FileName string
}

func (e FileNotFoundError) Error() string {
	return fmt.Sprintf("File not found: %s", e.FileName)
}

func simulateFileOperation(fileName string) error {
	if fileName == "nonexistent.txt" {
		return FileNotFoundError{FileName: fileName}
	}

	return nil
}

// Task 4: File I/O and JSON

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Pages  int    `json:"pages"`
}

func readAndProcessJSON(filePath string) error {
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var books []Book
	err = json.Unmarshal(jsonData, &books)
	if err != nil {
		return err
	}

	// Perform operations on the data (e.g., filter, sort)

	// Write the modified data back to a new JSON file
	newJSONData, err := json.MarshalIndent(books, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("modified_books.json", newJSONData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Task 5: Testing and Benchmarking

func TestFactorial(t *testing.T) {
	testCases := []struct {
		input    int
		expected int
	}{
		{0, 1},
		{1, 1},
		{5, 120},
		{10, 3628800},
	}

	for _, testCase := range testCases {
		resultChan := make(chan int)
		go factorial(testCase.input, resultChan)
		result := <-resultChan
		close(resultChan)

		if result != testCase.expected {
			t.Errorf("For input %d, expected %d, but got %d", testCase.input, testCase.expected, result)
		}
	}
}

func BenchmarkFactorial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultChan := make(chan int)
		go factorial(10, resultChan)
		<-resultChan
		close(resultChan)
	}
}

func BenchmarkFactorialLargeInput(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultChan := make(chan int)
		go factorial(20, resultChan)
		<-resultChan
		close(resultChan)
	}
}

func BenchmarkFactorialParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resultChan := make(chan int)
			go factorial(10, resultChan)
			<-resultChan
			close(resultChan)
		}
	})
}

func main() {
	// Task 1
	startTime := time.Now()

	num := 5
	resultChan := make(chan int)

	go factorial(num, resultChan)

	result := <-resultChan
	close(resultChan)

	elapsedTime := time.Since(startTime)

	fmt.Printf("Factorial of %d is %d\n", num, result)
	fmt.Printf("Execution time: %s\n", elapsedTime)

	// Task 2
	rect := Rectangle{width: 5, height: 3}
	circ := Circle{radius: 4}

	calculateShapeInfo(rect)
	calculateShapeInfo(circ)

	// Task 3
	fileName := "nonexistent.txt"
	err := simulateFileOperation(fileName)
	if err != nil {
		if fileNotFoundError, ok := err.(FileNotFoundError); ok {
			fmt.Println("Custom Error:", fileNotFoundError.Error())
		} else {
			fmt.Println("Generic Error:", err.Error())
		}
		return
	}

	// Run tests using the 'go test' command
	cmd := exec.Command("go", "test")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Tests Failed!")
		return
	}

	// Print a separator for clarity
	fmt.Println("---------------------------------------------------")

	// Run benchmarks
	fmt.Println("Running Benchmarks:")
	cmd = exec.Command("go", "test", "-bench=.")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Benchmarks Failed!")
		return
	}

	fmt.Println("Benchmarks Passed!")
}
