//
// Package tables implements lazy data transformation flow
//
package lazy

import (
	"reflect"
)

var refFalse = reflect.ValueOf(false)
var refTrue = reflect.ValueOf(true)

/*
Stream implements lazy stream for transformations
*/
type Stream struct {
	Ctx interface{} // a context of transformations, can be nil
	// if not nil must be readonly or synchronized in Func

	Tp reflect.Type // return type for the Func function

	// transformation function
	// can be nil if there is no transformation and Tp the same as result of Get
	// can be called concurrent
	// returns reflect.ValueOf(true) if result must not be used (filtered out for example)
	// returns reflect.ValueOf(false) if there are no more values
	Func func(index int, a reflect.Value, ctx interface{}) reflect.Value

	// the function getting values from stream source, can be nil if Src defined
	// can be call concurrent
	// returns reflect.ValueOf(true) if result must not be used (filtered out for example)
	// returns reflect.ValueOf(false) if there are no more values
	Get func(index int, ctx interface{}) reflect.Value

	// the source stream
	// can be nil if Get is defined
	Src *Stream

	// normally if Get/Src.Next returns boolean transformation does not applied
	// CatchAll = true means apply transformation to boolean value but ignore transformation result
	CatchAll bool
}

type Counter struct {
	Value int
}

func (c *Counter) Await(index int) {

}

func (c *Counter) Inc() {

}

func (c *Counter) GetInc() int {
	return 0
}

func (c *Counter) IncAwait(index int) {

}

/*
New creates new lazy transformation source from the channel of structs
*/
func New(c interface{}) *Stream {
	v := reflect.ValueOf(c)
	if v.Kind() == reflect.Chan && v.Elem().Kind() == reflect.Struct {
		scase := []reflect.SelectCase{{Dir: reflect.SelectRecv, Chan: v}}
		getf := func(index int, ctx interface{}) reflect.Value {
			ctx.(*Counter).Await(index)
			_, r, ok := reflect.Select(scase)
			ctx.(*Counter).Inc()
			if !ok {
				return reflect.ValueOf(false)
			}
			return r
		}
		return &Stream{Tp: v.Elem().Type(), Get: getf, Ctx: &Counter{0}}
	} else {
		panic("only chan struct{...} is allowed as an argument")
	}
}

/*
Next takes next value with in sequence with specified index

	if current index more then specified returns reflect.ValueOf(true)
	if current index less then specified awaits for index
	otherwise returns result of stream transformation
*/
func (z *Stream) Next(index int) reflect.Value {
	var r reflect.Value
	if z.Get != nil {
		r = z.Get(index, z.Ctx)
	} else {
		if z.Src == nil {
			panic("both Get and Src are nil, it's impossible to get next value")
		}
		r = z.Src.Next(index)
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
