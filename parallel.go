//go:build !go1.25

package parallel

import (
	"context"
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

// RunContext executes the given functions in parallel and waits for all to complete.
// It returns the first error encountered, if any, and cancels the context for remaining functions.
// Only the first error is returned; subsequent errors are discarded.
func RunContext(ctx context.Context, fns ...func(ctx context.Context) error) error {
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	errs := make(chan error, len(fns))
	defer close(errs)

	var wg sync.WaitGroup
	wg.Add(len(fns))
	for _, fn := range fns {
		go func(f func(ctx context.Context) error) {
			defer wg.Done()
			if err := f(ctx); err != nil {
				errs <- err
				cancelCtx()
			}
		}(fn)
	}
	wg.Wait()

	select {
	case err := <-errs:
		return err
	default:
		return nil
	}
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

// RunLimitContext executes the given functions in parallel with a concurrency limit.
// It ensures that at most 'limit' functions execute concurrently.
// It returns an error if limit is less than or equal to 0.
// It returns the first error encountered, if any, and cancels the context for remaining functions.
// Only the first error is returned; subsequent errors are discarded.
func RunLimitContext(ctx context.Context, limit int, fns ...func(ctx context.Context) error) error {
	if limit <= 0 {
		return ErrInvalidLimit
	}

	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	errs := make(chan error, len(fns))
	defer close(errs)

	var wg sync.WaitGroup
	sem := make(chan struct{}, limit)
	wg.Add(len(fns))
	for _, fn := range fns {
		sem <- struct{}{}

		go func(f func(ctx context.Context) error) {
			defer func() {
				<-sem
				wg.Done()
			}()
			if err := f(ctx); err != nil {
				errs <- err
				cancelCtx()
			}
		}(fn)
	}
	wg.Wait()

	select {
	case err := <-errs:
		return err
	default:
		return nil
	}
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

// RunForEachContext executes the given function for each item in the slice in parallel.
// It returns the first error encountered, if any, and cancels the context for remaining functions.
// Only the first error is returned; subsequent errors are discarded.
func RunForEachContext[T any](ctx context.Context, items []T, fn func(ctx context.Context, item T) error) error {
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	errs := make(chan error, len(items))
	defer close(errs)

	var wg sync.WaitGroup
	wg.Add(len(items))
	for _, item := range items {
		go func(it T) {
			defer wg.Done()
			if err := fn(ctx, it); err != nil {
				errs <- err
				cancelCtx()
			}
		}(item)
	}
	wg.Wait()

	select {
	case err := <-errs:
		return err
	default:
		return nil
	}
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

// RunForEachLimitContext executes the given function for each item in the slice in parallel with a concurrency limit.
// It ensures that at most 'limit' functions execute concurrently.
// It returns an error if limit is less than or equal to 0.
// It returns the first error encountered, if any, and cancels the context for remaining functions.
// Only the first error is returned; subsequent errors are discarded.
func RunForEachLimitContext[T any](ctx context.Context, limit int, items []T, fn func(ctx context.Context, item T) error) error {
	if limit <= 0 {
		return ErrInvalidLimit
	}

	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	errs := make(chan error, len(items))
	defer close(errs)

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
			if err := fn(ctx, it); err != nil {
				errs <- err
				cancelCtx()
			}
		}(item)
	}
	wg.Wait()

	select {
	case err := <-errs:
		return err
	default:
		return nil
	}
}
