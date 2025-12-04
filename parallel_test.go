package parallel

import (
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

// RunWithLimit tests

func TestRunWithLimitEmpty(t *testing.T) {
	RunWithLimit(5)
}

func TestRunWithLimitZero(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Expected panic for limit = 0")
		}
	}()

	RunWithLimit(0, func() {})
}

func TestRunWithLimitNegative(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Expected panic for limit < 0")
		}
	}()

	RunWithLimit(-1, func() {})
}

func TestRunWithLimit(t *testing.T) {
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

	RunWithLimit(3, fns...)

	if counter != 10 {
		t.Errorf("Expected counter to be 10, got %d", counter)
	}

	if maxConcurrent > 3 {
		t.Errorf("Expected max concurrent to be <= 3, got %d", maxConcurrent)
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

// RunForEachWithLimit tests

func TestRunForEachWithLimitEmpty(t *testing.T) {
	var counter int32
	RunForEachWithLimit(5, []int{}, func(item int) {
		atomic.AddInt32(&counter, 1)
	})

	if counter != 0 {
		t.Errorf("Expected counter to be 0, got %d", counter)
	}
}

func TestRunForEachWithLimitZero(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Expected panic for limit = 0")
		}
	}()

	RunForEachWithLimit(0, []int{1, 2, 3}, func(item int) {})
}

func TestRunForEachWithLimitNegative(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Expected panic for limit < 0")
		}
	}()

	RunForEachWithLimit(-1, []int{1, 2, 3}, func(item int) {})
}

func TestRunForEachWithLimit(t *testing.T) {
	var counter int32
	var concurrent int32
	var maxConcurrent int32

	items := make([]int, 10)
	for i := range items {
		items[i] = i
	}

	RunForEachWithLimit(3, items, func(item int) {
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
