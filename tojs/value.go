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
	"go/token"
	"strconv"
	"strings"
)

// Represents a value.
type value struct {
	useIota bool
	ident   string   // variable's identifier
	lit     []string // store the last literals (for array)
	len     int      // store length of array, to use in case of ellipsis (...)
	*bytes.Buffer
}

// Initializes a new type of "value".
func newValue(identifier string) *value {
	return &value{false, identifier, make([]string, 0), 0, new(bytes.Buffer)}
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
func (v *value) getValue(iface interface{}) error {
	// type Expr
	switch typ := iface.(type) {

	// http://golang.org/pkg/go/ast/#ArrayType || godoc go/ast ArrayType
	//  Len    Expr      // Ellipsis node for [...]T array types, nil for slice types
	//  Elt    Expr      // element type
	case *ast.ArrayType:
		array := iface.(*ast.ArrayType)

		if array.Len == nil { // slice
			break
		}

		if len(v.lit) == 0 {
			v.WriteString("new Array(")

			if v.len != 0 { // ellipsis
				v.WriteString(strconv.Itoa(v.len))
			} else {
				v.getValue(array.Len)
			}

			v.WriteString(")")
		} else {
			iArray := len(v.lit) - 1             // index of array
			vArray := "i" + strconv.Itoa(iArray) // variable's name for the loop

			v.WriteString(fmt.Sprintf("; for (var %s=0; %s<%s; %s++){ %s%s=new Array(",
				vArray, vArray, v.lit[iArray], vArray, v.ident, v.printArray()))
			v.getValue(array.Len)
			v.WriteString(")")
		}

		if _, ok := array.Elt.(*ast.ArrayType); ok {
			v.getValue(array.Elt)
		} else if len(v.lit) > 1 {
			v.WriteString("; " + strings.Repeat("}", len(v.lit)-1))
		}

	// http://golang.org/pkg/go/ast/#BasicLit || godoc go/ast BasicLit
	//  Value    string      // literal string
	case *ast.BasicLit:
		lit := iface.(*ast.BasicLit).Value

		v.WriteString(lit)
		v.lit = append(v.lit, lit)

	// http://golang.org/pkg/go/ast/#BinaryExpr || godoc go/ast BinaryExpr
	//  X     Expr        // left operand
	//  Op    token.Token // operator
	//  Y     Expr        // right operand
	case *ast.BinaryExpr:
		binaryExpr := iface.(*ast.BinaryExpr)

		v.getValue(binaryExpr.X)
		v.WriteString(" " + binaryExpr.Op.String() + " ")
		v.getValue(binaryExpr.Y)

	// http://golang.org/pkg/go/ast/#CallExpr || godoc go/ast CallExpr
	//  Fun      Expr      // function expression
	//  Args     []Expr    // function arguments; or nil
	case *ast.CallExpr:
		call := iface.(*ast.CallExpr)
		callIdent := call.Fun.(*ast.Ident).Name

		switch callIdent {
		default:
			panic(fmt.Sprintf("[getValue] call unimplemented: %s", callIdent))

		case "make":
			switch typ := call.Args[0].(type) {
			default:
				panic(fmt.Sprintf("[getValue] call of 'make' unimplemented: %T", typ))

			// For slice
			case *ast.ArrayType:
				v.WriteString("new Array(")
				v.getValue(call.Args[len(call.Args)-1]) // capacity
				v.WriteString(")")

			// The second argument (in Args), if any, is the capacity which
			// is not useful in JS since it is dynamic.
			case *ast.MapType:
				v.WriteString("{};") // or "new Object()"

			case *ast.ChanType:
				return fmt.Errorf("channel type, at line: %v",
					call.Args[0].(*ast.ChanType).Pos())
			}

		case "new":
			switch typ := call.Args[0].(type) {
			default:
				panic(fmt.Sprintf("[getValue] call of 'new' unimplemented: %T", typ))

			// Check if it is an array
			case *ast.ArrayType:
				for _, arg := range call.Args {
					v.getValue(arg)
				}
			}
		}

	// http://golang.org/pkg/go/ast/#CompositeLit || godoc go/ast CompositeLit
	//  Type   Expr      // literal type; or nil
	//  Elts   []Expr    // list of composite elements; or nil
	case *ast.CompositeLit:
		composite := iface.(*ast.CompositeLit)

		switch typ := composite.Type.(type) {
		case *ast.ArrayType:
			v.len = len(composite.Elts) // for ellipsis
			v.getValue(composite.Type)

			// For arrays initialized
			if len(composite.Elts) != 0 {
				if typ.Len == nil {
					v.WriteString("[")
				} else {
					v.WriteString(fmt.Sprintf("; %s = [", v.ident))
				}

				for i, el := range composite.Elts {
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
			lenElts := len(composite.Elts) - 1
			v.WriteString("{")

			for i, el := range composite.Elts {
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
		keyValueExpr := iface.(*ast.KeyValueExpr)

		v.WriteString("\n\t")
		v.getValue(keyValueExpr.Key)
		v.WriteString(": ")
		v.getValue(keyValueExpr.Value)

	// http://golang.org/pkg/go/ast/#UnaryExpr || godoc go/ast UnaryExpr
	//  Op    token.Token // operator
	//  X     Expr        // operand
	case *ast.UnaryExpr:
		unaryExpr := iface.(*ast.UnaryExpr)

		if unaryExpr.Op == token.ARROW { // channel
			return fmt.Errorf("channel operator, at line: %d",
				unaryExpr.OpPos)
		}

		v.WriteString(unaryExpr.Op.String())
		v.getValue(unaryExpr.X)

	default:
		panic(fmt.Sprintf("[getValue] unimplemented: %T, value: %v",
			iface, iface))
	}

	return nil
}
