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
	"go/ast"
	"go/token"
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

// Checks if a literal is a type.
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

//
// === Get

// Gets the identifiers.
//
// http://golang.org/pkg/go/ast/#Ident || godoc go/ast Ident
//  Name    string    // identifier name
func getName(spec *ast.ValueSpec) (names []string, skipName []bool) {
	skipName = make([]bool, len(spec.Names)) // for blank identifiers "_"

	for i, v := range spec.Names {
		if v.Name == "_" {
			skipName[i] = true
			continue
		}
		names = append(names, v.Name)
	}

	return
}
