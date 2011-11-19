// Copyright 2011  The "GotoScript" Authors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package tojs

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// Compiles a Go source file into JavaScript.
func Compile(filename string) error {
	bufConst := new(bytes.Buffer)
	//bufVar := new(bytes.Buffer)

	// If Go sintaxis is incorrect then there will be an error.
	node, err := parser.ParseFile(token.NewFileSet(), filename, nil, 0) //parser.ParseComments)
	if err != nil {
		return err
	}

	for _, decl := range node.Decls {
//	fmt.Printf("%T : %v\n", decl,decl)

		switch decl.(type) {
		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)

			// Constants
			if genDecl.Tok == token.CONST {
				getConst(bufConst, genDecl.Specs)
			}
		}
	}

	fmt.Print(bufConst.String())
	return nil
}

// Constants
//
// http://golang.org/doc/go_spec.html#Constant_declarations
// https://developer.mozilla.org/en/JavaScript/Reference/Statements/const
func getConst(buf *bytes.Buffer, spec []ast.Spec) {
	// http://golang.org/pkg/go/ast/#ValueSpec || godoc go/ast ValueSpec
	//   Names   []*Ident      // value names (len(Names) > 0)
	//   Values  []Expr        // initial values; or nil
	for _, s := range spec {
		buf.WriteString("const ")
		vSpec := s.(*ast.ValueSpec)

		// === Names
		// http://golang.org/pkg/go/ast/#Ident || godoc go/ast Ident
		//   Name    string    // identifier name
		for i, v := range vSpec.Names {
			if i != 0 {
				buf.WriteString(", " + v.Name)
			} else {
				buf.WriteString(v.Name)
			}
		}

		buf.WriteString(" = ")

		// === Values
		// http://golang.org/pkg/go/ast/#Expr || godoc go/ast Expr
		//   type Expr interface
		for i, v := range vSpec.Values {
			if i != 0 {
				buf.WriteString(", ")
			}
			getType(buf, v)
		}

		buf.WriteString(";\n")
	}
}

// Gets the type.
// It throws a panic message for types no added.
func getType(buf *bytes.Buffer, i interface{}) {
	// type Expr
	switch typ := i.(type) {

	// http://golang.org/pkg/go/ast/#BasicLit || godoc go/ast BasicLit
	//   Value    string      // literal string
	case *ast.BasicLit:
		buf.WriteString(i.(*ast.BasicLit).Value)

	// http://golang.org/pkg/go/ast/#UnaryExpr || godoc go/ast UnaryExpr
	//   Op    token.Token // operator
	//   X     Expr        // operand
	case *ast.UnaryExpr:
		unaryExpr := i.(*ast.UnaryExpr)
		buf.WriteString(unaryExpr.Op.String())
		getType(buf, unaryExpr.X)

	// http://golang.org/pkg/go/ast/#Ident || godoc go/ast Ident
	case *ast.Ident:
		if typ.Name == "iota" {
			buf.WriteString("0")
		} else {
			panic(fmt.Sprintf("[getType:*ast.Ident] name: %v", typ.Name))
		}

	default:
		panic(fmt.Sprintf("[getType:default] type: %T, value: %v", i, i))
	}
}
