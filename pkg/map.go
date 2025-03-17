package pkg

func Map[V, U any](fn func(V) U, values []V) []U {
	result := make([]U, len(values))
	for i, value := range values {
		result[i] = fn(value)
	}
	return result
}
