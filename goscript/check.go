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

/*var types = []string{
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
}*/

type check struct {
	isCallExpr, isCompositeLit bool
	isInt64, isUint64 bool
	isNegative bool
	isImplicit bool // type not indicated
}

// Initializes a new type of "check".
func newCheck() *check {
	return &check{
		false, false,
		false, false,
		false,
		false,
	}
}

// Checks if it has a valid type for JavaScript.
func (c *check) Type(expr ast.Expr) error {
	switch typ := expr.(type) {
	default:
		panic(fmt.Sprintf("[Type] unimplemented: %T", typ))

	case *ast.ArrayType:
		return c.Type(typ.Elt)

	case *ast.BasicLit:
	// Check after calculating the mathematical expressions. ToDo

	// Integer checking
	if typ.Kind == token.INT {
		// An integer type is "int", by default
		if c.isImplicit && !c.isInt64 && !c.isUint64 {
			c.isInt64 = true
		}

		if !c.isCallExpr {
			/*if err := c.checkInt(typ); err != nil {
				return err
			}*/
		}
	}

	case *ast.BinaryExpr:
		if err := c.Type(typ.X); err != nil {
			return err
		}
		if err := c.Type(typ.Y); err != nil {
			return err
		}

	case *ast.CallExpr:
		c.isCallExpr = true

		switch typ.Fun.(*ast.Ident).Name {
		case "make", "new":
			return c.Type(typ.Args[0])
		}

	// http://golang.org/pkg/go/ast/#ChanType || godoc go/ast ChanType
	case *ast.ChanType:
		return fmt.Errorf("%d: channel type", typ.Pos())

	case *ast.CompositeLit:
		return c.Type(typ.Type)

	case *ast.Ident:
		switch typ.Name {
		// Unsupported types
		case "complex64", "complex128": // "uintptr"
			return fmt.Errorf("%d: %s type", typ.Pos(), typ.Name)

		// To check if the number fills into a JS number
		case "int", "int64":
			c.isInt64 = true
		case "uint", "uint64":
			c.isUint64 = true
		}

	case *ast.InterfaceType: // ToDo: review

	case *ast.MapType:
		if err := c.Type(typ.Key); err != nil {
			return err
		}
		if err := c.Type(typ.Value); err != nil {
			return err
		}

	// http://golang.org/pkg/go/ast/#StarExpr || godoc go/ast StarExpr
	//  X    Expr      // operand
	case *ast.StarExpr:
		return c.Type(typ.X)

	case *ast.UnaryExpr:
		// Channel
		if typ.Op == token.ARROW {
			return fmt.Errorf("%d: channel operator", typ.Pos())
		}

		return c.Type(typ.X)

	// The type has not been indicated
	case nil:
		c.isImplicit = true
	}

	return nil
}

// Checks the maximum size of an integer for JavaScript.
func (c *check) maxInt(number *ast.BasicLit) error {
	var errConv, errMax bool

	if c.isInt64 {
		n, err := strconv.Atoi64(number.Value)
		if err != nil {
			errConv = true
		}

		if n > MAX_INT_JS {
			errMax = true
		}
	}
	if c.isUint64 {
		n, err := strconv.Atoui64(number.Value)
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

	// === To print
	intString := "integer"
	if c.isUint64 {
		intString = "unsigned " + intString
	}

	num := number.Value
	if c.isNegative {
		num = "-" + num
	}
	// ===

	if errMax {
		return fmt.Errorf("%d: %s does not safely fit in a JS number",
			number.Pos(), num)
	}
	// if errConv {
	return fmt.Errorf("%d: %s could not be converted to %s",
		number.Pos(), num, intString)
}
