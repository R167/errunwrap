// Check for error types that do not implement the Unwrap method
package errunwrap

import (
	"flag"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "errunwrap",
	Doc:      "report error types without an Unwrap method",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

var (
	flagSet      flag.FlagSet
	strictUnwrap bool
)

func init() {
	flagSet.BoolVar(&strictUnwrap, "strict-unwrap", false, "Strictly check that the Unwrap method returns an []error if multiple errors are wrapped")
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.TypeSpec)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch stmt := n.(type) {
		case *ast.TypeSpec:
			if stmt.Name == nil {
				return // ignore anonymous types
			}
			if stmt.Assign != 0 {
				return // ignore type aliases
			}

			// Get the type information
			info := pass.TypesInfo.Defs[stmt.Name]
			if info == nil {
				return
			}

			if !isErrorType(info.Type()) {
				return
			}

			errKind := wrapsError(pass, stmt.Type)
			if errKind == none {
				return
			}

			mustArray := errKind == array && strictUnwrap
			var expected string
			if errKind == array {
				// while we only fail on strictUnwrap, we still want to report the correct type
				expected = "[]error"
			} else {
				expected = "error"
			}

			unwrap, ok := unwrapMethod(info.Type())
			if !ok {
				pass.Reportf(stmt.Pos(), "error type %s does not implement Unwrap() %s", stmt.Name, expected)
				return
			}

			ret := unwrap.Type().(*types.Signature).Results()
			if ret.Len() != 1 {
				pass.Reportf(unwrap.Pos(), "Expected %s to return %s, got %s", unwrap.Name(), expected, ret)
				return
			}

			output := ret.At(0).Type().String()

			if !mustArray && output == "error" {
				return
			} else if output == "[]error" {
				return
			}

			pass.Reportf(unwrap.Pos(), "Expected %s to return %s, got %s", unwrap.Name(), expected, output)
		}
	})

	return nil, nil
}

var errType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

func isErrorType(t types.Type) bool {
	return types.Implements(t, errType) || types.Implements(types.NewPointer(t), errType)
}

// unwrapMethod checks if the type has an Unwrap\(\) error
func unwrapMethod(t types.Type) (*types.Func, bool) {
	// We always get a concrete type here, however we want to check the pointer receiver as well
	methods := types.NewMethodSet(types.NewPointer(t))

	for i := 0; i < methods.Len(); i++ {
		method := methods.At(i).Obj().(*types.Func)

		// Check if the method is the Unwrap method
		if method.Name() == "Unwrap" {
			return method, true
		}
	}
	return nil, false
}

// wrapsError checks if the type wraps or contains an error or []error
func wrapsError(pass *analysis.Pass, ts ast.Expr) errorCount {
	return wrapsErrorEmbed(pass, pass.TypesInfo.TypeOf(ts), true, false)
}

type errorCount int

const (
	none errorCount = iota
	one
	array
	// arguably, there's a place for a "many" here, but it's pretty niche
)

// wrapsErrorEmbed checks if the type wraps or contains an error or []error.
// When checkStructs is true, deep inspect each of its fields for errors, otherwise only check structs at the top level.
// This allows recursively checking embedded structs for errors while skipping non-embedded fields.
//
// Returns if the type wraps an error, or if it wraps an array of errors.
func wrapsErrorEmbed(pass *analysis.Pass, t types.Type, checkStructs bool, allowError bool) errorCount {
	switch t := t.(type) {
	case nil, *types.Basic:
		return none
	case *types.Named:
		// types.Named is special. I think it's just b/c we hit the `error` type and then it blows up, but I'm not sure.
		// Ultimately it's nice to be able to check if the underlying error that's wrapped is actually a []error or just error
		// so we can give a better message.
		if v := wrapsErrorEmbed(pass, t.Underlying(), checkStructs, allowError); v != none {
			return v
		}
		return defaultErrorCheck(t, allowError)
	case *types.Pointer:
		return wrapsErrorEmbed(pass, t.Elem(), checkStructs, allowError)
	case *types.Slice:
		if wrapsErrorEmbed(pass, t.Elem(), checkStructs, true) == none {
			return none
		}
		return array
	case *types.Struct:
		if !checkStructs {
			if isErrorType(t) {
				return one
			}
			return none
		}

		found := 0
		for i := 0; i < t.NumFields(); i++ {
			field := t.Field(i)
			recurse := field.Embedded() && checkStructs
			toCheck := field.Type()
			if field.Embedded() {
				toCheck = field.Type().Underlying()
			}
			switch wrapsErrorEmbed(pass, toCheck, recurse, true) {
			case array:
				return array
			case one:
				found++
			}
		}
		if found > 0 {
			return one
		}
	default:
		return defaultErrorCheck(t, allowError)
	}

	return none
}

func defaultErrorCheck(t types.Type, allowError bool) errorCount {
	if !allowError && types.IsInterface(t) {
		return none
	}
	if isErrorType(t) {
		return one
	}
	return none
}
