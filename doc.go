/*
Package parallel provides utilities for executing functions concurrently with optional concurrency control.

# Overview

The parallel package offers a simple and efficient way to execute multiple functions
concurrently in Go. It provides eight main functions:

  - Run: Execute functions in parallel without any concurrency limit, returning the first error
  - RunContext: Execute functions in parallel with context support and error handling
  - RunLimit: Execute functions in parallel with a maximum concurrency limit
  - RunLimitContext: Execute functions with concurrency limit, context support and error handling
  - RunForEach: Execute a function for each item in a slice in parallel, returning the first error
  - RunForEachContext: Execute a function for each item with context support and error handling
  - RunForEachLimit: Execute a function for each item with a concurrency limit, returning the first error
  - RunForEachLimitContext: Execute a function for each item with concurrency limit, context and error handling

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

	err := parallel.Run(
	    func() error { fmt.Println("Task 1"); return nil },
	    func() error { fmt.Println("Task 2"); return nil },
	    func() error { fmt.Println("Task 3"); return nil },
	)

Execute tasks with a concurrency limit of 3:

	err := parallel.RunLimit(3,
	    func() error { return processTask1() },
	    func() error { return processTask2() },
	    func() error { return processTask3() },
	    func() error { return processTask4() },
	    func() error { return processTask5() },
	)

Execute a function for each item in a slice:

	items := []int{1, 2, 3, 4, 5}
	err := parallel.RunForEach(items, func(item int) error {
	    return process(item)
	})

Execute a function for each item with a concurrency limit:

	err := parallel.RunForEachLimit(3, items, func(item int) error {
	    return process(item)
	})

Execute functions with context and error handling:

	ctx := context.Background()
	err := parallel.RunContext(ctx,
	    func(ctx context.Context) error {
	        return processTask1(ctx)
	    },
	    func(ctx context.Context) error {
	        return processTask2(ctx)
	    },
	)

Execute functions for each item with context and error handling:

	ctx := context.Background()
	err := parallel.RunForEachContext(ctx, items, func(ctx context.Context, item int) error {
	    return processItem(ctx, item)
	})

# Error Handling

All functions return an error if any of the executed functions returns an error.
The first error encountered is returned; subsequent errors are discarded.

Functions that accept a limit parameter also return an error if the limit is
less than or equal to 0:

	// These will return an error
	err := parallel.RunLimit(0, tasks...)
	err := parallel.RunForEachLimit(0, items, fn)
	err := parallel.RunLimitContext(ctx, 0, tasks...)
	err := parallel.RunForEachLimitContext(ctx, 0, items, fn)

The context-aware functions (RunContext, RunLimitContext, RunForEachContext,
and RunForEachLimitContext) additionally cancel the context when an error occurs:

  - When an error occurs, the context is cancelled for remaining functions
  - Only the first error is returned; subsequent errors are discarded
  - RunLimitContext and RunForEachLimitContext return an error if limit <= 0

# Performance Considerations

The package uses sync.WaitGroup for coordination and a semaphore pattern (buffered channel)
for concurrency limiting. On Go 1.25+, it uses the WaitGroup.Go method.

For optimal performance:
  - Use Run for CPU-bound tasks where you want maximum parallelism
  - Use RunLimit for I/O-bound tasks or when you need to control resource usage
  - Use RunForEach when you need to process a slice of items in parallel
  - Use RunForEachLimit when processing items with controlled concurrency
  - Use context-aware functions (RunContext, RunForEachContext, etc.) when you need cancellation support
  - The limit in RunLimit and RunForEachLimit should typically match your expected concurrency requirements
*/

package parallel
