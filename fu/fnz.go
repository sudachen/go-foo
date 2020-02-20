package fu

import "reflect"

/*
Fnz returns the first non zero value
*/
func Fnz(a ...interface{}) interface{} {
	for _, i := range a {
		if !reflect.ValueOf(i).IsZero() {
			return i
		}
	}
	return 0
}

/*
Fnzi returns the first non integer zero value
*/
func Fnzi(a ...int) int {
	for _, i := range a {
		if i != 0 {
			return i
		}
	}
	return 0
}

/*
Fnzl returns the first non zero long integer value
*/
func Fnzl(a ...int64) int64 {
	for _, i := range a {
		if i != 0 {
			return i
		}
	}
	return 0
}

/*
Fnzf returns the first non zero float value
*/
func Fnzf(a ...float32) float32 {
	for _, i := range a {
		if i != 0 {
			return i
		}
	}
	return 0
}

/*
Fnzd returns the first non zero double value
*/
func Fnzd(a ...float64) float64 {
	for _, i := range a {
		if i != 0 {
			return i
		}
	}
	return 0
}
