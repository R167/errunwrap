package genericssupport

type MagicWrapper[T any] struct {
	val T
}

func (m MagicWrapper[T]) String() T {
	return m.val
}

type UhOhError struct { // want `error type UhOhError does not implement Unwrap\(\) error`
	MagicWrapper[error]
}

func (e *UhOhError) Error() string {
	return e.val.Error()
}
