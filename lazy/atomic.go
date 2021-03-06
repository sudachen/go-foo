package lazy

import (
	"sync/atomic"
)

/*
AtomicCounter - hm, yes it's atomic counter
*/
type AtomicCounter struct {
	Value uint64
}

/*
PostInc increments counter and returns OLD value
*/
func (c *AtomicCounter) PostInc() uint64 {
	for {
		v := atomic.LoadUint64(&c.Value)
		if atomic.CompareAndSwapUint64(&c.Value, v, v+1) {
			return v
		}
	}
}

/*
Dec decrements counter and returns NEW value
*/
func (c *AtomicCounter) Dec() uint64 {
	for {
		v := atomic.LoadUint64(&c.Value)
		if v == 0 {
			panic("counter underflow")
		}
		if atomic.CompareAndSwapUint64(&c.Value, v, v-1) {
			return v - 1
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
Clear switches Integer to 0 atomically
*/
func (c *AtomicFlag) Clear() (r bool) {
	for atomic.LoadInt32(&c.Value) == 1 {
		r = atomic.CompareAndSwapInt32(&c.Value, 1, 0)
		if r {
			return
		}
	}
	return
}

/*
Set switches Integer to 1 atomically
*/
func (c *AtomicFlag) Set() (r bool) {
	for atomic.LoadInt32(&c.Value) == 0 {
		r = atomic.CompareAndSwapInt32(&c.Value, 0, 1)
		if r {
			break
		}
	}
	return
}

/*
State returns current state
*/
func (c *AtomicFlag) State() bool {
	v := atomic.LoadInt32(&c.Value)
	return bool(v != 0)
}
