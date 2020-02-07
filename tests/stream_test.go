package tests

import (
	"github.com/sudachen/go-fp/lazy"
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func Test_SimpleStream(t *testing.T) {
	getf := func(index int, ctx interface{}) reflect.Value {
		return reflect.ValueOf(false)
	}
	z := &lazy.Stream{Ctx: &lazy.Counter{}, Get: getf, Tp: reflect.TypeOf(struct{}{})}
	x := z.Collect().([]struct{})
	assert.Assert(t, len(x) == 0)
}
