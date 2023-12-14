package errunwrap

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestErrorType(t *testing.T) {
	t.Parallel()
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer,
		"errortype",
		"genericssupport",
		"embedable",
	)
}

func TestStrictMode(t *testing.T) {
	testdata := analysistest.TestData()
	strictUnwrap = true
	defer func() { strictUnwrap = false }()
	analysistest.Run(t, testdata, Analyzer,
		"strict",
	)
}

func TestWrapsError(t *testing.T) {
	t.Parallel()
}
