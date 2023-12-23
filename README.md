# errunwrap

[![CI](https://github.com/R167/errunwrap/actions/workflows/ci.yaml/badge.svg)](https://github.com/R167/errunwrap/actions/workflows/ci.yaml) [![Go Reference](https://pkg.go.dev/badge/github.com/R167/errunwrap.svg)](https://pkg.go.dev/github.com/R167/errunwrap) [![Go Report Card](https://goreportcard.com/badge/github.com/R167/errunwrap)](https://goreportcard.com/report/github.com/R167/errunwrap)

Check for errors which wrap other errors without implementing the Unwrap method.

[Go 1.13](https://go.dev/blog/go1.13-errors) introduced new features in the `errors` package with
the `Unwrap() error` method. This allows types of the form `type MyCustomError struct { error }` to
attach additional semantic meaning and context to errors. However, to use this functionality
correctly, in addition to the cases [errorlint](https://github.com/polyfloyd/go-errorlint) catches
to ensure error comparison is done using `errors.Is` and `errors.As`, developers must ALSO remember
to implement the `Unwrap() error` (or in go 1.20, `Unwrap() []error` if it wraps multiple errors)

## Usage

```bash
go install github.com/R167/errunwrap@latest

# Run the linter. Accepts standard go package specs
errunwrap ./...

# Run the linter in strict mode. This ensures any errors which wrap multiple
# errors implement Unwrap() []error
errunwrap -strict-unwrap ./...
```

If you want to run this along with other linters, you can use the analyzer itself
`github.com/R167/errunwrap/passes/errunwrap`. Refer to golang.org/x/tools/go/analysis for more
information on how to use analyzers.

## Examples

### Missing Unwrap

Errors which wrap other errors without implementing the `Unwrap() error` method will cause
`errors.Is` and `errors.As` to fail to traverse the error chain.

```go
type MissingUnwrapError { // linter will fail b/c Unwrap() is not implemented
  error
}

// errors.Is will fail to traverse the error chain
if errors.Is(MissingUnwrapError{error: ErrRowNotFound}, ErrRowNotFound) {
  // This code will never be reached
}
```

### Strict Unwrap

When using `-strict-unwrap`, the linter will enforce that errors which wrap multiple errors
(e.g. `type MultiError []error`) implement `Unwrap() []error` instead of `Unwrap() error`. This
allows `errors.Is` and `errors.As` to traverse the full error tree.

```go
// with -strict-unwrap
type MultiError []error

func (e MultiError) Unwrap() error { // linter will fail b/c Unwrap() error
  return e[0]
}
```

### Full example

```go
var ErrRowNotFound = errors.New("record is missing in the DB, but that's okay sometimes")

type MissingUnwrapError { // linter will fail b/c Unwrap() is not implemented
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
    // UH OH! b/c MissingUnwrapError doesn't implement Unwrap() error,
    // we will NEVER hit this code branch
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
