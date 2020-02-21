package fu

import (
	"reflect"
)

func Option(t interface{}, o []interface{}) reflect.Value {
	tv := reflect.ValueOf(t)
	for _, x := range o {
		v := reflect.ValueOf(x)
		if v.Type() == tv.Type() {
			return v
		}
	}
	return tv
}

func IfsOption(t interface{}, o []interface{}) interface{} {
	return Option(t, o).Interface()
}

func StrOption(t interface{}, o []interface{}) string {
	return Option(t, o).String()
}

func IntOption(t interface{}, o []interface{}) int {
	return int(Option(t, o).Int())
}

func FloatOption(t interface{}, o []interface{}) float64 {
	return Option(t, o).Float()
}

func BoolOption(t interface{}, o []interface{}) bool {
	return Option(t, o).Bool()
}

func MultiOption(o []interface{}, t ...interface{}) (reflect.Value, int) {
	for _, x := range o {
		for i, tv := range t {
			v := reflect.ValueOf(x)
			if v.Type() == reflect.TypeOf(tv) {
				return v, i
			}
		}
	}
	return reflect.ValueOf(t[0]), 0
}

func StrMultiOption(o []interface{}, t ...interface{}) (string, int) {
	v, i := MultiOption(o, t...)
	return v.String(), i
}

func AllStrOptions(o []interface{}, t ...interface{}) []string {
	r := []string{}
	for _, x := range o {
		for _, tv := range t {
			v := reflect.ValueOf(x)
			if v.Type() == reflect.TypeOf(tv) {
				r = append(r, v.String())
			}
		}
	}
	return r
}
