package utilities

// TODO: refactor this as a generic
func Min(values ...int) int {
	if len(values) == 0 {
		panic("Minimum value can only be computed on non empty slices/arrays")
	} else if len(values) == 1 {
		return values[0]
	}

	minValue := values[0]
	for value := range values[1:] {
		if value < minValue {
			minValue = value
		}
	}

	return minValue
}
