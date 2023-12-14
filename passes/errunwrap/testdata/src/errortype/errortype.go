package errortype

import "errors"

// Should be ignored
type AliasError = error

// don't check interfaces
type KindOfError error

type ArrayError []error

type StructError struct {
	error
}

type StructArrayError struct {
	errs []error
}

type GenericWithError[T any] struct {
	myErr error
	Ctx   T
}

func (e ArrayError) Error() string {
	return errors.Join(e...).Error()
}

func (e ArrayError) Unwrap() []error {
	return e
}

func (e *StructError) Unwrap() error {
	return e.error
}

func (e *StructArrayError) Error() string {
	return errors.Join(e.errs...).Error()
}

func (e *StructArrayError) Unwrap() []error {
	return e.errs
}

func (e *GenericWithError[T]) Error() string {
	return e.myErr.Error()
}

func (e *GenericWithError[T]) Unwrap() error {
	return e.myErr
}
