/*
Package parallel provides utilities for executing functions concurrently with optional concurrency control.

# Overview

The parallel package offers a simple and efficient way to execute multiple functions
concurrently in Go. It provides two main functions:

  - Run: Execute functions in parallel without any concurrency limit
  - RunWithLimit: Execute functions in parallel with a maximum concurrency limit

# Requirements

The package requires Go 1.18 or later and has zero external dependencies.
It automatically selects the appropriate implementation based on your Go version:

  - Go 1.18-1.24: Implementation with manual goroutine management
  - Go 1.25+: Implementation using WaitGroup.Go

The API remains identical regardless of which Go version you use, ensuring
seamless compatibility.

# Implementation Details

The package uses Go build tags to provide version-specific implementations:

  - parallel.go: Implementation for Go 1.18-1.24 (build tag: !go1.25)
  - parallel_go125.go: Implementation for Go 1.25+ (build tag: go1.25)

Build tags ensure the correct implementation is automatically selected at compile time,
with no runtime overhead or version detection required.

# Features

  - Thread-safe execution with proper synchronization primitives
  - Zero external dependencies - uses only Go standard library
  - Automatic version-specific optimization via build tags
  - Simple and intuitive API
  - Comprehensive test coverage

# Basic Usage

Execute multiple tasks in parallel without any limit:

	parallel.Run(
	    func() { fmt.Println("Task 1") },
	    func() { fmt.Println("Task 2") },
	    func() { fmt.Println("Task 3") },
	)

Execute tasks with a concurrency limit of 3:

	parallel.RunWithLimit(3,
	    func() { processTask1() },
	    func() { processTask2() },
	    func() { processTask3() },
	    func() { processTask4() },
	    func() { processTask5() },
	)

# Error Handling

The RunWithLimit function will panic if the limit parameter is less than or equal to 0.
This is a programming error that should be caught during development:

	// This will panic
	parallel.RunWithLimit(0, tasks...)

# Performance Considerations

The package uses sync.WaitGroup for coordination and a semaphore pattern (buffered channel)
for concurrency limiting. On Go 1.25+, it uses the WaitGroup.Go method.

For optimal performance:
  - Use Run for CPU-bound tasks where you want maximum parallelism
  - Use RunWithLimit for I/O-bound tasks or when you need to control resource usage
  - The limit in RunWithLimit should typically match your expected concurrency requirements
*/

package parallel
