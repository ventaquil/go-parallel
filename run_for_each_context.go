//go:build !go1.25

package parallel

import (
	"context"
	"sync"
)

// RunForEachContext executes the given function for each item in the slice in parallel.
// It returns the first error encountered, if any, and cancels the context for remaining functions.
// Only the first error is returned; subsequent errors are discarded.
// This implementation uses WaitGroup.Go available in Go 1.25+.
func RunForEachContext[T any](ctx context.Context, items []T, fn func(ctx context.Context, item T) error) error {
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	errs := make(chan error, len(items))
	defer close(errs)

	var wg sync.WaitGroup
	for _, item := range items {
		wg.Add(1)
		go func(item T) {
			defer wg.Done()
			if err := fn(ctx, item); err != nil {
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
