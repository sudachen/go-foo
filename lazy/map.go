package lazy

import (
	"reflect"
)

/*
Map transforms stream into new one

	rs = lazy.New([]int{0,1,2,3}).Map(func(r int)string{ return fmt.Sprint(r)}).Collect().([]string)
	rs -> {"0","1","2","3"}
*/
func (z *Stream) Map(f interface{}) *Stream {
	vf := reflect.ValueOf(f)
	vt := vf.Type()
	if !isTransformFunc(vt) {
		panic("only func(any)any is allowed as an argument")
	}
	fn := func(index int, v reflect.Value, _ interface{}) reflect.Value {
		return vf.Call([]reflect.Value{v})[0]
	}
	return &Stream{Func: fn, Src: z, Tp: vt.Out(0)}
}
