package fu

func Fnzi(a ...int) int {
	for _, i := range a {
		if i != 0 {
			return i
		}
	}
	return 0
}

func Fnzi64(a ...int64) int64 {
	for _, i := range a {
		if i != 0 {
			return i
		}
	}
	return 0
}

func Fnzf(a ...float32) float32 {
	for _, i := range a {
		if i != 0 {
			return i
		}
	}
	return 0
}

func Fnzf64(a ...float64) float64 {
	for _, i := range a {
		if i != 0 {
			return i
		}
	}
	return 0
}
