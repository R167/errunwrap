package embedable

type EmbeddableError struct { // want `error type EmbeddableError does not implement Unwrap method`
	anError error
}

func (e *EmbeddableError) Error() string {
	return e.anError.Error()
}

type EmbeddableError2 struct { // want `error type EmbeddableError2 does not implement Unwrap method`
	EmbeddableError
}

type EmbeddableError3 struct {
	EmbeddableError2
}

func (e *EmbeddableError3) Unwrap() error {
	return e.anError
}

type FieldError struct { // want `error type FieldError does not implement Unwrap method`
	aField EmbeddableError
}

func (e FieldError) Error() string {
	return e.aField.Error()
}

type basicErrorWrapper error

type NotAnError struct {
	aField basicErrorWrapper
}

func (e NotAnError) Error() string {
	return "not an error"
}
