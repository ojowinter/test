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
	var _names      []string
	var iValidNames []int // index of variables which are not in blank

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

	// Check if there is any variable to use
	for i, v := range _names {
		if v != BLANK {
			iValidNames = append(iValidNames, i)
		}
	}
	if len(iValidNames) == 0 {
		return
	}

	if tr.isSwitch {
		tr.varSwitch = _names[len(_names)-1]
	}

	// === Function
	if values != nil {
		if call, ok := values[0].(*ast.CallExpr); ok {

			// Function literal
			if _, ok := call.Fun.(*ast.SelectorExpr); ok {
				goto _noFunc
			}

			// Declaration of slice/array
			fun := call.Fun.(*ast.Ident).Name
			if fun == "make" || fun == "new" {
				goto _noFunc
			}

			// === Assign variable to the output of a function
			fun = tr.getExpression(call).String()

			if len(_names) == 1 {
				tr.WriteString(_names[0] + SP + sign + SP + fun + ";")
				return
			}

			// multiple variables
			str := fmt.Sprintf("_%s", SP+sign+SP+fun)

			for _, i := range iValidNames {
				str += fmt.Sprintf(",%s_[%d]", SP+_names[i]+SP+sign+SP, i)
			}

			tr.WriteString(str + ";")
			return
		}
	}

_noFunc:

	for _, i := range iValidNames {
		name := _names[i]

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

//
// === Functions

// Writes the function declaration.
func (tr *transform) writeFunc(name *ast.Ident, typ *ast.FuncType) {
	if name != nil {
		tr.WriteString(fmt.Sprintf("function %s(%s)%s", name, joinParams(typ), SP))
	} else { // Literal function
		tr.WriteString(fmt.Sprintf("function(%s)%s", joinParams(typ), SP))
	}

	// Return multiple values
	declResults, declReturn := joinResults(typ)

	if declResults != "" {
		tr.WriteString("{" + SP + declResults)
		tr.skipLbrace = true
		tr.results = declReturn
	} else {
		tr.results = ""
	}
}

// http://golang.org/pkg/go/ast/#FuncType || godoc go/ast FuncType
//  Func    token.Pos  // position of "func" keyword
//  Params  *FieldList // (incoming) parameters; or nil
//  Results *FieldList // (outgoing) results; or nil

// http://golang.org/pkg/go/ast/#FieldList || godoc go/ast FieldList
//  Opening token.Pos // position of opening parenthesis/brace, if any
//  List    []*Field  // field list; or nil
//  Closing token.Pos // position of closing parenthesis/brace, if any

// http://golang.org/pkg/go/ast/#Field || godoc go/ast Field
//  Doc     *CommentGroup // associated documentation; or nil
//  Names   []*Ident      // field/method/parameter names; or nil if anonymous field
//  Type    Expr          // field/method/parameter type
//  Tag     *BasicLit     // field tag; or nil
//  Comment *CommentGroup // line comments; or nil

// Gets the parameters.
func joinParams(f *ast.FuncType) string {
	isFirst := true
	s := ""

	//if f.Params == nil {
		//return s
	//}

	for _, list := range f.Params.List {
		for _, v := range list.Names {
			if !isFirst {
				s += "," + SP
			}
			s += v.Name

			if isFirst {
				isFirst = false
			}
		}
	}

	return s
}

// Gets the results to use both in the declaration and in its return.
func joinResults(f *ast.FuncType) (decl, ret string) {
	isFirst := true
	isMultiple := false

	if f.Results == nil {
		return
	}

	for _, list := range f.Results.List {
		if list.Names == nil {
			continue
		}

		init := initValue(list.Type, "")

		for _, v := range list.Names {
			if !isFirst {
				decl += "," + SP
				ret += "," + SP
				isMultiple = true
			} else {
				isFirst = false
			}

			decl += fmt.Sprintf("%s=%s", v.Name+SP, SP+init)
			ret += v.Name
		}
	}

	if decl != "" {
		decl = "var " + decl + ";"
	}

	if isMultiple {
		ret = "[" + ret + "]"
	}
	ret = "return " + ret + ";"

	return
}
