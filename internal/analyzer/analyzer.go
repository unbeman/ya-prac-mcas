// Package analyzer defines GetMultiAnalyzer function for creating custom linter.
package analyzer

import (
	"github.com/fatih/errwrap/errwrap"
	"github.com/go-critic/go-critic/checkers/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"

	"github.com/unbeman/ya-prac-mcas/internal/analyzer/exitmain"
)

func GetMultiAnalyzer() []*analysis.Analyzer {
	analyzers := []*analysis.Analyzer{
		asmdecl.Analyzer,
		copylock.Analyzer,
		errorsas.Analyzer,
		findcall.Analyzer,
		framepointer.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		structtag.Analyzer,
		tests.Analyzer,
		unreachable.Analyzer,

		errwrap.Analyzer,
		analyzer.Analyzer,

		exitmain.Analyzer,
	}
	analyzers = addStaticCheckAnalyzers(analyzers)
	analyzers = addStyleCheckAnalyzers(analyzers)
	return analyzers
}

func addStaticCheckAnalyzers(analyzers []*analysis.Analyzer) []*analysis.Analyzer {
	for _, v := range staticcheck.Analyzers {
		analyzers = append(analyzers, v.Analyzer)
	}
	return analyzers
}

func addStyleCheckAnalyzers(analyzers []*analysis.Analyzer) []*analysis.Analyzer {
	for _, v := range stylecheck.Analyzers {
		analyzers = append(analyzers, v.Analyzer)
	}
	return analyzers
}
