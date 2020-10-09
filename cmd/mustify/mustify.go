package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"
)

var src = `package foo

import (
	"fmt"
	"time"
)

func (f int) bar(size int) (bool, error) {
	fmt.Println(time.Now())
	return true
}`

func printFuncDecl(decl ast.Decl) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, token.NewFileSet(), decl)

	out := buf.String()
	head := strings.Split(out, "\n")[0]

	return head
}

// removeError removes last return value if it is an error
func removeError(fn *ast.FuncDecl) {
	results := fn.Type.Results
	lr := len(results.List)
	if lr > 0 {
		last := results.List[lr-1].Type.(*ast.Ident)
		if last.Name == "error" {
			results.List = results.List[:lr-1]
		}
	}
}

func main() {
	fset := token.NewFileSet() // positions are relative to fset

	// Parse src but stop after processing the imports.
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the imports from the file's AST.
	for _, s := range f.Decls {
		fn, ok := s.(*ast.FuncDecl)
		if !ok {
			continue
		}

		//fmt.Println(fn.Name, fn.Type.Params.List[0].Names[0].Name)
		removeError(fn)
		fmt.Println(printFuncDecl(s))

	}

}
