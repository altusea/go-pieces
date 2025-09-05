package main

import (
	"fmt"
	"go-pieces/pkg/math"
)

func main() {
	fmt.Println("Hello, World!")

	sum := math.Add(10, 5)
	diff := math.Subtract(10, 5)
	product := math.Multiply(10, 5)
	quotient, err := math.Divide(10, 5)

	fmt.Printf("10 + 5 = %d\n", sum)
	fmt.Printf("10 - 5 = %d\n", diff)
	fmt.Printf("10 * 5 = %d\n", product)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("10 / 5 = %d\n", quotient)
	}

	_, err = math.Divide(10, 0)
	if err != nil {
		fmt.Printf("10 / 0 = Error: %v\n", err)
	}
}
