//go:build !go1.25

package parallel

import (
	"sync"
)

// Run executes the given functions in parallel and waits for all to complete.
// It returns when all functions have finished execution.
func Run(fns ...func()) {
	var wg sync.WaitGroup
	wg.Add(len(fns))
	for _, fn := range fns {
		go func(f func()) {
			defer wg.Done()
			f()
		}(fn)
	}
	wg.Wait()
}

// RunWithLimit executes the given functions in parallel with a concurrency limit.
// It ensures that at most 'limit' functions execute concurrently.
// It panics if limit is less than or equal to 0.
func RunWithLimit(limit int, fns ...func()) {
	if limit <= 0 {
		panic("parallel: limit must be greater than 0")
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, limit)
	wg.Add(len(fns))
	for _, fn := range fns {
		sem <- struct{}{}

		go func(f func()) {
			defer func() {
				<-sem
				wg.Done()
			}()
			f()
		}(fn)
	}
	wg.Wait()
}

// RunForEach executes the given function for each item in the slice in parallel.
// It waits for all executions to complete before returning.
func RunForEach[T any](items []T, fn func(item T)) {
	var wg sync.WaitGroup
	wg.Add(len(items))
	for _, item := range items {
		go func(it T) {
			defer wg.Done()
			fn(it)
		}(item)
	}
	wg.Wait()
}

// RunForEachWithLimit executes the given function for each item in the slice in parallel with a concurrency limit.
// It ensures that at most 'limit' functions execute concurrently.
// It panics if limit is less than or equal to 0.
func RunForEachWithLimit[T any](limit int, items []T, fn func(item T)) {
	if limit <= 0 {
		panic("parallel: limit must be greater than 0")
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, limit)
	wg.Add(len(items))
	for _, item := range items {
		sem <- struct{}{}

		go func(it T) {
			defer func() {
				<-sem
				wg.Done()
			}()
			fn(it)
		}(item)
	}
	wg.Wait()
}
