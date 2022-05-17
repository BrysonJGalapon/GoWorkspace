package expression

func reduce(f func(x, y float64) float64, values []float64) (result float64) {
	if len(values) < 2 {
		panic("can't reduce a slice with less than 2 elements")
	}

	result = values[0]

	for i := 1; i < len(values); i++ {
		result = f(result, values[i])
	}

	return result
}

func add(x, y float64) float64 {
	return x + y
}

func subtract(x, y float64) float64 {
	return x - y
}

func multiply(x, y float64) float64 {
	return x * y
}

func divide(x, y float64) float64 {
	return x / y
}
