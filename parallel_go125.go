//go:build go1.25

package parallel

import (
	"sync"
)

// Run executes the given functions in parallel and waits for all to complete.
// It returns when all functions have finished execution.
// This implementation uses WaitGroup.Go available in Go 1.25+.
func Run(fns ...func()) {
	var wg sync.WaitGroup

	for _, fn := range fns {
		wg.Go(fn)
	}

	wg.Wait()
}

// RunWithLimit executes the given functions in parallel with a concurrency limit.
// It ensures that at most 'limit' functions execute concurrently.
// It panics if limit is less than or equal to 0.
// This implementation uses WaitGroup.Go available in Go 1.25+.
func RunWithLimit(limit int, fns ...func()) {
	if limit <= 0 {
		panic("parallel: limit must be greater than 0")
	}

	var wg sync.WaitGroup

	sem := make(chan struct{}, limit)

	for _, fn := range fns {
		sem <- struct{}{}

		wg.Go(func() {
			defer func() {
				<-sem
			}()
			fn()
		})
	}

	wg.Wait()
}

// RunForEach executes the given function for each item in the slice in parallel.
// It waits for all executions to complete before returning.
// This implementation uses WaitGroup.Go available in Go 1.25+.
func RunForEach[T any](items []T, fn func(item T)) {
	var wg sync.WaitGroup

	for _, item := range items {
		wg.Go(func() {
			fn(item)
		})
	}

	wg.Wait()
}

// RunForEachWithLimit executes the given function for each item in the slice in parallel with a concurrency limit.
// It ensures that at most 'limit' functions execute concurrently.
// It panics if limit is less than or equal to 0.
// This implementation uses WaitGroup.Go available in Go 1.25+.
func RunForEachWithLimit[T any](limit int, items []T, fn func(item T)) {
	if limit <= 0 {
		panic("parallel: limit must be greater than 0")
	}

	var wg sync.WaitGroup

	sem := make(chan struct{}, limit)

	for _, item := range items {
		sem <- struct{}{}

		wg.Go(func() {
			defer func() {
				<-sem
			}()
			fn(item)
		})
	}

	wg.Wait()
}
