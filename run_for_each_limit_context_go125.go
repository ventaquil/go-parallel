//go:build go1.25

package parallel

import (
	"context"
	"sync"
)

// RunForEachLimitContext executes the given function for each item in the slice in parallel with a concurrency limit.
// It ensures that at most 'limit' functions execute concurrently.
// It returns an error if limit is less than or equal to 0.
// It returns the first error encountered, if any, and cancels the context for remaining functions.
// Only the first error is returned; subsequent errors are discarded.
// This implementation uses WaitGroup.Go available in Go 1.25+.
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
	for _, item := range items {
		sem <- struct{}{}

		wg.Go(func() {
			defer func() {
				<-sem
			}()
			if err := fn(ctx, item); err != nil {
				errs <- err
				cancelCtx()
			}
		})
	}
	wg.Wait()

	select {
	case err := <-errs:
		return err
	default:
		return nil
	}
}
