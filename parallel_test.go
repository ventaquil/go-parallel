package parallel

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	var counter int32

	fns := make([]func(), 5)
	for i := 0; i < 5; i++ {
		fns[i] = func() {
			atomic.AddInt32(&counter, 1)
		}
	}

	Run(fns...)

	if counter != 5 {
		t.Errorf("Expected counter to be 5, got %d", counter)
	}
}

func TestRunEmpty(t *testing.T) {
	// Should not panic with empty input
	Run()
}

func TestRunWithLimit(t *testing.T) {
	var counter int32
	var concurrent int32
	var maxConcurrent int32

	fns := make([]func(), 10)
	for i := 0; i < 10; i++ {
		fns[i] = func() {
			current := atomic.AddInt32(&concurrent, 1)

			// Track max concurrent executions
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

	RunWithLimit(3, fns...)

	if counter != 10 {
		t.Errorf("Expected counter to be 10, got %d", counter)
	}

	if maxConcurrent > 3 {
		t.Errorf("Expected max concurrent to be <= 3, got %d", maxConcurrent)
	}
}

func TestRunWithLimitZero(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for limit <= 0, but didn't panic")
		}
	}()

	fns := make([]func(), 3)
	for i := 0; i < 3; i++ {
		fns[i] = func() {}
	}

	// Should panic when limit <= 0
	RunWithLimit(0, fns...)
}

func TestRunWithLimitNegative(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for negative limit, but didn't panic")
		}
	}()

	fns := make([]func(), 3)
	for i := 0; i < 3; i++ {
		fns[i] = func() {}
	}

	// Should panic when limit is negative
	RunWithLimit(-1, fns...)
}

func TestRunWithLimitEmpty(t *testing.T) {
	// Should not panic with empty input
	RunWithLimit(5)
}

func TestRunConcurrency(t *testing.T) {
	var mu sync.Mutex
	var order []int

	fns := []func(){
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
	}

	Run(fns...)

	// All functions should have executed
	if len(order) != 3 {
		t.Errorf("Expected 3 executions, got %d", len(order))
	}

	// Function 3 should complete first (no sleep), then 2, then 1
	if order[0] != 3 {
		t.Errorf("Expected first completion to be function 3, got %d", order[0])
	}
}
