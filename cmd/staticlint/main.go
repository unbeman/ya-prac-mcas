// Staticlint is a tool for static analysis of Go programs based on different analyzers:
//
// Some from go/analysis/pass;
//
// All from staticcheck;
//
// All from stylecheck;
//
// External errwrap;
//
// External ruleguard;
//
// Custom check that detects os.Exit calls in main function.
//
// Run with help flag for more details:
//
//	staticlint help

package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/unbeman/ya-prac-mcas/internal/analyzer"
)

func main() {
	multichecker.Main(analyzer.GetMultiAnalyzer()...)
}
