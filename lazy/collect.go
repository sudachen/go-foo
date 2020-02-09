package lazy

import (
	"reflect"
	"sync"
)

/*
Collect executes all lazy transformations and collects result to array
*/
func (z *Stream) Collect() interface{} {
	r := reflect.MakeSlice(reflect.SliceOf(z.Tp), 0, 0)
	index := int64(0)
	for {
		v := z.Next(index)
		index++
		if v.Kind() == reflect.Bool {
			if !v.Bool() {
				break
			}
		} else {
			r = reflect.Append(r, v)
		}
	}
	return r.Interface()
}

/*
ConqCollect executes all lazy transformations and collects result to array.
	concurrency - count of goroutines executing transformations
*/
func (z *Stream) ConqCollect(concurrency int) interface{} {
	r := reflect.MakeSlice(reflect.SliceOf(z.Tp), 0, 0)
	index := &AtomicCounter{0}
	wc := &WaitCounter{Value: 0}
	gw := sync.WaitGroup{}
	gw.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer gw.Done()
			for {
				n := index.Inc()
				v := z.Next(n)
				wc.Wait(n)
				if v.Kind() != reflect.Bool {
					r = reflect.Append(r, v)
				}
				wc.Inc()
				if v.Kind() == reflect.Bool && !v.Bool() {
					break
				}
			}
		}()
	}
	gw.Wait()
	return r.Interface()
}
