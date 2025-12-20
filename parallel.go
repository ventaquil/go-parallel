//go:build !go1.25

package parallel

import (
	"errors"
	"sync"
)

// ErrInvalidLimit is returned when a limit parameter is less than or equal to 0.
var ErrInvalidLimit = errors.New("limit must be greater than 0")

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

// RunLimit executes the given functions in parallel with a concurrency limit.
// It ensures that at most 'limit' functions execute concurrently.
// It returns an error if limit is less than or equal to 0.
func RunLimit(limit int, fns ...func()) error {
	if limit <= 0 {
		return ErrInvalidLimit
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

	return nil
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

// RunForEachLimit executes the given function for each item in the slice in parallel with a concurrency limit.
// It ensures that at most 'limit' functions execute concurrently.
// It returns an error if limit is less than or equal to 0.
func RunForEachLimit[T any](limit int, items []T, fn func(item T)) error {
	if limit <= 0 {
		return ErrInvalidLimit
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

	return nil
}
