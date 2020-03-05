package fu

import "reflect"

func MapInterface(m map[string]reflect.Value) map[string]interface{} {
	r := map[string]interface{}{}
	for k, v := range m {
		r[k] = v.Interface()
	}
	return r
}
