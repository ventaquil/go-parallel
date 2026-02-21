package parallel

import "context"

// Run executes the given functions in parallel and waits for all to complete.
// It returns when all functions have finished execution.
// This implementation uses WaitGroup.Go available in Go 1.25+.
func Run(fns ...func() error) error {
	return RunForEach(fns, func(fn func() error) error {
		return fn()
	})
}

// RunContext executes the given functions in parallel and waits for all to complete.
// It returns the first error encountered, if any, and cancels the context for remaining functions.
// Only the first error is returned; subsequent errors are discarded.
// This implementation uses WaitGroup.Go available in Go 1.25+.
func RunContext(ctx context.Context, fns ...func(ctx context.Context) error) error {
	return RunForEachContext(ctx, fns, func(ctx context.Context, fn func(ctx context.Context) error) error {
		return fn(ctx)
	})
}

// RunLimit executes the given functions in parallel with a concurrency limit.
// It ensures that at most 'limit' functions execute concurrently.
// It returns an error if limit is less than or equal to 0.
// This implementation uses WaitGroup.Go available in Go 1.25+.
func RunLimit(limit int, fns ...func() error) error {
	return RunForEachLimitContext(context.Background(), limit, fns, func(_ context.Context, fn func() error) error {
		return fn()
	})
}

// RunLimitContext executes the given functions in parallel with a concurrency limit.
// It ensures that at most 'limit' functions execute concurrently.
// It returns an error if limit is less than or equal to 0.
// It returns the first error encountered, if any, and cancels the context for remaining functions.
// Only the first error is returned; subsequent errors are discarded.
// This implementation uses WaitGroup.Go available in Go 1.25+.
func RunLimitContext(ctx context.Context, limit int, fns ...func(ctx context.Context) error) error {
	return RunForEachLimitContext(ctx, limit, fns, func(ctx context.Context, fn func(ctx context.Context) error) error {
		return fn(ctx)
	})
}

// RunForEach executes the given function for each item in the slice in parallel.
// It waits for all executions to complete before returning.
// This implementation uses WaitGroup.Go available in Go 1.25+.
func RunForEach[T any](items []T, fn func(item T) error) error {
	return RunForEachContext(context.Background(), items, func(_ context.Context, item T) error {
		return fn(item)
	})
}

// RunForEachLimit executes the given function for each item in the slice in parallel with a concurrency limit.
// It ensures that at most 'limit' functions execute concurrently.
// It returns an error if limit is less than or equal to 0.
// This implementation uses WaitGroup.Go available in Go 1.25+.
func RunForEachLimit[T any](limit int, items []T, fn func(item T) error) error {
	return RunForEachLimitContext(context.Background(), limit, items, func(_ context.Context, item T) error {
		return fn(item)
	})
}
