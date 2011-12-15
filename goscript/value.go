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

// Represents a value.
type value struct {
	name       string // variable's name
	useIota    bool
	isNegative bool

	eltsLen int      // store length of array, to use in case of ellipsis (...)
	lit     []string // store the last literals (for array)

	*bytes.Buffer // sintaxis translated
}

// Initializes a new type of "value".
func newValue(identifier string) *value {
	return &value{
		identifier,
		false,
		false,
		0,
		make([]string, 0),
		new(bytes.Buffer),
	}
}

// Returns the values of an array formatted like "[i0][i1]..."
func (v *value) printArray() string {
	a := ""

	for i := 0; i < len(v.lit); i++ {
		vArray := "i" + strconv.Itoa(i)
		a = fmt.Sprintf("%s[%s]", a, vArray)
	}
	return a
}

// * * *

// Gets the value.
// It throws a panic message for types no added.
func (tr *transform) getValue(expr ast.Expr) {
	// type Expr
	switch typ := expr.(type) {

	// http://golang.org/pkg/go/ast/#ArrayType || godoc go/ast ArrayType
	//  Len    Expr      // Ellipsis node for [...]T array types, nil for slice types
	//  Elt    Expr      // element type
	case *ast.ArrayType:
		if typ.Len == nil { // slice
			break
		}

		if len(tr.src.lit) == 0 {
			tr.src.WriteString("new Array(")

			if tr.src.eltsLen != 0 { // ellipsis
				tr.src.WriteString(strconv.Itoa(tr.src.eltsLen))
			} else {
				tr.getValue(typ.Len)
			}

			tr.src.WriteString(")")
		} else {
			iArray := len(tr.src.lit) - 1        // index of array
			vArray := "i" + strconv.Itoa(iArray) // variable's name for the loop

			tr.src.WriteString(fmt.Sprintf("; for (var %s=0; %s<%s; %s++){ %s%s=new Array(",
				vArray, vArray, tr.src.lit[iArray], vArray, tr.src.name, tr.src.printArray()))
			tr.getValue(typ.Len)
			tr.src.WriteString(")")
		}

		if _, ok := typ.Elt.(*ast.ArrayType); ok {
			tr.getValue(typ.Elt)
		} else if len(tr.src.lit) > 1 {
			tr.src.WriteString("; " + strings.Repeat("}", len(tr.src.lit)-1))
		}

	// http://golang.org/pkg/go/ast/#BasicLit || godoc go/ast BasicLit
	//  Kind     token.Token // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING
	//  Value    string      // literal string
	case *ast.BasicLit:
		tr.src.WriteString(typ.Value)

		// === Add the value
		sign := ""

		if tr.src.isNegative {
			sign = "-"
		}
		tr.src.lit = append(tr.src.lit, sign+typ.Value)

	// http://golang.org/pkg/go/ast/#BinaryExpr || godoc go/ast BinaryExpr
	//  X     Expr        // left operand
	//  Op    token.Token // operator
	//  Y     Expr        // right operand
	case *ast.BinaryExpr:
		tr.getValue(typ.X)
		tr.src.WriteString(" " + typ.Op.String() + " ")
		tr.getValue(typ.Y)

	// http://golang.org/pkg/go/ast/#CallExpr || godoc go/ast CallExpr
	//  Fun      Expr      // function expression
	//  Args     []Expr    // function arguments; or nil
	case *ast.CallExpr:
		callIdent := typ.Fun.(*ast.Ident).Name

		// Conversion: []byte()
		if t, ok := typ.Fun.(*ast.ArrayType); ok {
			if t.Elt.(*ast.Ident).Name == "byte" {
				tr.getValue(typ.Args[0])
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
				tr.src.WriteString("new Array(")
				tr.getValue(typ.Args[len(typ.Args)-1]) // capacity
				tr.src.WriteString(")")

			// The second argument (in Args), if any, is the capacity which
			// is not useful in JS since it is dynamic.
			case *ast.MapType:
				tr.src.WriteString("{};") // or "new Object()"
			}

		case "new":
			switch argType := typ.Args[0].(type) {
			default:
				panic(fmt.Sprintf("[getValue] call of 'new' unimplemented: %T", argType))

			case *ast.ArrayType:
				for _, arg := range typ.Args {
					tr.getValue(arg)
				}
			}

		// Conversion
		case "uint", "uint8", "uint16", "uint32",
			"int", "int8", "int16", "int32",
			"float32", "float64", "byte", "rune", "string":
			tr.getValue(typ.Args[0])
		}

	// http://golang.org/pkg/go/ast/#CompositeLit || godoc go/ast CompositeLit
	//  Type   Expr      // literal type; or nil
	//  Elts   []Expr    // list of composite elements; or nil
	case *ast.CompositeLit:
		switch compoType := typ.Type.(type) {
		default:
			panic(fmt.Sprintf("[getValue] 'CompositeLit' unimplemented: %s", compoType))

		case *ast.ArrayType:
			tr.src.eltsLen = len(typ.Elts) // for ellipsis
			tr.getValue(typ.Type)
			//tr.src.pos = 

			// For arrays initialized
			if len(typ.Elts) != 0 {
				if compoType.Len == nil {
					tr.src.WriteString("[")
				} else {
					tr.src.WriteString(fmt.Sprintf("; %s = [", tr.src.name))
				}

				for i, el := range typ.Elts {
					if i != 0 {
						tr.src.WriteString(",")
					}
					tr.getValue(el)
				}
				tr.src.WriteString("]")
			}

		// http://golang.org/pkg/go/ast/#MapType || godoc go/ast MapType
		//  Key   Expr
		//  Value Expr
		case *ast.MapType:
			lenElts := len(typ.Elts) - 1
			tr.src.WriteString("{")

			for i, el := range typ.Elts {
				tr.getValue(el)

				if i != lenElts {
					tr.src.WriteString(", ")
				}
			}
			tr.src.WriteString("};")
		}

	// http://golang.org/pkg/go/ast/#Ellipsis || godoc go/ast Ellipsis
	//  Elt      Expr      // ellipsis element type (parameter lists only); or nil
	//case *ast.Ellipsis:

	// http://golang.org/pkg/go/ast/#Ident || godoc go/ast Ident
	//  Name    string    // identifier name
	case *ast.Ident:
		if typ.Name == "iota" {
			tr.src.WriteString("%d")
			tr.src.useIota = true
			break
		}
		// Undefined value in array / slice
		if len(tr.src.lit) != 0 && typ.Name == "_" {
			break
		}

		tr.src.WriteString(typ.Name)

	// http://golang.org/pkg/go/ast/#KeyValueExpr || godoc go/ast KeyValueExpr
	//  Key   Expr
	//  Value Expr
	case *ast.KeyValueExpr:
		tr.getValue(typ.Key)
		tr.src.WriteString(":")
		tr.getValue(typ.Value)

	// http://golang.org/pkg/go/ast/#ParenExpr || godoc go/ast ParenExpr
	//  X      Expr      // parenthesized expression
	case *ast.ParenExpr:
		tr.getValue(typ.X)

	// http://golang.org/pkg/go/ast/#StructType || godoc go/ast StructType
	//  Struct     token.Pos  // position of "struct" keyword
	//  Fields     *FieldList // list of field declarations
	//  Incomplete bool       // true if (source) fields are missing in the Fields list
	/*case *ast.StructType:*/

	// http://golang.org/pkg/go/ast/#StarExpr || godoc go/ast StarExpr
	//  X    Expr      // operand
	case *ast.StarExpr:
		tr.getValue(typ.X)

	// http://golang.org/pkg/go/ast/#UnaryExpr || godoc go/ast UnaryExpr
	//  Op    token.Token // operator
	//  X     Expr        // operand
	case *ast.UnaryExpr:
		switch typ.Op {
		case token.SUB:
			tr.src.isNegative = true
		}

		tr.src.WriteString(typ.Op.String())
		tr.getValue(typ.X)

	default:
		panic(fmt.Sprintf("[getValue] unimplemented: %T, value: %v",
			expr, expr))
	}
}
