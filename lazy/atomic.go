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
func (c *WaitCounter) Wait(index int64) {
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
AtomicCounter - hm, yes it's atomic counter
*/
type AtomicCounter struct {
	Value int64
}

/*
Inc increments counter
*/
func (c *AtomicCounter) Inc() int64 {
	for {
		v := atomic.LoadInt64(&c.Value)
		if atomic.CompareAndSwapInt64(&c.Value, v, v+1) {
			return v
		}
	}
}

/*
AtomicFlag - hm, yes it's atomic flag
*/
type AtomicFlag struct {
	Value int32
}

/*
Off Switches Value to 0
*/
func (c *AtomicFlag) Off() {
	for {
		v := atomic.LoadInt32(&c.Value)
		if v == 0 || atomic.CompareAndSwapInt32(&c.Value, v, 0) {
			break
		}
	}
}

/*
On Switches Value to 1
*/
func (c *AtomicFlag) On() {
	for {
		v := atomic.LoadInt32(&c.Value)
		if v != 0 || atomic.CompareAndSwapInt32(&c.Value, v, 1) {
			break
		}
	}
}

/*
State returns current state
*/
func (c *AtomicFlag) State() bool {
	v := atomic.LoadInt32(&c.Value)
	return bool(v != 0)
}
