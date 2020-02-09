package tests

import (
	"github.com/sudachen/go-fp/lazy"
	"gotest.tools/assert"
	"gotest.tools/assert/cmp"
	"testing"
)

func Test_Counter1(t *testing.T) {
	assert.Assert(t, cmp.Panics(func() {
		wc := lazy.WaitCounter{Value: 10}
		wc.Wait(1)
	}))
	wc := lazy.WaitCounter{Value: 0}
	wc.Inc()
	wc.Wait(1)
	wc.Inc()
	assert.Assert(t, wc.Value == 2)
}

func Test_Counter2(t *testing.T) {
	f := lazy.AtomicFlag{1}
	assert.Assert(t, f.State() == true)
	f.Off()
	assert.Assert(t, f.State() == false)
	f.On()
	assert.Assert(t, f.State() == true)
	f.Off()
	assert.Assert(t, f.State() == false)

	f = lazy.AtomicFlag{0}
	assert.Assert(t, f.State() == false)
	f.Off()
	assert.Assert(t, f.State() == false)
	f.On()
	assert.Assert(t, f.State() == true)
}
