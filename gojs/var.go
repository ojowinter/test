// Copyright 2011  The "GoScript" Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gojs

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

// Constants
//
// http://golang.org/doc/go_spec.html#Constant_declarations
// https://developer.mozilla.org/en/JavaScript/Reference/Statements/const
func (tr *transform) getConst(spec []ast.Spec, isGlobal bool) {
	isMultipleLine := false
	iotaExpr := make([]string, 0) // iota expressions

	if len(spec) > 1 {
		isMultipleLine = true
	}

	// godoc go/ast ValueSpec
	//  Doc     *CommentGroup // associated documentation; or nil
	//  Names   []*Ident      // value names (len(Names) > 0)
	//  Type    Expr          // value type; or nil
	//  Values  []Expr        // initial values; or nil
	//  Comment *CommentGroup // line comments; or nil
	for _, s := range spec {
		vSpec := s.(*ast.ValueSpec)

		// Type checking
		if tr.getExpression(vSpec.Type).hasError {
			continue
		}

		tr.addLine(vSpec.Pos())
		isFirst := true

		for i, ident := range vSpec.Names {
			if ident.Name == "_" {
				iotaExpr = append(iotaExpr, "")
				continue
			}

			value := strconv.Itoa(ident.Obj.Data.(int)) // possible value of iota

			if vSpec.Values != nil {
				v := vSpec.Values[i]

				expr := tr.getExpression(v)
				if expr.hasError {
					continue
				}

				if expr.useIota {
					exprStr := expr.String()
					value = strings.Replace(exprStr, IOTA, value, -1)
					iotaExpr = append(iotaExpr, exprStr)
				} else {
					value = expr.String()
				}
			} else {
				if tr.hasError {
					continue
				}
				value = strings.Replace(iotaExpr[i], IOTA, value, -1)
			}

			if isGlobal {
				tr.addIfExported(ident)
			}

			// === Write
			if isFirst {
				isFirst = false

				if !isGlobal && isMultipleLine {
					tr.WriteString(strings.Repeat(TAB, tr.tabLevel))
				}
				tr.WriteString(fmt.Sprintf("const %s=%s", ident.Name+SP, SP+value))
			} else {
				tr.WriteString(fmt.Sprintf(",%s=%s", SP+ident.Name+SP, SP+value))
			}
		}

		// It is possible that there is only a blank identifier
		if !isFirst {
			tr.WriteString(";")
		}
	}
}

// Variables
//
// http://golang.org/doc/go_spec.html#Variable_declarations
// https://developer.mozilla.org/en/JavaScript/Reference/Statements/var
// https://developer.mozilla.org/en/JavaScript/Reference/Statements/let
//
// TODO: use let for local variables
func (tr *transform) getVar(spec []ast.Spec, isGlobal bool) {
	isMultipleLine := false

	if len(spec) > 1 {
		isMultipleLine = true
	}

	// godoc go/ast ValueSpec
	for _, s := range spec {
		vSpec := s.(*ast.ValueSpec)

		// Type checking
		if tr.getExpression(vSpec.Type).hasError {
			continue
		}

		tr.addLine(vSpec.Pos())
		// Pass token.DEFINE to know that it is a new variable
		tr.writeVar(vSpec.Names, vSpec.Values, vSpec.Type, token.DEFINE,
			isGlobal, isMultipleLine)
	}
}

