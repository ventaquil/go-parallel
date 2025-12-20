package parallel

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Run tests

func TestRunEmpty(t *testing.T) {
	Run()
}

func TestRun(t *testing.T) {
	var counter int32

	Run(
		func() { atomic.AddInt32(&counter, 1) },
		func() { atomic.AddInt32(&counter, 1) },
		func() { atomic.AddInt32(&counter, 1) },
	)

	if counter != 3 {
		t.Errorf("Expected counter to be 3, got %d", counter)
	}
}

func TestRunConcurrency(t *testing.T) {
	var mu sync.Mutex
	var order []int

	Run(
		func() {
			time.Sleep(20 * time.Millisecond)
			mu.Lock()
			order = append(order, 1)
			mu.Unlock()
		},
		func() {
			time.Sleep(10 * time.Millisecond)
			mu.Lock()
			order = append(order, 2)
			mu.Unlock()
		},
		func() {
			mu.Lock()
			order = append(order, 3)
			mu.Unlock()
		},
	)

	if len(order) != 3 {
		t.Errorf("Expected 3 executions, got %d", len(order))
	}

	if order[0] != 3 {
		t.Errorf("Expected first completion to be 3, got %d", order[0])
	}
}

// RunContext tests

func TestRunContextSuccess(t *testing.T) {
	ctx := context.Background()
	var counter int32

	err := RunContext(ctx,
		func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		},
		func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		},
		func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		},
	)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if counter != 3 {
		t.Errorf("Expected counter to be 3, got %d", counter)
	}
}

func TestRunContextError(t *testing.T) {
	ctx := context.Background()
	testErr := errors.New("test error")

	err := RunContext(ctx,
		func(ctx context.Context) error {
			time.Sleep(10 * time.Millisecond)
			return nil
		},
		func(ctx context.Context) error {
			return testErr
		},
		func(ctx context.Context) error {
			time.Sleep(20 * time.Millisecond)
			return nil
		},
	)

	if err != testErr {
		t.Errorf("Expected error %v, got %v", testErr, err)
	}
}

func TestRunContextCancellation(t *testing.T) {
	ctx := context.Background()
	testErr := errors.New("test error")
	var executed int32

	err := RunContext(ctx,
		func(ctx context.Context) error {
			atomic.AddInt32(&executed, 1)
			return testErr
		},
		func(ctx context.Context) error {
			time.Sleep(50 * time.Millisecond)
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				atomic.AddInt32(&executed, 1)
				return nil
			}
		},
	)

	if err != testErr {
		t.Errorf("Expected error %v, got %v", testErr, err)
	}
}

// RunLimit tests

func TestRunLimitEmpty(t *testing.T) {
	err := RunLimit(5)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestRunLimitZero(t *testing.T) {
	err := RunLimit(0, func() {})
	if err == nil {
		t.Error("Expected error for limit = 0")
	}
}

func TestRunLimitNegative(t *testing.T) {
	err := RunLimit(-1, func() {})
	if err == nil {
		t.Error("Expected error for limit < 0")
	}
}

func TestRunLimit(t *testing.T) {
	var counter int32
	var concurrent int32
	var maxConcurrent int32

	fns := make([]func(), 10)
	for i := 0; i < 10; i++ {
		fns[i] = func() {
			current := atomic.AddInt32(&concurrent, 1)
			for {
				max := atomic.LoadInt32(&maxConcurrent)
				if current <= max || atomic.CompareAndSwapInt32(&maxConcurrent, max, current) {
					break
				}
			}
			time.Sleep(10 * time.Millisecond)
			atomic.AddInt32(&counter, 1)
			atomic.AddInt32(&concurrent, -1)
		}
	}

	RunLimit(3, fns...)

	if counter != 10 {
		t.Errorf("Expected counter to be 10, got %d", counter)
	}

	if maxConcurrent > 3 {
		t.Errorf("Expected max concurrent to be <= 3, got %d", maxConcurrent)
	}
}

// RunLimitContext tests

func TestRunLimitContextSuccess(t *testing.T) {
	ctx := context.Background()
	var counter int32

	err := RunLimitContext(ctx, 2,
		func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		},
		func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		},
		func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		},
	)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if counter != 3 {
		t.Errorf("Expected counter to be 3, got %d", counter)
	}
}

func TestRunLimitContextZero(t *testing.T) {
	ctx := context.Background()

	err := RunLimitContext(ctx, 0, func(ctx context.Context) error {
		return nil
	})

	if err == nil {
		t.Error("Expected error for limit = 0")
	}
}

func TestRunLimitContextError(t *testing.T) {
	ctx := context.Background()
	testErr := errors.New("test error")

	err := RunLimitContext(ctx, 2,
		func(ctx context.Context) error {
			time.Sleep(10 * time.Millisecond)
			return nil
		},
		func(ctx context.Context) error {
			return testErr
		},
	)

	if err != testErr {
		t.Errorf("Expected error %v, got %v", testErr, err)
	}
}

// RunForEach tests

func TestRunForEachEmpty(t *testing.T) {
	var counter int32
	RunForEach([]int{}, func(item int) {
		atomic.AddInt32(&counter, 1)
	})

	if counter != 0 {
		t.Errorf("Expected counter to be 0, got %d", counter)
	}
}

func TestRunForEach(t *testing.T) {
	var counter int32

	RunForEach([]int{1, 2, 3, 4, 5}, func(item int) {
		atomic.AddInt32(&counter, int32(item))
	})

	if counter != 15 {
		t.Errorf("Expected counter to be 15, got %d", counter)
	}
}

