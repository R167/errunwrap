package strict

import (
	"errors"
	"net/http"
)

type SimpleUnwrapError []error

func (e SimpleUnwrapError) Error() string {
	return errors.Join(e...).Error()
}

func (e SimpleUnwrapError) Unwrap() error { // want `Expected Unwrap to return \[\]error, got error`
	if len(e) == 0 {
		return nil
	}
	return e[0]
}

type CorrectUnwrapError []error

func (e CorrectUnwrapError) Error() string {
	return errors.Join(e...).Error()
}

func (e CorrectUnwrapError) Unwrap() []error {
	return e
}

type WrongSingleError struct {
	error
}

// Currently we let this slide, but we should _maybe_ catch it.
func (e WrongSingleError) Unwrap() []error {
	return []error{e.error}
}

type CorrectSingleError struct {
	error
}

func (e CorrectSingleError) Unwrap() error {
	return e.error
}

type CompletelyWrongError struct {
	error
}

func (e CompletelyWrongError) Unwrap() http.Handler { // want `Expected Unwrap to return error, got net/http.Handler`
	return nil
}
