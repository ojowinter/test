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

		skipName := make([]bool, len(vSpec.Names)) // for blank identifiers "_"
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

		// It is possible that there is only a blank identifier
		if isFirst {
			continue
		}

		isFirst = true
		for i, v := range values {
			if skipName[i] {
				continue
			}

			// TODO: calculate expression using "exp/types"
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

		skipName := make([]bool, len(vSpec.Names))
		names := make([]string, 0)
		values := make([]string, 0)

		// === Names
		// http://golang.org/pkg/go/ast/#Ident || godoc go/ast Ident
		for i, v := range vSpec.Names {
			// Mark blank identifier
			if v.Name == "_" {
				skipName[i] = true
				continue
			}
			names = append(names, v.Name)
		}

		// === Values
		// http://golang.org/pkg/go/ast/#Expr || godoc go/ast Expr
		for i, v := range vSpec.Values {
			//var expr string

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
		isFirst := true
		for i, v := range names {
			if skipName[i] {
				continue
			}

			if isFirst {
				isFirst = false
				buf.WriteString("var " + v)
			} else {
				buf.WriteString(", " + v)
			}
		}

		isFirst = true
		for i, v := range values {
			if skipName[i] {
				continue
			}

			// TODO: calculate expression using "exp/types"
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
	ident   string // variable's identifier
	lit     string // store the last literal
	*bytes.Buffer
}

// Initializes a new type of "value".
func newValue(identifier string) *value {
	return &value{false, identifier, "", new(bytes.Buffer)}
}

// Gets the value.
// It throws a panic message for types no added.
func (v *value) getValue(iface interface{}) {
	// type Expr
	switch typ := iface.(type) {

	// http://golang.org/pkg/go/ast/#BasicLit || godoc go/ast BasicLit
	//  Value    string      // literal string
	case *ast.BasicLit:
		lit := iface.(*ast.BasicLit).Value

		v.WriteString(lit)
		v.lit = lit

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

	// http://golang.org/pkg/go/ast/#CompositeLit || godoc go/ast CompositeLit
	//  Type   Expr      // literal type; or nil
	//  Elts   []Expr    // list of composite elements; or nil
	case *ast.CompositeLit:
		composite := iface.(*ast.CompositeLit)

		v.getValue(composite.Type)
		fmt.Println(composite.Elts)

	// http://golang.org/pkg/go/ast/#ArrayType || godoc go/ast ArrayType
	//  Len    Expr      // Ellipsis node for [...]T array types, nil for slice types
	//  Elt    Expr      // element type
	case *ast.ArrayType:
		array := iface.(*ast.ArrayType)

		if v.lit == "" {
			v.WriteString("new Array(")
			v.getValue(array.Len)
			v.WriteString(")")
		} else {
			v.WriteString(fmt.Sprintf("; for (i=0; i<%s; i++) %s[i]=new Array(",
				v.lit, v.ident))
			v.getValue(array.Len)
			v.WriteString(")")
		}

		if _, ok:= array.Elt.(*ast.ArrayType); ok {
			v.getValue(array.Elt)
		}

	default:
		panic(fmt.Sprintf("[getValue:default] unimplemented: %T, value: %v",
			iface, iface))
	}
}
