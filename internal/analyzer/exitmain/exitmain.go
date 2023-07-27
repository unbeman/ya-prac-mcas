// Package exitmain defines Analyzer that checks using os.Exit calls in main function of main package.
package exitmain

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const (
	mainName     = "main"
	exitFuncName = "Exit"
	exitPkgName  = "os"
)

var Analyzer = &analysis.Analyzer{
	Name:     "exitmain",
	Doc:      "reports usage of os.Exit in main function",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != mainName { // package main
		return nil, nil
	}

	var mainFunc *ast.FuncDecl

fileLoop:
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Name.Name == mainName { // func main()
				mainFunc = funcDecl
				break fileLoop
			}
		}
	}

	if mainFunc == nil {
		return nil, nil
	}

	for _, stmt := range mainFunc.Body.List {
		if expr, isExpr := stmt.(*ast.ExprStmt); isExpr {
			if call, isCall := expr.X.(*ast.CallExpr); isCall {
				if selExpr, isSelExpr := call.Fun.(*ast.SelectorExpr); isSelExpr {
					if ident, isIdent := selExpr.X.(*ast.Ident); isIdent {
						if selExpr.Sel.Name == exitFuncName && ident.Name == exitPkgName { // os.Exit()
							pass.Reportf(call.Lparen, "call of os.Exit in main function")
						}
					}
				}
			}
		}
	}

	return nil, nil
}
