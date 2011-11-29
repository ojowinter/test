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
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

// Maximum size for integers in JavaScript.
const (
	MAX_UINT_JS = 1<<53 - 1
	MAX_INT_JS  = 1<<52 - 1
)

var types = []string{
	"bool", "string",

	// Numeric types
	"uint8", "uint16", "uint32", "uint64",
	"int8", "int16", "int32", "int64",
	"float32", "float64",
	"complex64", "complex128",
	"byte", "rune", "uint", "int", "uintptr",
}

// JavaScript does not support numbers of 64 bits.
var invalidTypes = []string{
	"uint64", "int64", "float64", "complex64", "complex128", //"uintptr",
}

// Checks if the literal is a type.
func isType(tok token.Token, lit string) bool {
	if tok != token.IDENT {
		return false
	}

	for _, v := range types {
		if v == lit {
			return true
		}
	}
	return false
}

// Checks if the literal is a valid type for JavaScript.
func isValidType(lit string) bool {
	for _, v := range invalidTypes {
		if v == lit {
			return false
		}
	}
	return true
}

// * * *

// Checks if it is a valid type, when it is used an explicit type.
func checkType(expr ...ast.Expr) error {
	for _, v := range expr {
		switch typ := v.(type) {
		default:
			panic(fmt.Sprintf("[checkType] unimplemented: %T", typ))

		// The type has not been indicated
		case nil:

		// Elt    Expr      // element type
		case *ast.ArrayType:
			return checkType(typ.Elt)

		// Name    string    // identifier name
		case *ast.Ident:
			if ok := isValidType(typ.Name); !ok {
				return fmt.Errorf("Unsupported %q type: line %d",
					typ.Name, typ.Pos())
			}

		case *ast.InterfaceType:

		// X    Expr      // operand
		case *ast.StarExpr:
			return checkType(typ.X)
		}
	}

	return nil
}

// Checks the maximum size of an integer for JavaScript.
func checkInt(num string, isNegative bool, pos token.Pos) error {
	var errConv, errMax bool
	intString := "integer" // to print

	if isNegative {
		n, err := strconv.Atoi64(num)
		if err != nil {
			errConv = true
		}

		if n > MAX_INT_JS {
			errMax = true
		}

		num = "-" + num // to print the sign
	} else {
		intString = "unsigned " + intString

		n, err := strconv.Atoui64(num)
		if err != nil {
			errConv = true
		}

		if n > MAX_UINT_JS {
			errMax = true
		}
	}

	if errConv {
		return fmt.Errorf("%q could not be converted to %s: line %d",
			num, intString, pos)
	}
	if errMax {
		return fmt.Errorf("%q does not safely fit in a JavaScript number: line %d",
			num, pos)
	}
	return nil
}
