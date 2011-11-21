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
		//fmt.Printf("%T : %v\n", decl,decl) // TODO: to delete when it's finished

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

		skipName := make([]bool, len(vSpec.Names)) // for identifiers like "_"
		names := make([]string, 0)                 // identifiers
		values := make([]string, 0)

		// === Names
		// http://golang.org/pkg/go/ast/#Ident || godoc go/ast Ident
		//  Name    string    // identifier name
		for i, v := range vSpec.Names {
			names = append(names, v.Name)

			if v.Name == "_" {
				skipName[i] = true
			}
		}

		// === Values
		// http://golang.org/pkg/go/ast/#Expr || godoc go/ast Expr
		//  type Expr interface

		if len(vSpec.Values) != 0 {
			for i, v := range vSpec.Values {
				var expr string

				val := newValue()
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
		isFirst := true
		for i, v := range names {
			if skipName[i] {
				continue
			}

			if isFirst {
				isFirst = false
				buf.WriteString("const " + v)
			} else {
				buf.WriteString(", " + v)
			}
		}

		// It is possible that it is only an identifier "_"
		if isFirst {
			continue
		}

		isFirst = true
		for i, v := range values {
			if skipName[i] {
				continue
			}

			if isFirst {
				isFirst = false
				buf.WriteString(" = " + v)
			} else {
				buf.WriteString(", " + v)
			}
		}

		buf.WriteString(";\n")
	}
}

// * * *

// Represents a value.
type value struct {
	useIota bool
	*bytes.Buffer
}

// Initializes a new type of "value".
func newValue() *value {
	return &value{false, new(bytes.Buffer)}
}

// Gets the value.
// It throws a panic message for types no added.
func (v *value) getValue(iface interface{}) {
	// type Expr
	switch typ := iface.(type) {

	// http://golang.org/pkg/go/ast/#BasicLit || godoc go/ast BasicLit
	//  Value    string      // literal string
	case *ast.BasicLit:
		v.WriteString(iface.(*ast.BasicLit).Value)

	// http://golang.org/pkg/go/ast/#UnaryExpr || godoc go/ast UnaryExpr
	//  Op    token.Token // operator
	//  X     Expr        // operand
	case *ast.UnaryExpr:
		unaryExpr := iface.(*ast.UnaryExpr)

		v.WriteString(unaryExpr.Op.String())
		v.getValue(unaryExpr.X)

	// http://golang.org/pkg/go/ast/#BinaryExpr || godoc go/ast BinaryExpr
	//  X     Expr        // left operand
	//  Op    token.Token // operator
	//  Y     Expr        // right operand
	case *ast.BinaryExpr:
		binaryExpr := iface.(*ast.BinaryExpr)

		// TODO: calculate expression
		v.getValue(binaryExpr.X)
		v.WriteString(" " + binaryExpr.Op.String() + " ")
		v.getValue(binaryExpr.Y)

	// http://golang.org/pkg/go/ast/#Ident || godoc go/ast Ident
	case *ast.Ident:
		if typ.Name == "iota" {
			v.WriteString("%d")
			v.useIota = true
		} else {
			v.WriteString(typ.Name)
		}

	default:
		panic(fmt.Sprintf("[getValue:default] type: %T, value: %v", iface, iface))
	}
}
