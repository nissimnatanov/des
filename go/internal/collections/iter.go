package collections

func MapIter[T, S any](si func(yield func(S) bool), convert func(S) T) func(yield func(T) bool) {
	return func(yield func(T) bool) {
		si(func(s S) bool {
			return yield(convert(s))
		})
	}
}

func MapSlice[T, S any](items []S, convert func(S) T) []T {
	result := make([]T, len(items))
	for i, item := range items {
		result[i] = convert(item)
	}
	return result
}
