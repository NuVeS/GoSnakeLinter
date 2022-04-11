// проверяет чтоб не было записи змейкой_как_эта
package snakecheck

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "snakelint",
	Doc:  "reports snake_case",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			ident, ok := n.(*ast.Ident)
			if !ok {
				return true
			}

			if len(ident.Name) <= 0 {
				return false
			}

			containsSnake := func(expr string) bool {
				return strings.Contains(expr, "_")
			}

			if containsSnake(ident.Name) {
				pass.Reportf(ident.Pos(), "snake found %q",
					render(pass.Fset, ident))
				return true
			}

			return false
		})
	}

	return nil, nil
}

func render(fset *token.FileSet, x interface{}) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, x); err != nil {
		panic(err)
	}
	return buf.String()
}
