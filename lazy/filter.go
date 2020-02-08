package lazy

import (
	"reflect"
)

/*
Filter transforms stream filtering out all records not passed predicate

	struct R{Index int}
	rs = lazy.New([]R{0,1,2,3}).Filter(func(r R)bool{ return r.Index%2 == 0}).Collect().([]R)
	rs -> {0,2}
*/
func (z *Stream) Filter(f interface{}) *Stream {
	vf := reflect.ValueOf(f)
	vt := vf.Type()
	if !isFilterFunc(vt) {
		panic("only func(any)bool is allowed as an argument")
	}
	fn := func(index int, v reflect.Value) reflect.Value {
		r := vf.Call([]reflect.Value{v})[0]
		if r.Bool() {
			return v
		}
		return reflect.ValueOf(true)
	}
	return &Stream{Func: fn, Src: z, Tp: z.Tp}
}
