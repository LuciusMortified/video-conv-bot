package ptr

func AnyRef[T any](value T) *T {
	return &value
}
