// +build ignore

// Unmain reads a go file from stdin and removes the main function.
package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"

	"golang.org/x/tools/imports"
)

func main() {
	del := make(map[string]bool)
	for _, f := range append([]string{"main"}, os.Args[1:]...) {
		del[f] = true
	}

	fset := token.NewFileSet()
	f, e := parser.ParseFile(fset, "", os.Stdin, parser.ParseComments)
	fatal(e)

	// remove main function
	j := 0
	for _, a := range f.Decls {
		if d, o := a.(*ast.FuncDecl); o && del[d.Name.Name] {
		} else {
			f.Decls[j] = a
			j++
		}
	}
	f.Decls = f.Decls[:j]

	// fix imports
	var buf bytes.Buffer
	fatal(format.Node(&buf, fset, f))
	out, e := imports.Process("k.go", buf.Bytes(), nil)
	fatal(e)
	os.Stdout.Write(out)
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
