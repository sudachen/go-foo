package internal

import (
	"os"
	"sync/atomic"
	"time"
)

type NaiveRandom struct {
	value int32
}

func (nr *NaiveRandom) Reseed() {
	atomic.StoreInt32(&nr.value,int32(time.Now().UnixNano() + int64(os.Getpid())))
}

func (nr *NaiveRandom) Next() int32 {
	var r int32
	for {
		r = atomic.LoadInt32(&nr.value)
		rx := r
		if r == 0 {
			r = int32(time.Now().UnixNano() + int64(os.Getpid()))
		}
		r = r*1664525 + 1013904223 // constants from Numerical Recipes
		if atomic.CompareAndSwapInt32(&nr.value, rx, r) {
			break
		}
	}
	return r
}