// Types
//
// http://golang.org/doc/go_spec.html#Type_declarations
func (tr *transform) getType(spec []ast.Spec, isGlobal bool) {
	// godoc go/ast TypeSpec
	//  Doc     *CommentGroup // associated documentation; or nil
	//  Name    *Ident        // type name
	//  Type    Expr          // *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
	//  Comment *CommentGroup // line comments; or nil
	for _, s := range spec {
		tSpec := s.(*ast.TypeSpec)

		// Type checking
		if tr.getExpression(tSpec.Type).hasError {
			continue
		}

		switch typ := tSpec.Type.(type) {
		// godoc go/ast Ident
		//  NamePos token.Pos // identifier position
		//  Name    string    // identifier name
		//  Obj     *Object   // denoted object; or nil
		case *ast.Ident:

		// godoc go/ast StructType
		//  Struct     token.Pos  // position of "struct" keyword
		//  Fields     *FieldList // list of field declarations
		//  Incomplete bool       // true if (source) fields are missing in the Fields list
		//
		// godoc go/ast FieldList
		//  Opening token.Pos // position of opening parenthesis/brace, if any
		//  List    []*Field  // field list; or nil
		//  Closing token.Pos // position of closing parenthesis/brace, if any
		case *ast.StructType:
			if typ.Incomplete {
				panic("list of fields incomplete ???")
			}

			var fieldNames, fieldLines, fieldsInit string
			//!anonField := make([]bool, 0) // anonymous field

			firstPos := tr.getLine(typ.Fields.Opening)
			posOldField := firstPos
			posNewField := 0
			isFirst := true

			// godoc go/ast Field
			//  Doc     *CommentGroup // associated documentation; or nil
			//  Names   []*Ident      // field/method/parameter names; or nil if anonymous field
			//  Type    Expr          // field/method/parameter type
			//  Tag     *BasicLit     // field tag; or nil
			//  Comment *CommentGroup // line comments; or nil
			for _, field := range typ.Fields.List {
				if _, ok := field.Type.(*ast.FuncType); ok {
					tr.addError("%s: function type in struct",
						tr.fset.Position(field.Pos()))
					continue
				}

				// Type checking
				if tr.getExpression(field.Type).hasError {
					continue
				}
				if field.Names == nil {
					tr.addError("%s: anonymous field in struct",
						tr.fset.Position(field.Pos()))
					continue
				}

				init, _ := tr.initValue(true, field.Type)

				for _, v := range field.Names {
					name := v.Name
					if name == "_" {
						continue
					}

					if !isFirst {
						fieldNames += "," + SP
						fieldsInit += "," + SP
					}
					fieldNames += name
					fieldsInit += init
					//!anonField = append(anonField, false)

					// === Printing of fields
					posNewField = tr.getLine(v.Pos())

					if posNewField != posOldField {
						fieldLines += strings.Repeat(NL, posNewField - posOldField)
						fieldLines += strings.Repeat(TAB, tr.tabLevel + 1)
					} else {
						fieldLines += SP
					}
					fieldLines += fmt.Sprintf("this.%s=%s;", name, name)
					posOldField = posNewField
					// ===

					if isFirst {
						isFirst = false
					}
				}
			}

			// The right brace
			posNewField = tr.getLine(typ.Fields.Closing)

			if posNewField != posOldField {
				fieldLines += strings.Repeat(NL, posNewField - posOldField)
				fieldLines += strings.Repeat(TAB, tr.tabLevel)
			} else {
				fieldLines += SP
			}

			// Empty structs
			if fieldLines == SP {
				fieldLines = ""
			}

			// Write
			tr.addLine(tSpec.Pos())
			tr.WriteString(fmt.Sprintf(
				"function %s(%s)%s{%s}", tSpec.Name, fieldNames, SP, fieldLines))

			// Store the name of new type with its values initialized
			tr.types[tr.funcId][tr.blockId][tSpec.Name.Name] = fieldsInit

			tr.line += posNewField - firstPos // update the global position

		default:
			panic(fmt.Sprintf("unimplemented: %T", typ))
		}

		if tr.hasError {
			continue
		}
		if isGlobal {
			tr.addIfExported(tSpec.Name)
		}
	}
}

//
// === Utility

