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
    err := parallel.Run(
        func() error { fmt.Println("Task 1"); return nil },
        func() error { fmt.Println("Task 2"); return nil },
        func() error { fmt.Println("Task 3"); return nil },
    )
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
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
    // Execute tasks with maximum 3 running concurrently
    tasks := []func() error{
        func() error { fmt.Println("Task 1"); return nil },
        func() error { fmt.Println("Task 2"); return nil },
        func() error { fmt.Println("Task 3"); return nil },
        func() error { fmt.Println("Task 4"); return nil },
        func() error { fmt.Println("Task 5"); return nil },
    }
    
    err := parallel.RunLimit(3, tasks...)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

### Process Slice Items in Parallel

Execute a function for each item in a slice in parallel:

```go
package main

import (
    "fmt"
    "github.com/ventaquil/go-parallel"
)

func main() {
    items := []int{1, 2, 3, 4, 5}
    
    err := parallel.RunForEach(items, func(item int) error {
        fmt.Printf("Processing item: %d\n", item)
        return nil
    })
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

### Process Slice Items with Concurrency Limit

Execute a function for each item with controlled concurrency:

```go
package main

import (
    "fmt"
    "github.com/ventaquil/go-parallel"
)

func main() {
    items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    // Process items with maximum 3 running concurrently
    err := parallel.RunForEachLimit(3, items, func(item int) error {
        fmt.Printf("Processing item: %d\n", item)
        return nil
    })
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

### Context-Aware Execution with Error Handling

Execute functions with context support and error handling:

```go
package main

import (
    "context"
    "fmt"
    "github.com/ventaquil/go-parallel"
)

func main() {
    ctx := context.Background()
    
    err := parallel.RunContext(ctx,
        func(ctx context.Context) error {
            fmt.Println("Task 1")
            return nil
        },
        func(ctx context.Context) error {
            fmt.Println("Task 2")
            return nil
        },
    )
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

### Process Items with Context and Concurrency Limit

Combine concurrency control with context and error handling:

```go
package main

import (
    "context"
    "fmt"
    "github.com/ventaquil/go-parallel"
)

func main() {
    ctx := context.Background()
    items := []int{1, 2, 3, 4, 5}
    
    err := parallel.RunForEachLimitContext(ctx, 3, items,
        func(ctx context.Context, item int) error {
            fmt.Printf("Processing item: %d\n", item)
            return nil
        },
    )
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## API

### `Run(fns ...func() error) error`

Executes the given functions in parallel and waits for all to complete. Returns the first error encountered, if any.

### `RunLimit(limit int, fns ...func() error) error`

Executes the given functions in parallel with a concurrency limit. Ensures that at most `limit` functions execute concurrently. Returns an error if limit is less than or equal to 0.

### `RunForEach[T any](items []T, fn func(item T) error) error`

Executes the given function for each item in the slice in parallel. Returns the first error encountered, if any.

### `RunForEachLimit[T any](limit int, items []T, fn func(item T) error) error`

Executes the given function for each item in the slice in parallel with a concurrency limit. Ensures that at most `limit` functions execute concurrently. Returns an error if limit is less than or equal to 0.

### `RunContext(ctx context.Context, fns ...func(ctx context.Context) error) error`

Executes the given functions in parallel with context support. Returns the first error encountered, if any, and cancels the context for remaining functions. Only the first error is returned; subsequent errors are discarded.

### `RunLimitContext(ctx context.Context, limit int, fns ...func(ctx context.Context) error) error`

Executes the given functions in parallel with a concurrency limit and context support. Returns an error if limit is less than or equal to 0. Returns the first error encountered from any function.

### `RunForEachContext[T any](ctx context.Context, items []T, fn func(ctx context.Context, item T) error) error`

Executes the given function for each item in the slice in parallel with context support. Returns the first error encountered, if any, and cancels the context for remaining functions.

### `RunForEachLimitContext[T any](ctx context.Context, limit int, items []T, fn func(ctx context.Context, item T) error) error`

Executes the given function for each item in the slice in parallel with a concurrency limit and context support. Returns an error if limit is less than or equal to 0. Returns the first error encountered from any function.

## Implementation

The package automatically uses the appropriate implementation based on your Go version:

- Go 1.18-1.24: Implementation using `WaitGroup.Add`/`WaitGroup.Done`
- Go 1.25+: Implementation using `WaitGroup.Go`

## License

This package is licensed under the MIT License.
