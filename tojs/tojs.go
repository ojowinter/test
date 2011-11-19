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
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
)

// Compiles a Go source file into JavaScript.
func Compile(filename string) error {
	//bufConst := bufio.NewWriter(os.Stdout)
	//bufVar := bufio.NewWriter(os.Stdout)

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
				getConst(genDecl.Specs)
			}
		}
	}

	//bufConst.Flush()
	return nil
}

// Constants
//
// http://golang.org/doc/go_spec.html#Constant_declarations
// https://developer.mozilla.org/en/JavaScript/Reference/Statements/const
func getConst(spec []ast.Spec) {
	// Each "ValueSpec" can have several identifiers with its values
	for _, s := range spec {
		vSpec := s.(*ast.ValueSpec) // spec of token.CONST

fmt.Print("const ")

		// Names
		for i, v := range vSpec.Names {
			if i != 0 {
				fmt.Print(", " + v.Name)
			} else {
				fmt.Print(v.Name)
			}
		}

fmt.Print(" = ")

		// Values
		for i, v := range vSpec.Values {
			if i != 0 {
				fmt.Print(", ")
			}
			getType(v)
		}

fmt.Print(";\n")
	}
}

// Gets the type.
// It throws a panic message for types no added.
func getType(i interface{}) {
	// Get the type
	switch typ := i.(type) {
	case *ast.BasicLit:
fmt.Print(i.(*ast.BasicLit).Value)
	case *ast.UnaryExpr:
		unaryExpr := i.(*ast.UnaryExpr)
fmt.Print(unaryExpr.Op)
		getType(unaryExpr.X)
	case *ast.Ident:
		if typ.Name == "iota" {
			fmt.Print("iota")
		} else {
			panic(fmt.Sprintf("[getType:*ast.Ident] name: %v", typ.Name))
		}
	default:
		panic(fmt.Sprintf("[getType:default] type: %T, value: %v", i, i))
	}
}
