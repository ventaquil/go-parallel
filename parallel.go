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
