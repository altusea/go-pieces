# Go Pieces

A simple Go project for learning Go programming basics.

## Project Structure

```
go-pieces/
├── cmd/           # Command line applications
├── internal/      # Internal packages
├── pkg/           # Public packages
│   └── math/      # Math utility package
├── main.go        # Main application entry point
├── go.mod         # Go module file
├── .gitignore     # Git ignore file
└── README.md      # This file
```

## Features

- Basic math utility functions (Add, Subtract, Multiply, Divide)
- Error handling for division by zero
- Clean project structure following Go conventions

## Getting Started

### Prerequisites

- Go installed on your system

### Running the Application

```bash
# Run the main application
go run main.go

# Build the application
go build -o go-pieces main.go

# Run tests
go test ./...
```

## Usage Example

```go
package main

import (
    "fmt"
    "go-pieces/pkg/math"
)

func main() {
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
}
```

## Math Package

The `pkg/math` package provides basic arithmetic operations:

- `Add(a, b int) int` - Addition
- `Subtract(a, b int) int` - Subtraction
- `Multiply(a, b int) int` - Multiplication
- `Divide(a, b int) (int, error)` - Division with error handling

## Contributing

This is a learning project. Feel free to explore and modify the code to understand Go programming concepts.