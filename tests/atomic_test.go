package tests

import (
	"github.com/sudachen/go-fp/lazy"
	"gotest.tools/assert"
	"testing"
)

func Test_Atomic1(t *testing.T) {
	f := lazy.AtomicFlag{1}
	assert.Assert(t, f.State() == true)
	f.Clear()
	assert.Assert(t, f.State() == false)
	f.Set()
	assert.Assert(t, f.State() == true)
	f.Clear()
	assert.Assert(t, f.State() == false)

	f = lazy.AtomicFlag{0}
	assert.Assert(t, f.State() == false)
	f.Clear()
	assert.Assert(t, f.State() == false)
	f.Set()
	assert.Assert(t, f.State() == true)
}
