# go-parallel

[![CI](https://img.shields.io/github/actions/workflow/status/ventaquil/go-parallel/ci.yml?style=flat-square)](https://github.com/ventaquil/go-parallel/actions/workflows/ci.yml)
[![Go Report Card](https://img.shields.io/badge/Go_Report_Card-View%20report-blue?style=flat-square)](https://goreportcard.com/report/github.com/ventaquil/go-parallel)
[![GoDoc](https://img.shields.io/badge/GoDoc-reference-blue?style=flat-square)](https://godoc.org/github.com/ventaquil/go-parallel)

A simple Go package for executing functions in parallel with optional concurrency limits.

## Requirements

Minimum Go version: 1.18.

## Installation

```bash
go get github.com/ventaquil/go-parallel
```

## Usage

### Basic Parallel Execution

Execute multiple functions in parallel and wait for all to complete:

```go
package main

import (
    "fmt"
    "github.com/ventaquil/go-parallel"
)

func main() {
    parallel.Run(
        func() { fmt.Println("Task 1") },
        func() { fmt.Println("Task 2") },
        func() { fmt.Println("Task 3") },
    )
}
```

### Parallel Execution with Concurrency Limit

Execute functions in parallel but limit the number of concurrent executions:

```go
package main

import (
    "fmt"
    "github.com/ventaquil/go-parallel"
)

func main() {
    // Execute 10 tasks with maximum 3 running concurrently
    tasks := []func(){
        func() { fmt.Println("Task 1") },
        func() { fmt.Println("Task 2") },
        func() { fmt.Println("Task 3") },
        func() { fmt.Println("Task 4") },
        func() { fmt.Println("Task 5") },
    }
    
    parallel.RunWithLimit(3, tasks...)
}
```

## Implementation

The package automatically uses the appropriate implementation based on your Go version:

- Go 1.18-1.24: Implementation using `WaitGroup.Add`/`WaitGroup.Done`
- Go 1.25+: Implementation using `WaitGroup.Go`

## License

This package is licensed under the MIT License.
