package lazy

import "sync/atomic"

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
Clear switches Value to 0 atomically
*/
func (c *AtomicFlag) Clear() {
	for {
		v := atomic.LoadInt32(&c.Value)
		if v == 0 || atomic.CompareAndSwapInt32(&c.Value, v, 0) {
			break
		}
	}
}

/*
Set switches Value to 1 atomically
*/
func (c *AtomicFlag) Set() {
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
