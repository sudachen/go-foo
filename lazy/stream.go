/*
Package lazy implements lazy transformation flow
*/
package lazy

import (
	"reflect"
)

/*
Stream implements lazy stream for transformations
*/
type Stream struct {
	Ctx interface{} // a context of gathering/transformations, can be nil
	// if not nil must be readonly or synchronized in Func and/or in Get

	Tp reflect.Type // return type for the Func function

	// transformation function
	// can be nil if there is no transformation and Tp the same as result of Get
	// can be called concurrently
	// returns reflect.ValueOf(true) if result must not be used (filtered out for example)
	// returns reflect.ValueOf(false) if there are no more values
	Func func(index int, a reflect.Value, ctx interface{}) reflect.Value

	// the function getting values from any source, can be nil if Src defined
	// can be called concurrently
	// returns reflect.ValueOf(true) if result must not be used (filtered out for example)
	// returns reflect.ValueOf(false) if there are no more values
	Getf func(index int, ctx interface{}) reflect.Value

	// the source stream
	// can be nil if Get is defined
	Src *Stream

	// normally if Get/Src.Next returns boolean transformation does not applied
	// CatchAll = true means apply transformation to boolean value but ignore transformation result
	CatchAll bool
}

/*
New creates new lazy transformation source from the channel of structs
*/
func New(c interface{}) *Stream {
	v := reflect.ValueOf(c)
	vt := v.Type()
	if v.Kind() == reflect.Chan {
		scase := []reflect.SelectCase{{Dir: reflect.SelectRecv, Chan: v}}
		getf := func(index int, ctx interface{}) reflect.Value {
			ctx.(*WaitCounter).Wait(index)
			_, r, ok := reflect.Select(scase)
			ctx.(*WaitCounter).Inc()
			if !ok {
				return reflect.ValueOf(false)
			}
			return r
		}
		return &Stream{Tp: vt.Elem(), Getf: getf, Ctx: &WaitCounter{Value: 0}}
	} else if v.Kind() == reflect.Slice {
		l := v.Len()
		tp := vt.Elem()
		getf := func(index int, ctx interface{}) reflect.Value {
			if index < l {
				return v.Index(index)
			}
			return reflect.ValueOf(false)
		}
		return &Stream{Tp: tp, Getf: getf}
	} else {
		panic("only `chan any` and []any are allowed as an argument")
	}
}

/*
Get gets value with specified index from the stream

	if index out of the stream range then returns reflect.ValueOf(false)
	can awaits for the index in the source stream
	returns result of stream transformation
*/
func (z *Stream) Get(index int) reflect.Value {
	var r reflect.Value
	if z.Getf != nil {
		r = z.Getf(index, z.Ctx)
	} else {
		if z.Src == nil {
			panic("both Get and Src are nil, it's impossible to get next value")
		}
		r = z.Src.Get(index)
	}
	if r.Kind() == reflect.Bool {
		if z.CatchAll && z.Func != nil {
			z.Func(index, r, z.Ctx)
		}
		return r
	}
	if z.Func != nil {
		r = z.Func(index, r, z.Ctx)
	}
	return r
}

// isFilterFunc checks is value a predicate function or not
func isFilterFunc(vt reflect.Type) bool {
	return vt.Kind() == reflect.Func &&
		vt.NumIn() == 1 && vt.NumOut() == 1 &&
		vt.Out(0).Kind() == reflect.Bool
}

// isTransformFunc checks is value a transform function or not
func isTransformFunc(vt reflect.Type) bool {
	return vt.Kind() == reflect.Func &&
		vt.NumIn() == 1 && vt.NumOut() == 1
}
