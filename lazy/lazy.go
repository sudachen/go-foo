package lazy

import (
	"github.com/sudachen/go-foo/fu"
	"math"
	"reflect"
	"runtime"
)

const STOP = math.MaxUint64

type Stream func(index uint64) (reflect.Value, error)
type Source func() Stream
type Sink func(reflect.Value) error
type Parallel int

func (zf Source) Map(f interface{}) Source {
	return func() Stream {
		z := zf()
		return func(index uint64) (v reflect.Value, err error) {
			if v, err = z(index); err != nil || v.Kind() == reflect.Bool {
				return v, err
			}
			fv := reflect.ValueOf(f)
			return fv.Call([]reflect.Value{v})[0], nil
		}
	}
}

func (zf Source) Filter(f interface{}) Source {
	return func() Stream {
		z := zf()
		return func(index uint64) (v reflect.Value, err error) {
			if v, err = z(index); err != nil || v.Kind() == reflect.Bool {
				return v, err
			}
			fv := reflect.ValueOf(f)
			if fv.Call([]reflect.Value{v})[0].Bool() {
				return
			}
			return reflect.ValueOf(true), nil
		}
	}
}

func (zf Source) Parallel(concurrency ...int) Source {
	return func() Stream {
		z := zf()
		ccrn := fu.Fnzi(fu.Fnzi(concurrency...), runtime.NumCPU())
		type C struct {
			reflect.Value
			error
		}
		index := AtomicCounter{0}
		wc := WaitCounter{Value: 0}
		c := make(chan C)
		stop := make(chan struct{})
		alive := AtomicCounter{uint64(ccrn)}
		for i := 0; i < ccrn; i++ {
			go func() {
			loop:
				for !wc.Stopped() {
					n := index.PostInc() // returns old value
					v, err := z(n)
					if n < STOP && wc.Wait(n) {
						select {
						case c <- C{v, err}:
						case <-stop:
							wc.Stop()
							break loop
						}
						wc.Inc()
					}
				}
				if alive.Dec() == 0 { // return new value
					close(c)
				}
			}()
		}
		return func(index uint64) (reflect.Value, error) {
			if index == STOP {
				close(stop)
				return z(STOP)
			}
			if x, ok := <-c; ok {
				return x.Value, x.error
			}
			return reflect.ValueOf(false), nil
		}
	}
}

func (zf Source) First(n int) Source {
	return func() Stream {
		z := zf()
		count := AtomicCounter{0}
		wc := WaitCounter{Value: 0}
		return func(index uint64) (v reflect.Value, err error) {
			v, err = z(index)
			if index != STOP && wc.Wait(index) {
				if err == nil && v.Kind() != reflect.Bool {
					if count.PostInc() < uint64(n) {
						wc.Inc()
						return
					}
				}
				wc.Stop()
			}
			return reflect.ValueOf(false), nil
		}
	}
}

func (zf Source) Drain(sink func(reflect.Value) error) (err error) {
	z := zf()
	var v reflect.Value
	var i uint64
	for {
		if v, err = z(i); err != nil {
			break
		}
		i++
		if v.Kind() != reflect.Bool {
			if err = sink(v); err != nil {
				break
			}
		} else if !v.Bool() {
			break
		}
	}
	z(STOP)
	e := sink(reflect.ValueOf(err == nil))
	return fu.Fnze(err, e)
}

func Chan(c interface{}, stop ...chan struct{}) Source {
	return func() Stream {
		scase := []reflect.SelectCase{{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(c)}}
		wc := WaitCounter{Value: 0}
		return func(index uint64) (v reflect.Value, err error) {
			if index == STOP {
				wc.Stop()
				for _, s := range stop {
					close(s)
				}
			}
			if wc.Wait(index) {
				_, r, ok := reflect.Select(scase)
				if wc.Inc() && ok {
					return r, nil
				}
			}
			return reflect.ValueOf(false), nil
		}
	}
}

func List(list interface{}) Source {
	return func() Stream {
		v := reflect.ValueOf(list)
		l := uint64(v.Len())
		flag := AtomicFlag{Value: 1}
		return func(index uint64) (reflect.Value, error) {
			if index < l && flag.State() {
				return v.Index(int(index)), nil
			}
			return reflect.ValueOf(false), nil
		}
	}
}

const iniCollectLength = 13

func (z Source) Collect() (r interface{}, err error) {
	length := 0
	values := reflect.ValueOf((interface{})(nil))
	err = z.Drain(func(v reflect.Value) error {
		if length == 0 {
			values = reflect.MakeSlice(reflect.SliceOf(v.Type()), 0, iniCollectLength)
		}
		if v.Kind() != reflect.Bool {
			values = reflect.Append(values, v)
			length++
		}
		return nil
	})
	if err != nil {
		return
	}
	return values.Interface(), nil
}

func (z Source) LuckyCollect() interface{} {
	t, err := z.Collect()
	if err != nil {
		panic(err)
	}
	return t
}

func (z Source) Count() (count int, err error) {
	err = z.Drain(func(v reflect.Value) error {
		if v.Kind() != reflect.Bool {
			count++
		}
		return nil
	})
	return
}

func (z Source) LuckyCount() int {
	c, err := z.Count()
	if err != nil {
		panic(err)
	}
	return c
}

func (zf Source) RandFilter(seed int, prob float64, t bool) Source {
	z := zf()
	return func() Stream {
		nr := fu.NaiveRandom{Value: uint32(seed)}
		wc := WaitCounter{Value: 0}
		return func(index uint64) (v reflect.Value, err error) {
			v, err = z(index)
			if index == STOP {
				wc.Stop()
			}
			if wc.Wait(index) {
				if v.Kind() != reflect.Bool {
					p := nr.Float()
					if (t && p <= prob) || (!t && p > prob) {
						v = reflect.ValueOf(true) // skip
					}
				}
				wc.Inc()
			}
			return
		}
	}
}

func (z Source) RandSkip(seed int, prob float64) Source {
	return z.RandFilter(seed, prob, true)
}

func (z Source) Rand(seed int, prob float64) Source {
	return z.RandFilter(seed, prob, false)
}

func Error(err error) Stream {
	return func(_ uint64) (reflect.Value, error) {
		return reflect.Value{}, err
	}
}
