// Copyright 2011  The "GoScript" Authors
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

package goscript

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

// Represents an expression.
type expression struct {
	ident      string // variable's name
	useIota    bool
	isNegative bool

	eltsLen int      // store length of array, to use in case of ellipsis (...)
	lit     []string // store the last literals (for array)

	*bytes.Buffer // sintaxis translated
}

// Initializes a new type of "expression".
func newExpression(identifier string) *expression {
	return &expression{
		identifier,
		false,
		false,
		0,
		make([]string, 0),
		new(bytes.Buffer),
	}
}

// Returns the Go expression in JavaScript.
func getExpression(ident string, expr ast.Expr) string {
	e := newExpression(ident)

	e.transform(expr)
	return e.String()
}

// Returns the values of an array formatted like "[i0][i1]..."
func (e *expression) printArray() string {
	a := ""

	for i := 0; i < len(e.lit); i++ {
		vArray := "i" + strconv.Itoa(i)
		a = fmt.Sprintf("%s[%s]", a, vArray)
	}
	return a
}

// Transforms the Go expression.
// It throws a panic message for types no added.
func (e *expression) transform(expr ast.Expr) {
	switch typ := expr.(type) {

	// http://golang.org/pkg/go/ast/#ArrayType || godoc go/ast ArrayType
	//  Len    Expr      // Ellipsis node for [...]T array types, nil for slice types
	//  Elt    Expr      // element type
	case *ast.ArrayType:
		if typ.Len == nil { // slice
			break
		}

		if len(e.lit) == 0 {
			e.WriteString("new Array(")

			if e.eltsLen != 0 { // ellipsis
				e.WriteString(strconv.Itoa(e.eltsLen))
			} else {
				e.transform(typ.Len)
			}

			e.WriteString(")")
		} else {
			iArray := len(e.lit) - 1             // index of array
			vArray := "i" + strconv.Itoa(iArray) // variable's name for the loop

			e.WriteString(fmt.Sprintf("; for (var %s=0; %s<%s; %s++){ %s%s=new Array(",
				vArray, vArray, e.lit[iArray], vArray, e.ident, e.printArray()))
			e.transform(typ.Len)
			e.WriteString(")")
		}

		if _, ok := typ.Elt.(*ast.ArrayType); ok {
			e.transform(typ.Elt)
		} else if len(e.lit) > 1 {
			e.WriteString(";" + SP + strings.Repeat("}", len(e.lit)-1))
		}

	// http://golang.org/pkg/go/ast/#BasicLit || godoc go/ast BasicLit
	//  Kind     token.Token // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING
	//  Value    string      // literal string
	case *ast.BasicLit:
		e.WriteString(typ.Value)

		// === Add the value
		sign := ""

		if e.isNegative {
			sign = "-"
		}
		e.lit = append(e.lit, sign+typ.Value)

	// http://golang.org/pkg/go/ast/#BinaryExpr || godoc go/ast BinaryExpr
	//  X     Expr        // left operand
	//  Op    token.Token // operator
	//  Y     Expr        // right operand
	case *ast.BinaryExpr:
		e.transform(typ.X)
		e.WriteString(SP + typ.Op.String() + SP)
		e.transform(typ.Y)

	// http://golang.org/pkg/go/ast/#CallExpr || godoc go/ast CallExpr
	//  Fun      Expr      // function expression
	//  Args     []Expr    // function arguments; or nil
	case *ast.CallExpr:
		callIdent := typ.Fun.(*ast.Ident).Name

		// Conversion: []byte()
		if t, ok := typ.Fun.(*ast.ArrayType); ok {
			if t.Elt.(*ast.Ident).Name == "byte" {
				e.transform(typ.Args[0])
			} else {
				panic(fmt.Sprintf("[getValue] call of conversion unimplemented: []%T()", t))
			}
			break
		}

		switch callIdent {
		default:
			panic(fmt.Sprintf("[getValue] call unimplemented: %s", callIdent))

		case "make":
			switch argType := typ.Args[0].(type) {
			default:
				panic(fmt.Sprintf("[getValue] call of 'make' unimplemented: %T", argType))

			// For slice
			case *ast.ArrayType:
				e.WriteString("new Array(")
				e.transform(typ.Args[len(typ.Args)-1]) // capacity
				e.WriteString(")")

			// The second argument (in Args), if any, is the capacity which
			// is not useful in JS since it is dynamic.
			case *ast.MapType:
				e.WriteString("{};") // or "new Object()"
			}

		case "new":
			switch argType := typ.Args[0].(type) {
			default:
				panic(fmt.Sprintf("[getValue] call of 'new' unimplemented: %T", argType))

			case *ast.ArrayType:
				for _, arg := range typ.Args {
					e.transform(arg)
				}
			}

		// Conversion
		case "uint", "uint8", "uint16", "uint32",
			"int", "int8", "int16", "int32",
			"float32", "float64", "byte", "rune", "string":
			e.transform(typ.Args[0])
		}

	// http://golang.org/pkg/go/ast/#CompositeLit || godoc go/ast CompositeLit
	//  Type   Expr      // literal type; or nil
	//  Elts   []Expr    // list of composite elements; or nil
	case *ast.CompositeLit:
		switch compoType := typ.Type.(type) {
		default:
			panic(fmt.Sprintf("[getValue] 'CompositeLit' unimplemented: %s", compoType))

		case *ast.ArrayType:
			e.eltsLen = len(typ.Elts) // for ellipsis
			e.transform(typ.Type)
			//e.pos = 

			// For arrays initialized
			if len(typ.Elts) != 0 {
				if compoType.Len == nil {
					e.WriteString("[")
				} else {
					e.WriteString(fmt.Sprintf(";%s%s%s=%s[", SP, e.ident, SP, SP))
				}

				for i, el := range typ.Elts {
					if i != 0 {
						e.WriteString(",")
					}
					e.transform(el)
				}
				e.WriteString("]")
			}

		// http://golang.org/pkg/go/ast/#MapType || godoc go/ast MapType
		//  Key   Expr
		//  Value Expr
		case *ast.MapType:
			lenElts := len(typ.Elts) - 1
			e.WriteString("{")

			for i, el := range typ.Elts {
				e.transform(el)

				if i != lenElts {
					e.WriteString("," + SP)
				}
			}
			e.WriteString("};")
		}

	// http://golang.org/pkg/go/ast/#Ellipsis || godoc go/ast Ellipsis
	//  Elt      Expr      // ellipsis element type (parameter lists only); or nil
	//case *ast.Ellipsis:

	// http://golang.org/pkg/go/ast/#Ident || godoc go/ast Ident
	//  Name    string    // identifier name
	case *ast.Ident:
		name := typ.Name

		if name == "iota" {
			e.WriteString("%d")
			e.useIota = true
			break
		}
		// Undefined value in array / slice
		if len(e.lit) != 0 && name == "_" {
			break
		}
		// https://developer.mozilla.org/en/JavaScript/Reference/Global_Objects/undefined
		if name == "nil" {
			name = "undefined"
		}

		e.WriteString(name)

	// http://golang.org/pkg/go/ast/#KeyValueExpr || godoc go/ast KeyValueExpr
	//  Key   Expr
	//  Value Expr
	case *ast.KeyValueExpr:
		e.transform(typ.Key)
		e.WriteString(":")
		e.transform(typ.Value)

	// http://golang.org/pkg/go/ast/#ParenExpr || godoc go/ast ParenExpr
	//  X      Expr      // parenthesized expression
	case *ast.ParenExpr:
		e.transform(typ.X)

	// http://golang.org/pkg/go/ast/#StructType || godoc go/ast StructType
	//  Struct     token.Pos  // position of "struct" keyword
	//  Fields     *FieldList // list of field declarations
	//  Incomplete bool       // true if (source) fields are missing in the Fields list
	/*case *ast.StructType:*/

	// http://golang.org/pkg/go/ast/#StarExpr || godoc go/ast StarExpr
	//  X    Expr      // operand
	case *ast.StarExpr:
		e.transform(typ.X)

	// http://golang.org/pkg/go/ast/#UnaryExpr || godoc go/ast UnaryExpr
	//  Op    token.Token // operator
	//  X     Expr        // operand
	case *ast.UnaryExpr:
		switch typ.Op {
		case token.SUB:
			e.isNegative = true
		}

		e.WriteString(typ.Op.String())
		e.transform(typ.X)

	default:
		panic(fmt.Sprintf("[getValue] unimplemented: %T, value: %v",
			expr, expr))
	}
}
