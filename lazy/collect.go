package lazy

import "reflect"

func (z *Stream) Collect() interface{} {
	r := reflect.MakeSlice(reflect.SliceOf(z.Tp), 0, 0)
	return r.Interface()
}

func (z *Stream) ConqCollect(concurrency int) interface{} {
	r := reflect.MakeSlice(reflect.SliceOf(z.Tp), 0, 0)
	return r.Interface()
}
