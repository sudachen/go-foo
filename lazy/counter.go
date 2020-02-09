package lazy

import (
	"sync"
	"sync/atomic"
)

/*
WaitCounter implements barrier counter for lazy flow execution synchronization
*/
type WaitCounter struct {
	Value int64
	cond  sync.Cond
	mu    sync.Mutex
}

/*
Wait waits until counter Value is not equal to specified index
*/
func (c *WaitCounter) Wait(index int64) bool {
	r := true
	c.mu.Lock()
	if c.cond.L == nil {
		c.cond.L = &c.mu
	}
	for c.Value != index {
		if c.Value > index {
			panic("index continuity broken")
		}
		if c.Value >= 0 {
			c.cond.Wait()
		} else {
			r = false
			break
		}
	}
	c.mu.Unlock()
	return r
}

/*
Inc increments index and notifies waiting goroutines
*/
func (c *WaitCounter) Inc() bool {
	r := false
	c.mu.Lock()
	if c.cond.L == nil {
		c.cond.L = &c.mu
	}
	if c.Value >= 0 {
		atomic.AddInt64(&c.Value, 1)
		r = true
	}
	c.mu.Unlock()
	c.cond.Broadcast()
	return r
}

/*
Stop sets Value to -1 and notifies waiting goroutines. It means also counter will not increment more
*/
func (c *WaitCounter) Stop() {
	c.mu.Lock()
	if c.cond.L == nil {
		c.cond.L = &c.mu
	}
	atomic.StoreInt64(&c.Value, -1)
	c.mu.Unlock()
	c.cond.Broadcast()
}

/*
Stopped returns true if counter is stopped and will not increment more
*/
func (c *WaitCounter) Stopped() bool {
	return atomic.LoadInt64(&c.Value) == -1
}
