package utilities

// Min TODO: refactor this as a generic
// Min defines our own implementation of Min since, apparently, golang does not offer its own builtin
// Accepts an arbitrary amount of int values through variadic argument.
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
