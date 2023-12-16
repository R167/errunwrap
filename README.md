# errunwrap

Check for errors which wrap other errors without implementing the Unwrap method.

## Why?

[Go 1.13](https://go.dev/blog/go1.13-errors) introduced new features in the `errors` package with
the `Unwrap() error` method. This allows types of the form `type MyCustomError struct { error }` to
attach additional semantic meaning and context to errors. However, to use this functionality
correctly, in addition to the cases [errorlint](https://github.com/polyfloyd/go-errorlint) catches
to ensure error comparison is done using `errors.Is` and `errors.As`, developers must ALSO remember
to implement the `Unwrap() error` (or in go 1.20, `Unwrap() []error` if it wraps multiple errors)

## Example

```go
var ErrRowNotFound = errors.New("record is missing in the DB, but that's okay sometimes")

type MissingUnwrapError { // linter will fail here b/c MissingUnwrapError
  error
}

type StatusError struct {
  Err        error
  StatusCode int
}

func (e StatusError) Unwrap() error {
  return e.Err
}

func (e StatusError) Error() string {
  return fmt.Sprintf("%d: %s", e.StatusCode, e.Err)
}

func handleRequest() int {
  var err error = MissingUnwrapError{
    error: StatusError{
      error: ErrBad
      StatusCode: 418
    }
  }

  if err == nil {
    return 200
  }

  if errors.Is(err, ErrRowNotFound) {
    // UH OH! b/c MissingUnwrapError doesn't implement Unwrap() error, we will NEVER hit this code branch
    // errors.Is is unable to look inside the contents of StatusError.error
    return http.StatusNotFound
  }
  var statusErr StatusError
  if errors.As(err, &StatusError) {
    // OH NO AGAIN! Once again MissingUnwrapError foils us from inspecting the status code :sadpanda:
    return statusErr.StatusCode
  }

  // Fallback to general purpose error code
  return 500
}
```
