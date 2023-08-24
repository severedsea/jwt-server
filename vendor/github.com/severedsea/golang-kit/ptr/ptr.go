package ptr

// Reference returns the pointer reference to the value in the argument.
func Reference[T any](v T) *T {
	return &v
}

// Value returns the dereference of the pointer in the argument
func Value[T any](v *T) T {
	if v == nil {
		var r T

		return r
	}

	return *v
}
