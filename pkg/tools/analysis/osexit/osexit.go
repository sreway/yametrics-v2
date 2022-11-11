// Package osexit defines an Analyzer for the presence of the os.Exit function in the main package
package osexit

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const Doc = "Checks for the presence of the os.Exit function in the main package"

var Analyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  Doc,
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if pass.Pkg.Name() != "main" {
			continue
		}
		ast.Inspect(file, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				var (
					id *ast.Ident
					p  string
				)
				switch fun := call.Fun.(type) {
				case *ast.Ident:
					id = fun
				case *ast.SelectorExpr:
					id = fun.Sel
					p = fmt.Sprintf("%v", fun.X)
				}

				if id != nil && id.Name == "Exit" && p == "os" {
					pass.Report(analysis.Diagnostic{
						Pos:     call.Lparen,
						Message: "Found call of os.Exit on main package",
					})
				}
			}
			return true
		})
	}

	return nil, nil
}
