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
	bufVar := new(bytes.Buffer)

	// If Go sintaxis is incorrect then there will be an error.
	node, err := parser.ParseFile(token.NewFileSet(), filename, nil, 0) //parser.ParseComments)
	if err != nil {
		return err
	}

	for _, decl := range node.Decls {
		//fmt.Printf("%T : %v\n", decl,decl) // TODO: to delete when it's finished

		switch decl.(type) {
		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)

			switch genDecl.Tok {
			case token.CONST:
				getConst(bufConst, genDecl.Specs)
			case token.VAR:
				getVar(bufVar, genDecl.Specs)
			}
		}
	}

	fmt.Print(bufConst.String())
	fmt.Print(bufVar.String())
	return nil
}

//
// === Get

// Maximum number of expressions to get.
// The expressions are the values after of "=".
const MAX_EXPRESSION = 10

// Constants
//
// http://golang.org/doc/go_spec.html#Constant_declarations
// https://developer.mozilla.org/en/JavaScript/Reference/Statements/const
func getConst(buf *bytes.Buffer, spec []ast.Spec) {
	iotas := make([]int, MAX_EXPRESSION)
	lastValues := make([]string, MAX_EXPRESSION)

	// http://golang.org/pkg/go/ast/#ValueSpec || godoc go/ast ValueSpec
	//  Names   []*Ident      // value names (len(Names) > 0)
	//  Values  []Expr        // initial values; or nil
	for _, s := range spec {
		vSpec := s.(*ast.ValueSpec)

		if len(vSpec.Values) > MAX_EXPRESSION {
			panic("length of 'iotas' is lesser than 'vSpec.Values'")
		}

		names, skipName := getName(vSpec)
		values := make([]string, 0)

		// === Values
		// http://golang.org/pkg/go/ast/#Expr || godoc go/ast Expr
		//  type Expr interface
		if len(vSpec.Values) != 0 {
			for i, v := range vSpec.Values {
				var expr string

				val := newValue(names[i])
				val.getValue(v)

				if val.useIota {
					expr = fmt.Sprintf(val.String(), iotas[i])
					iotas[i]++
				} else {
					expr = val.String()
				}

				if !skipName[i] {
					values = append(values, expr)
					lastValues[i] = val.String()
				}
			}
		} else { // get last value of iota
			for i := 0; i < len(names); i++ {
				expr := fmt.Sprintf(lastValues[i], iotas[0])
				values = append(values, expr)
			}
			iotas[0]++
		}

		// === Write
		// TODO: calculate expression using "exp/types"
		isFirst := true
		for i, v := range names {
			if skipName[i] {
				continue
			}

			if isFirst {
				isFirst = false
				buf.WriteString(fmt.Sprintf("const %s = %s", v, values[i]))
			} else {
				buf.WriteString(fmt.Sprintf(", %s = %s", v, values[i]))
			}
		}

		// It is possible that there is only a blank identifier
		if isFirst {
			continue
		}

		buf.WriteString(";\n")
	}
}

// Variables
//
// http://golang.org/doc/go_spec.html#Variable_declarations
// https://developer.mozilla.org/en/JavaScript/Reference/Statements/var
// https://developer.mozilla.org/en/JavaScript/Reference/Statements/let
//
// TODO: use let for local variables
func getVar(buf *bytes.Buffer, spec []ast.Spec) {
	// http://golang.org/pkg/go/ast/#ValueSpec || godoc go/ast ValueSpec
	for _, s := range spec {
		vSpec := s.(*ast.ValueSpec)
		names, skipName := getName(vSpec)
		values := make([]string, 0)

		// === Values
		// http://golang.org/pkg/go/ast/#Expr || godoc go/ast Expr
		for i, v := range vSpec.Values {
			// Skip when it is not a function
			if skipName[i] {
				if _, ok := v.(*ast.CallExpr); !ok {
					continue
				}
			}

			val := newValue(names[i])
			val.getValue(v)
			if !skipName[i] {
				values = append(values, val.String())
			}
		}

		// === Write
		// TODO: calculate expression using "exp/types"
		isFirst := true
		for i, n := range names {
			if skipName[i] {
				continue
			}

			if isFirst {
				isFirst = false
				buf.WriteString("var " + n)
			} else {
				buf.WriteString(", " + n)
			}

			if len(values) != 0 {
				buf.WriteString(" = " + values[i])
			}
		}

		buf.WriteString(";\n")
	}
}
