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
