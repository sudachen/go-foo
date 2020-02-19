package fu

func Ife(expr bool, x interface{}, y interface{}) interface{} {
	if expr {
		return x
	}
	return y
}

func Ifei(expr bool, x int, y int) int {
	if expr {
		return x
	}
	return y
}

func Ifei64(expr bool, x int64, y int64) int64 {
	if expr {
		return x
	}
	return y
}

func Ifef(expr bool, x float32, y float32) float32 {
	if expr {
		return x
	}
	return y
}

func Ifef64(expr bool, x float64, y float64) float64 {
	if expr {
		return x
	}
	return y
}
