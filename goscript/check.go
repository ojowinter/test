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

// Gets both type name and position.
// If returns an empty literal, then it has not been indicated.
func getType(expr ...ast.Expr) (name []string, pos []token.Pos, err error) {
	for _, e := range expr {
		switch typ := e.(type) {
		default:
			panic(fmt.Sprintf("[checkType] unimplemented: %T", typ))

		// The type has not been indicated
		case nil:
			//return "", 0//, typ.Pos()

		// 
		case *ast.BasicLit:
			//return "", 0//

		// Get the type data
		case *ast.Ident:
			name = append(name, typ.Name)
			pos = append(pos, typ.Pos())

		// * * *

		case *ast.ArrayType:
			return getType(typ.Elt)

		case *ast.BinaryExpr:
			return getType(typ.X, typ.Y)

		case *ast.CallExpr:
			switch typ.Fun.(*ast.Ident).Name {
			case "make", "new":
				return getType(typ.Args[0])
			}

		case *ast.ChanType:
			err = fmt.Errorf("%d: channel type", typ.Pos())

		case *ast.CompositeLit:
			return getType(typ.Type)

		case *ast.InterfaceType: // ToDo: review

		case *ast.MapType:
			return getType(typ.Key, typ.Value)

		// http://golang.org/pkg/go/ast/#StarExpr || godoc go/ast StarExpr
		//  X    Expr      // operand
		case *ast.StarExpr:
			return getType(typ.X)

		case *ast.UnaryExpr:
			// Channel
			if typ.Op == token.ARROW {
				err = fmt.Errorf("%d: channel operator", typ.Pos())
				break
			}
			return getType(typ.X)
		}
	}

	return
}

// * * *

// Checks if it has a valid type for JavaScript.
func checkType(expr ast.Expr) error {
	name, pos, err := getType(expr)
	if err != nil {
		return err
	}

	for i, v := range name {
		if v == "complex64" || v == "complex128" { // || v == "uintptr"
			return fmt.Errorf("%d: %s type", pos[i], v)
		}
	}
	return nil
}

// Checks the maximum size of an integer for JavaScript.
func checkInt(value, type_ string, isNegative bool, pos token.Pos) error {
	var errConv, errMax, isInt bool
	intString := "integer" // to print

//println(type_)

	// Check overflow
	if isNegative {
		switch type_ {
		case "uint", "uint8", "uint16", "uint32", "uint64":
			return fmt.Errorf("%d: -%s overflows %s", pos, value, type_)
		}
	}

	switch type_ {
	case "int", "int64", "": // an integer type is "int", by default
		isInt = true
	case "uint", "uint64":
		intString = "unsigned " + intString
	default:
		return nil
	}

	if isInt {
		n, err := strconv.Atoi64(value)
		if err != nil {
			errConv = true
		}

		if n > MAX_INT_JS {
			errMax = true
		}
	} else {
		n, err := strconv.Atoui64(value)
		if err != nil {
			errConv = true
		}

		if n > MAX_UINT_JS {
			errMax = true
		}
	}

	if !errConv && !errMax {
		return nil
	}

	// To print the sign
	if isNegative {
		value = "-" + value
	}

	if errMax {
		return fmt.Errorf("%d: %s does not safely fit in a JS number",
			pos, value)
	}
	// if errConv {
	return fmt.Errorf("%d: %s could not be converted to %s",
		pos, value, intString)
}
