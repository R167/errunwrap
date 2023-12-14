package errortype

type StructsAllTheWayDown struct { // want `error type StructsAllTheWayDown does not implement Unwrap\(\) \[\]error`
	errs []StructsAllTheWayDown
}

func (e StructsAllTheWayDown) Error() string {
	return "errMsg"
}

// TODO(R167): We don't examine nested structs in this case. We should, but it gets tricky
// because we need to make sure we don't get stuck in a recursive loop + analyzing all the
// children of an explicitly wrapped struct could get "weird".
type NestedInlineStructErrors struct {
	errs []struct {
		err error
	}
}

func (e NestedInlineStructErrors) Error() string {
	return "errMsg"
}