// Writes variables for both declarations and assignments.
func (tr *transform) writeVar(names interface{}, values []ast.Expr, type_ interface{}, operator token.Token, isGlobal, isMultipleLine bool) {
	var sign string
	var isNewVar, isBitClear bool

	if !isGlobal && isMultipleLine {
		tr.WriteString(strings.Repeat(TAB, tr.tabLevel))
	}

	// === Operator
	switch operator {
	case token.DEFINE:
		isNewVar = true
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
	var _names        []string
	var iValidNames   []int // index of variables which are not in blank
	var nameIsPointer []bool

	switch t := names.(type) {
	case []*ast.Ident:
		_names = make([]string, len(t))
		nameIsPointer = make([]bool, len(t))

		for i, v := range t {
			expr := tr.getExpression(v)

			_names[i] = expr.String()
			nameIsPointer[i] = expr.isPointer
		}
	case []ast.Expr: // like avobe
		_names = make([]string, len(t))
		nameIsPointer = make([]bool, len(t))

		for i, v := range t {
			expr := tr.getExpression(v)

			_names[i] = expr.String()
			nameIsPointer[i] = expr.isPointer
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
			if len(iValidNames) == 1 {
				i := iValidNames[0]
				tr.WriteString(fmt.Sprintf("%s[%d];",
					_names[i] + SP + sign + SP + fun, i))
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
	expr := tr.newExpression(nil)
	typeIsPointer := false
	isFirst := true

	for _, i := range iValidNames {
		name := _names[i]
		value := ""

		if isGlobal {
			tr.addIfExported(name)
		}
		tr.lastVarName = name

		// === Name
		if !isNewVar {
			name += tagPointer('P', tr.funcId, tr.blockId, name)
		}

		if isFirst {
			tr.WriteString(name)
			isFirst = false
		} else {
			tr.WriteString("," + SP + name)
		}
		tr.WriteString(SP + sign + SP)

		// === Value
		if values != nil {
			valueOfValidName := values[i]

			// If the expression is an anonymous function, then
			// it is written in the main buffer.
			expr = tr.newExpression(name)
			expr.transform(valueOfValidName)

			if _, ok := valueOfValidName.(*ast.FuncLit); !ok {
				exprStr := expr.String()

				if isBitClear {
					exprStr = "~(" + exprStr + ")"
				}
				value = exprStr

				_, typeIsPointer = tr.initValue(false, type_)
			}

		} else { // Initialization explicit
			value, typeIsPointer = tr.initValue(true, type_)
		}

		if isNewVar {
			tr.vars[tr.funcId][tr.blockId][name] = typeIsPointer

			// Could be addressed ahead
			if !expr.isPointer && !expr.isAddress && !typeIsPointer {
				value = tagPointer('L', tr.funcId, tr.blockId, name) +
					value +
					tagPointer('R', tr.funcId, tr.blockId, name)
			}
		}

		tr.WriteString(value)
	}

	if !isFirst && !expr.skipSemicolon && !tr.skipSemicolon {
		tr.WriteString(";")
	}
	if tr.skipSemicolon {
		tr.skipSemicolon = false
	}
}

// Returns the value initialized to zero according to its type if "init", and a
// boolean indicating if it is a pointer.
func (tr *transform) initValue(init bool, typ interface{}) (value string, typeIsPointer bool) {
	var ident *ast.Ident

	switch t := typ.(type) {
	case nil:
		return "", false
	case *ast.ArrayType:
		tr.skipSemicolon = true
		return tr.getExpression(t).String(), false

	case *ast.Ident:
		ident = t
	case *ast.StarExpr:
		tr.initIsPointer = true
		return tr.initValue(init, t.X)
	default:
		panic(fmt.Sprintf("initValue(): unexpected type: %T", typ))
	}

	if !init {
		return "", tr.initIsPointer
	}

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
	case "complex64", "complex128":
		value = "(0+0i)"
	default:
		value = ident.Name
		value = fmt.Sprintf("new %s(%s)", value, tr.getInitType(value))
	}

	if tr.initIsPointer {
		value = "[" + value + "]"
		typeIsPointer = true
		tr.initIsPointer = false
	}
	return
}

// Returns the fields of a custom type.
func (tr *transform) getTypeFields(fields []string) (args, allFields string) {
	for i, f := range fields {
		if i == 0 {
			args = f
		} else {
			args += "," + SP + f
			allFields += SP
		}

		allFields += fmt.Sprintf("this.%s=%s;", f, f)
	}
	return
}

// Returns the initializations of a custom type.
func (tr *transform) getInitType(name string) string {
	// In the actual function
	if tr.funcId != 0 {
		for block := tr.blockId; block >= 1; block-- {
			if _, ok := tr.types[tr.funcId][block][name]; ok {
				return tr.types[tr.funcId][block][name]
			}
		}
	}

	// Finally, search in the global variables (funcId = 0).
	for block := tr.blockId; block >= 0; block-- { // block until 0
		if _, ok := tr.types[0][block][name]; ok {
			return tr.types[0][block][name]
		}
	}
	//fmt.Printf("Function %d, block %d, name %s\n", tr.funcId, tr.blockId, name)
	panic("getInitType: type not found: " + name)
}