func TestRunForEachGenerics(t *testing.T) {
	var mu sync.Mutex
	var results []string

	RunForEach([]string{"a", "b", "c"}, func(item string) {
		mu.Lock()
		results = append(results, item)
		mu.Unlock()
	})

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}
}

func TestRunForEachConcurrency(t *testing.T) {
	var mu sync.Mutex
	var order []int

	RunForEach([]int{10, 20, 30}, func(item int) {
		time.Sleep(time.Duration(40-item) * time.Millisecond)
		mu.Lock()
		order = append(order, item)
		mu.Unlock()
	})

	if len(order) != 3 {
		t.Errorf("Expected 3 results, got %d", len(order))
	}

	if order[0] != 30 {
		t.Errorf("Expected first result to be 30, got %d", order[0])
	}
}

// RunForEachContext tests

func TestRunForEachContextSuccess(t *testing.T) {
	ctx := context.Background()
	var sum int32

	err := RunForEachContext(ctx, []int{1, 2, 3, 4, 5}, func(ctx context.Context, item int) error {
		atomic.AddInt32(&sum, int32(item))
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if sum != 15 {
		t.Errorf("Expected sum to be 15, got %d", sum)
	}
}

func TestRunForEachContextError(t *testing.T) {
	ctx := context.Background()
	testErr := errors.New("test error")

	err := RunForEachContext(ctx, []int{1, 2, 3}, func(ctx context.Context, item int) error {
		if item == 2 {
			return testErr
		}
		return nil
	})

	if err != testErr {
		t.Errorf("Expected error %v, got %v", testErr, err)
	}
}

func TestRunForEachContextEmpty(t *testing.T) {
	ctx := context.Background()

	err := RunForEachContext(ctx, []int{}, func(ctx context.Context, item int) error {
		return errors.New("should not be called")
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// RunForEachLimit tests

func TestRunForEachLimitEmpty(t *testing.T) {
	var counter int32
	err := RunForEachLimit(5, []int{}, func(item int) {
		atomic.AddInt32(&counter, 1)
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if counter != 0 {
		t.Errorf("Expected counter to be 0, got %d", counter)
	}
}

func TestRunForEachWithLimitZero(t *testing.T) {
	err := RunForEachLimit(0, []int{1, 2, 3}, func(item int) {})
	if err == nil {
		t.Error("Expected error for limit = 0")
	}
}

func TestRunForEachLimitNegative(t *testing.T) {
	err := RunForEachLimit(-1, []int{1, 2, 3}, func(item int) {})
	if err == nil {
		t.Error("Expected error for limit < 0")
	}
}

func TestRunForEachLimit(t *testing.T) {
	var counter int32
	var concurrent int32
	var maxConcurrent int32

	items := make([]int, 10)
	for i := range items {
		items[i] = i
	}

	RunForEachLimit(3, items, func(item int) {
		current := atomic.AddInt32(&concurrent, 1)
		for {
			max := atomic.LoadInt32(&maxConcurrent)
			if current <= max || atomic.CompareAndSwapInt32(&maxConcurrent, max, current) {
				break
			}
		}
		time.Sleep(10 * time.Millisecond)
		atomic.AddInt32(&counter, 1)
		atomic.AddInt32(&concurrent, -1)
	})

	if counter != 10 {
		t.Errorf("Expected counter to be 10, got %d", counter)
	}

	if maxConcurrent > 3 {
		t.Errorf("Expected max concurrent to be <= 3, got %d", maxConcurrent)
	}
}

// RunForEachLimitContext tests

func TestRunForEachLimitContextSuccess(t *testing.T) {
	ctx := context.Background()
	var sum int32

	err := RunForEachLimitContext(ctx, 2, []int{1, 2, 3, 4, 5}, func(ctx context.Context, item int) error {
		atomic.AddInt32(&sum, int32(item))
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if sum != 15 {
		t.Errorf("Expected sum to be 15, got %d", sum)
	}
}

func TestRunForEachLimitContextZero(t *testing.T) {
	ctx := context.Background()

	err := RunForEachLimitContext(ctx, 0, []int{1, 2, 3}, func(ctx context.Context, item int) error {
		return nil
	})

	if err == nil {
		t.Error("Expected error for limit = 0")
	}
}

func TestRunForEachLimitContextError(t *testing.T) {
	ctx := context.Background()
	testErr := errors.New("test error")

	err := RunForEachLimitContext(ctx, 2, []int{1, 2, 3, 4}, func(ctx context.Context, item int) error {
		if item == 3 {
			return testErr
		}
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	if err != testErr {
		t.Errorf("Expected error %v, got %v", testErr, err)
	}
}

func TestRunForEachLimitContextConcurrency(t *testing.T) {
	ctx := context.Background()
	var concurrent int32
	var maxConcurrent int32

	items := make([]int, 10)
	for i := range items {
		items[i] = i
	}

	err := RunForEachLimitContext(ctx, 3, items, func(ctx context.Context, item int) error {
		current := atomic.AddInt32(&concurrent, 1)
		for {
			max := atomic.LoadInt32(&maxConcurrent)
			if current <= max || atomic.CompareAndSwapInt32(&maxConcurrent, max, current) {
				break
			}
		}
		time.Sleep(10 * time.Millisecond)
		atomic.AddInt32(&concurrent, -1)
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if maxConcurrent > 3 {
		t.Errorf("Expected max concurrent to be <= 3, got %d", maxConcurrent)
	}
}
