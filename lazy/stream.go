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
	Tp reflect.Type // return type for the Func function

	// transformation function
	// can be nil if there is no transformation and Tp the same as result of Getf (or Src.Tp)
	// can be called concurrently
	// returns reflect.ValueOf(true) if result must not be used (filtered out for example)
	// returns reflect.ValueOf(false) if there are no more values
	Func func(index int64, a reflect.Value) reflect.Value

	// the function getting values from any source, can be nil if Src defined
	// can be called concurrently
	// returns reflect.ValueOf(true) if result must not be used (filtered out for example)
	// returns reflect.ValueOf(false) if there are no more values
	Getf func(index int64) reflect.Value

	// the source stream
	// can be nil if Get is defined
	Src *Stream

	// normally if Get/Src.Next returns boolean transformation does not applied
	// CatchAll = true means apply transformation to boolean value but ignore transformation result
	CatchAll bool

	// to stop producing new values
	// can be nil if transforms values only
	Stopf func()
}

/*
New creates new lazy transformation source from the channel of structs
*/
func New(c interface{}, stop ...chan struct{}) *Stream {
	v := reflect.ValueOf(c)
	vt := v.Type()
	if v.Kind() == reflect.Chan {
		scase := []reflect.SelectCase{{Dir: reflect.SelectRecv, Chan: v}}
		wc := &WaitCounter{Value: 0}
		getf := func(index int64) reflect.Value {
			if wc.Wait(index) {
				_, r, ok := reflect.Select(scase)
				if wc.Inc() && ok {
					return r
				}
			}
			return reflect.ValueOf(false)
		}
		stopf := func() {
			wc.Stop()
			for _, s := range stop {
				close(s)
			}
		}
		return &Stream{Tp: vt.Elem(), Getf: getf, Stopf: stopf}
	} else if v.Kind() == reflect.Slice {
		l := int64(v.Len())
		tp := vt.Elem()
		flag := &AtomicFlag{Value: 1}
		getf := func(index int64) reflect.Value {
			if index < l && flag.State() {
				return v.Index(int(index))
			}
			return reflect.ValueOf(false)
		}
		return &Stream{Tp: tp, Getf: getf, Stopf: func() { flag.Clear() }}
	} else {
		panic("only `chan any` and []any are allowed as an argument")
	}
}

/*
Next gets (next?) value with specified index from the stream

	if index out of the stream range then returns reflect.ValueOf(false)
	can awaits for the index in the source stream
	returns result of stream transformation
*/
func (z *Stream) Next(index int64) reflect.Value {
	var r reflect.Value
	if z.Getf != nil {
		r = z.Getf(index)
	} else {
		if z.Src == nil {
			panic("both Get and Src are nil, it's impossible to get next value")
		}
		r = z.Src.Next(index)
	}
	if r.Kind() == reflect.Bool {
		if z.CatchAll && z.Func != nil {
			z.Func(index, r)
		}
		return r
	}
	if z.Func != nil {
		r = z.Func(index, r)
	}
	return r
}

/*
Close stops stream to produce new values
*/
func (z *Stream) Close() {
	for x := z; x != nil; x = x.Src {
		if x.Stopf != nil {
			x.Stopf()
		}
	}
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
