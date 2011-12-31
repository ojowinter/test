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
	"go/ast"
)

// Gets the parameters.
//
// http://golang.org/pkg/go/ast/#FuncType || godoc go/ast FuncType
//  Func    token.Pos  // position of "func" keyword
//  Params  *FieldList // (incoming) parameters; or nil
//  Results *FieldList // (outgoing) results; or nil
func joinParams(f *ast.FuncType) string {
	isFirst := true
	s := ""

	for _, list := range f.Params.List {
		for _, id := range list.Names {
			if !isFirst {
				s += "," + SP
			}

			s += id.Name
			if isFirst {
				isFirst = false
			}
		}
	}

	return s
}

// Returns the initialization value in Go.
func initValue(val *ast.ValueSpec) string {
	var ident *ast.Ident

	switch typ := val.Type.(type) {
	case *ast.Ident:
		ident = typ
	case *ast.StarExpr:
		ident = typ.X.(*ast.Ident)
	default:
		panic("another type of value")
	}

	switch ident.Name {
	case "bool":
		return "false"
	case "string":
		return EMPTY
	case "uint", "uint8", "uint16", "uint32", "uint64",
		"int", "int8", "int16", "int32", "int64",
		"float32", "float64",
		"byte", "rune", "uintptr":
		return "0"
	//case "complex64", "complex128":
		//return "(0+0i)"
	}
	panic("unreachable")
}
