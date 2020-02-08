package lazy

import (
	"sync"
	"sync/atomic"
)

/*
WaitCounter implements barrier counter for lazy flow execution synchronization
*/
type WaitCounter struct {
	Value int
	cond  sync.Cond
	mu    sync.Mutex
}

/*
Wait waits until counter index is not equal to specified
*/
func (c *WaitCounter) Wait(index int) {
	c.mu.Lock()
	if c.cond.L == nil {
		c.cond.L = &c.mu
	}
	for c.Value != index {
		if c.Value > index {
			panic("index continuity broken")
		}
		c.cond.Wait()
	}
	c.mu.Unlock()
}

/*
Inc increments index and notifies waiting goroutines
*/
func (c *WaitCounter) Inc() {
	c.mu.Lock()
	if c.cond.L == nil {
		c.cond.L = &c.mu
	}
	c.Value++
	c.mu.Unlock()
	c.cond.Broadcast()
}

/*
AtcomicCounter - hm, yes it's atomic counter
*/
type AtomicCounter struct {
	Value int32
}

/*
Inc increments counter
*/
func (c *AtomicCounter) Inc() int {
	return int(atomic.AddInt32(&c.Value, 1) - 1)
}
