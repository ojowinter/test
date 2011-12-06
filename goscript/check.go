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
	"os"
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
	isNegative bool
	isImplicit bool // implicit type?

	type_ string
}

// Initializes a new type of "check".
func newCheck(expr ast.Expr) (*check, error) {
	chk := new(check)

	switch expr.(type) {
	case nil:
		chk.isImplicit = true
	}

	return chk, chk.Type(expr)
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

	// An integer type is "int", by default
	if c.isImplicit && typ.Kind == token.INT {
		fmt.Fprintf(os.Stderr, "warning: %d: implicit integer type\n", typ.Pos())

		/*if !c.isCallExpr {
			if err := c.maxInt(typ); err != nil {
				return err
			}
		}*/
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
		ident := typ.Fun.(*ast.Ident).Name

		switch ident {
		case "make", "new":
			return c.Type(typ.Args[0])

		case "int", "uint", "int64", "uint64":
			return fmt.Errorf("%d: conversion of type %s", typ.Pos(), ident)
		}

	// http://golang.org/pkg/go/ast/#ChanType || godoc go/ast ChanType
	case *ast.ChanType:
		return fmt.Errorf("%d: channel type", typ.Pos())

	case *ast.CompositeLit:
		return c.Type(typ.Type)

	case *ast.Ident:
		switch typ.Name {
		// Unsupported types
		case "complex64", "complex128",
		"int64", "uint64", "int", "uint": // "uintptr"
			return fmt.Errorf("%d: %s type", typ.Pos(), typ.Name)
		}

	case *ast.InterfaceType: // ToDo: review

	case *ast.MapType:
		if err := c.Type(typ.Key); err != nil {
			return err
		}
		if err := c.Type(typ.Value); err != nil {
			return err
		}

	case *ast.ParenExpr:
		return c.Type(typ.X)

	// http://golang.org/pkg/go/ast/#StarExpr || godoc go/ast StarExpr
	//  X    Expr      // operand
	case *ast.StarExpr:
		return c.Type(typ.X)

	case *ast.StructType:
		

	case *ast.UnaryExpr:
		// Channel
		if typ.Op == token.ARROW {
			return fmt.Errorf("%d: channel operator", typ.Pos())
		}

		return c.Type(typ.X)

	// The type has not been indicated
	case nil:
		
	}

	return nil
}
