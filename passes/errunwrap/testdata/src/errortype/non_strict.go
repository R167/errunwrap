package errortype

import "errors"

type SimpleUnwrapError []error

func (e SimpleUnwrapError) Error() string {
	return errors.Join(e...).Error()
}

func (e SimpleUnwrapError) Unwrap() error {
	if len(e) == 0 {
		return nil
	}
	return e[0]
}

type MissingUnwrapError []error // want `error type MissingUnwrapError does not implement Unwrap\(\) \[\]error`

func (e MissingUnwrapError) Error() string {
	return errors.Join(e...).Error()
}

type ResolvesNestedErrors struct { // want `error type ResolvesNestedErrors does not implement Unwrap\(\) \[\]error`
	errs MissingUnwrapError
}

func (e ResolvesNestedErrors) Error() string {
	return e.errs.Error()
}
