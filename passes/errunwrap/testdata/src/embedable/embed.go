package embedable

type EmbeddableError struct { // want `error type EmbeddableError does not implement Unwrap\(\) error`
	anError error
}

func (e *EmbeddableError) Error() string {
	return e.anError.Error()
}

type EmbeddableError2 struct { // want `error type EmbeddableError2 does not implement Unwrap\(\) error`
	EmbeddableError
}

type EmbeddableError3 struct {
	EmbeddableError2
}

func (e *EmbeddableError3) Unwrap() error {
	return e.anError
}

type EmbeddableError4 struct {
	EmbeddableError3
}

type FieldError struct { // want `error type FieldError does not implement Unwrap\(\) error`
	aField EmbeddableError
}

func (e FieldError) Error() string {
	return e.aField.Error()
}

type thingWhichWrapsAnError struct {
	anError error
}

type NotAnError struct {
	aField thingWhichWrapsAnError
}

func (e NotAnError) Error() string {
	return "errMsg"
}

// type delcarations of interfaces are NOT expected to implement Unwrap, only concrete types.
type errInterface error

type moreErrInterface interface {
	Error() string
}

// pointers can't have receivers so they can't implement Unwrap
type concreteError *error
