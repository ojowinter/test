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
)

// Writes names and values for both declarations and assignments.
func (tr *transform) writeValues(names interface{}, values []ast.Expr,
type_ interface{}, sign string, isGlobal bool) {
	var _names []string
	isFirst := true
	skipSemicolon := false

	switch t := names.(type) {
	case []*ast.Ident:
		_names = make([]string, len(t))

		for i, v := range t {
			_names[i] = tr.getExpression(v).String()
		}
	}

	// Call (function)
	if values != nil {
		if call, ok := values[0].(*ast.CallExpr); ok {
			funcName := call.Fun.(*ast.Ident).Name
			if funcName == "make" || funcName == "new" {
				goto _noFunc
			}

			str := fmt.Sprintf("var _%s;", SP+sign+SP+tr.getExpression(call).String())

			for i, name := range _names {
				if name == BLANK {
					continue
				}

				if !isFirst {
					str += ","
				} else {
					isFirst = false
				}

				str += fmt.Sprintf("%s_[%d]", SP+name+SP+sign+SP, i)
			}

			if !isFirst {
				tr.WriteString(str + ";")
			}

			return
		}
	}

_noFunc:

	// number of names = number of values
	for i, name := range _names {
		if name == BLANK {
			continue
		}

		// === Name
		if isFirst {
			tr.WriteString("var " + name)
			isFirst = false
		} else {
			tr.WriteString("," + SP + name)
		}
		tr.WriteString(SP + sign + SP)

		if isGlobal {
			tr.addIfExported(name)
		}

		// === Value
		if values != nil {
			_value := values[i]

			// If the expression is an anonymous function, then
			// it is written in the main buffer.
			expr := tr.newExpression(name)
			expr.transform(values[i])

			/*if expr.hasError {
				return
			}*/

			if _, ok := _value.(*ast.FuncLit); !ok {
				tr.WriteString(initValue(type_, expr.String()))
			} else {
//				wasFunc = true
			}

			if expr.skipSemicolon {
				skipSemicolon = true
			}

		} else { // Initialization explicit
			tr.WriteString(initValue(type_, ""))
		}

	}

	if !isFirst && !skipSemicolon {
		tr.WriteString(";")
	}
}

// Returns the value, which is initialized if were necessary.
// A pointer is formatted like an array.
func initValue(type_ interface{}, value string) string {
	var ident *ast.Ident
	var isPointer bool

	switch typ := type_.(type) {
	case nil:
		return value
	case *ast.Ident:
		ident = typ
	case *ast.StarExpr:
		ident = typ.X.(*ast.Ident)
		isPointer = true
	default:
		panic(fmt.Sprintf("unexpected type of value: %T", typ))
	}

	if value == "" {
		switch ident.Name {
		case "bool":
			value = "false"
		case "string":
			value = EMPTY
		case "uint", "uint8", "uint16", "uint32", "uint64",
			"int", "int8", "int16", "int32", "int64",
			"float32", "float64",
			"byte", "rune", "uintptr":
			value = "0"
		//case "complex64", "complex128":
			//value = "(0+0i)"
		}
	}

	if isPointer {
		return "[" + value + "]"
	}
	return value
}
