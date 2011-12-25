// Copyright 2011  The "GoJscript" Authors
//
// Use of this source code is governed by the BSD 2-Clause License
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

package gojs

import (
	"fmt"
	"go/ast"
	"go/token"
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

// Checks the expression and adds the error (if any).
func (tr *transform) CheckAndAddError(expr ast.Expr) error {
	err := tr.checkType(expr)
	if err != nil {
		tr.addError(err)
	}

	return err
}

// Checks if "expr" has a valid type for JavaScript.
func (tr *transform) checkType(expr ast.Expr) error {
	switch typ := expr.(type) {
	default:
		panic(fmt.Sprintf("unimplemented: %T", typ))

	case *ast.ArrayType:
		return tr.checkType(typ.Elt)

	case *ast.BasicLit:

	case *ast.BinaryExpr:
		if err := tr.checkType(typ.X); err != nil {
			return err
		}
		if err := tr.checkType(typ.Y); err != nil {
			return err
		}

	case *ast.CallExpr:
		call := typ.Fun.(*ast.Ident).Name

		switch call {
		case "make", "new":
			return tr.checkType(typ.Args[0])

		case "int64", "uint64":
			return fmt.Errorf("%s: conversion of type %s", tr.fset.Position(typ.Pos()), call)

		// golang.org/pkg/builtin/
		case "complex":
			return fmt.Errorf("%s: built-in function %s()", tr.fset.Position(typ.Pos()), call)
		}

	// http://golang.org/pkg/go/ast/#ChanType || godoc go/ast ChanType
	case *ast.ChanType:
		return fmt.Errorf("%s: channel type", tr.fset.Position(typ.Pos()))

	case *ast.CompositeLit:
		return tr.checkType(typ.Type)

	case *ast.Ident:
		switch typ.Name {
		// Unsupported types
		case "int64", "uint64", "complex64", "complex128": // "uintptr"
			return fmt.Errorf("%s: %s type", tr.fset.Position(typ.Pos()), typ.Name)
		}

	case *ast.InterfaceType: // ToDo: review

	case *ast.MapType:
		if err := tr.checkType(typ.Key); err != nil {
			return err
		}
		if err := tr.checkType(typ.Value); err != nil {
			return err
		}

	case *ast.ParenExpr:
		return tr.checkType(typ.X)

	// http://golang.org/pkg/go/ast/#StarExpr || godoc go/ast StarExpr
	//  X    Expr      // operand
	case *ast.StarExpr:
		return tr.checkType(typ.X)

	case *ast.StructType:

	case *ast.UnaryExpr:
		// Channel
		if typ.Op == token.ARROW {
			return fmt.Errorf("%s: channel operator", tr.fset.Position(typ.Pos()))
		}

		return tr.checkType(typ.X)

	// The type has not been indicated
	case nil:
	}

	return nil
}
