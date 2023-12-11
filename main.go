package main

import (
	"github.com/R167/errunwrap/passes/errunwrap"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(errunwrap.Analyzer)
}
