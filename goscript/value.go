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

// Gets the value.
// It throws a panic message for types no added.
func (v *value) getValue(iface interface{}) {
	// type Expr
	switch typ := iface.(type) {

	// http://golang.org/pkg/go/ast/#ArrayType || godoc go/ast ArrayType
	//  Len    Expr      // Ellipsis node for [...]T array types, nil for slice types
	//  Elt    Expr      // element type
	case *ast.ArrayType:
		if typ.Len == nil { // slice
			break
		}

		if len(v.lit) == 0 {
			v.WriteString("new Array(")

			if v.eltsLen != 0 { // ellipsis
				v.WriteString(strconv.Itoa(v.eltsLen))
			} else {
				v.getValue(typ.Len)
			}

			v.WriteString(")")
		} else {
			iArray := len(v.lit) - 1             // index of array
			vArray := "i" + strconv.Itoa(iArray) // variable's name for the loop

			v.WriteString(fmt.Sprintf("; for (var %s=0; %s<%s; %s++){ %s%s=new Array(",
				vArray, vArray, v.lit[iArray], vArray, v.name, v.printArray()))
			v.getValue(typ.Len)
			v.WriteString(")")
		}

		if _, ok := typ.Elt.(*ast.ArrayType); ok {
			v.getValue(typ.Elt)
		} else if len(v.lit) > 1 {
			v.WriteString("; " + strings.Repeat("}", len(v.lit)-1))
		}

	// http://golang.org/pkg/go/ast/#BasicLit || godoc go/ast BasicLit
	//  Kind     token.Token // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING
	//  Value    string      // literal string
	case *ast.BasicLit:
		v.WriteString(typ.Value)

		// === Add the value
		sign := ""

		if v.isNegative {
			sign = "-"
		}
		v.lit = append(v.lit, sign+typ.Value)

	// http://golang.org/pkg/go/ast/#BinaryExpr || godoc go/ast BinaryExpr
	//  X     Expr        // left operand
	//  Op    token.Token // operator
	//  Y     Expr        // right operand
	case *ast.BinaryExpr:
		v.getValue(typ.X)
		v.WriteString(" " + typ.Op.String() + " ")
		v.getValue(typ.Y)

	// http://golang.org/pkg/go/ast/#CallExpr || godoc go/ast CallExpr
	//  Fun      Expr      // function expression
	//  Args     []Expr    // function arguments; or nil
	case *ast.CallExpr:
		callIdent := typ.Fun.(*ast.Ident).Name

		// Conversion: []byte()
		if t, ok := typ.Fun.(*ast.ArrayType); ok {
			if t.Elt.(*ast.Ident).Name == "byte" {
				v.getValue(typ.Args[0])
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
				v.WriteString("new Array(")
				v.getValue(typ.Args[len(typ.Args)-1]) // capacity
				v.WriteString(")")

			// The second argument (in Args), if any, is the capacity which
			// is not useful in JS since it is dynamic.
			case *ast.MapType:
				v.WriteString("{};") // or "new Object()"
			}

		case "new":
			switch argType := typ.Args[0].(type) {
			default:
				panic(fmt.Sprintf("[getValue] call of 'new' unimplemented: %T", argType))

			case *ast.ArrayType:
				for _, arg := range typ.Args {
					v.getValue(arg)
				}
			}

		// Conversion
		case "uint8", "uint16", "uint32",
			"int8", "int16", "int32",
			"float32", "byte", "rune", "string":
			v.getValue(typ.Args[0])
		}

	// http://golang.org/pkg/go/ast/#CompositeLit || godoc go/ast CompositeLit
	//  Type   Expr      // literal type; or nil
	//  Elts   []Expr    // list of composite elements; or nil
	case *ast.CompositeLit:
		switch compoType := typ.Type.(type) {
		default:
			panic(fmt.Sprintf("[getValue] 'CompositeLit' unimplemented: %s", compoType))

		case *ast.ArrayType:
			v.eltsLen = len(typ.Elts) // for ellipsis
			v.getValue(typ.Type)

			// For arrays initialized
			if len(typ.Elts) != 0 {
				if compoType.Len == nil {
					v.WriteString("[")
				} else {
					v.WriteString(fmt.Sprintf("; %s = [", v.name))
				}

				for i, el := range typ.Elts {
					if i != 0 {
						v.WriteString(",")
					}
					v.getValue(el)
				}
				v.WriteString("]")
			}

		// http://golang.org/pkg/go/ast/#MapType || godoc go/ast MapType
		//  Key   Expr
		//  Value Expr
		case *ast.MapType:
			lenElts := len(typ.Elts) - 1

			v.WriteString("{")

			for i, el := range typ.Elts {
				v.getValue(el)

				if i != lenElts {
					v.WriteString(",")
				}
			}
			v.WriteString("\n};")
		}

	// http://golang.org/pkg/go/ast/#Ellipsis || godoc go/ast Ellipsis
	//  Elt      Expr      // ellipsis element type (parameter lists only); or nil
	//case *ast.Ellipsis:

	// http://golang.org/pkg/go/ast/#Ident || godoc go/ast Ident
	//  Name    string    // identifier name
	case *ast.Ident:
		if typ.Name == "iota" {
			v.WriteString("%d")
			v.useIota = true
			break
		}
		// Undefined value in array / slice
		if len(v.lit) != 0 && typ.Name == "_" {
			break
		}

		v.WriteString(typ.Name)

	// http://golang.org/pkg/go/ast/#KeyValueExpr || godoc go/ast KeyValueExpr
	//  Key   Expr
	//  Value Expr
	case *ast.KeyValueExpr:
		v.WriteString("\n\t")
		v.getValue(typ.Key)
		v.WriteString(": ")
		v.getValue(typ.Value)

	// http://golang.org/pkg/go/ast/#ParenExpr || godoc go/ast ParenExpr
	//  X      Expr      // parenthesized expression
	case *ast.ParenExpr:
		v.getValue(typ.X)

	// http://golang.org/pkg/go/ast/#StructType || godoc go/ast StructType
	//  Struct     token.Pos  // position of "struct" keyword
	//  Fields     *FieldList // list of field declarations
	//  Incomplete bool       // true if (source) fields are missing in the Fields list
	/*case *ast.StructType:*/

	// http://golang.org/pkg/go/ast/#StarExpr || godoc go/ast StarExpr
	//  X    Expr      // operand
	case *ast.StarExpr:
		v.getValue(typ.X)

	// http://golang.org/pkg/go/ast/#UnaryExpr || godoc go/ast UnaryExpr
	//  Op    token.Token // operator
	//  X     Expr        // operand
	case *ast.UnaryExpr:
		switch typ.Op {
		case token.SUB:
			v.isNegative = true
		}

		v.WriteString(typ.Op.String())
		v.getValue(typ.X)

	default:
		panic(fmt.Sprintf("[getValue] unimplemented: %T, value: %v",
			iface, iface))
	}
}
