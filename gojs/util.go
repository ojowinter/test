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

// Writes names and values for both declarations and assignments.
func (tr *transform) writeValues(names interface{}, values []ast.Expr, type_ interface{}, operator token.Token, isGlobal bool) {
	var sign string
	var skipSemicolon, isBitClear bool
	isFirst := true

	// === Operator
	switch operator {
	case token.DEFINE:
		tr.WriteString("var ")
		sign = "="
	case token.ASSIGN,
		token.ADD_ASSIGN, token.SUB_ASSIGN, token.MUL_ASSIGN, token.QUO_ASSIGN,
		token.REM_ASSIGN,
		token.AND_ASSIGN, token.OR_ASSIGN, token.XOR_ASSIGN, token.SHL_ASSIGN,
		token.SHR_ASSIGN:

		sign = operator.String()
	case token.AND_NOT_ASSIGN:
		sign = "&="
		isBitClear = true

	default:
		panic(fmt.Sprintf("operator unimplemented: %s", operator.String()))
	}

	// === Names
	var _names []string

	switch t := names.(type) {
	case []*ast.Ident:
		_names = make([]string, len(t))

		for i, v := range t {
			_names[i] = tr.getExpression(v).String()
		}
	case []ast.Expr: // like avobe
		_names = make([]string, len(t))

		for i, v := range t {
			_names[i] = tr.getExpression(v).String()
		}
	default:
		panic("unreachable")
	}

	if tr.isSwitch {
		tr.switchInit = _names[len(_names)-1]
	}

	// === Function
	if values != nil {
		if call, ok := values[0].(*ast.CallExpr); ok {

			// Anonymous function
			if _, ok := call.Fun.(*ast.SelectorExpr); ok {
				goto _noFunc
			}

			// Declaration of slice/array
			fun := call.Fun.(*ast.Ident).Name
			if fun == "make" || fun == "new" {
				goto _noFunc
			}

			// Handle return of multiple values
			str := fmt.Sprintf("_%s;", SP+sign+SP+tr.getExpression(call).String())

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
			tr.WriteString(name)
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

			if _, ok := _value.(*ast.FuncLit); !ok {
				exprStr := expr.String()

				if isBitClear {
					exprStr = "~(" + exprStr + ")"
				}

				tr.WriteString(initValue(type_, exprStr))
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
