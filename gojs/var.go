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
	tr.isConst = true

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

	tr.isConst = false
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

		// It is necessary to add the first variable before of checking
		tr.lastVarName = vSpec.Names[0].Name

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
				isPointer := false

				if _, ok := field.Type.(*ast.FuncType); ok {
					tr.addError("%s: function type in struct",
						tr.fset.Position(field.Pos()))
					continue
				}
				if field.Names == nil {
					tr.addError("%s: anonymous field in struct",
						tr.fset.Position(field.Pos()))
					continue
				}
				// Type checking
				if expr := tr.getExpression(field.Type); expr.hasError {
					continue
				} else if expr.isPointer {
					isPointer = true
				}

				zero, _ := tr.zeroValue(true, field.Type)

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
					fieldsInit += zero
					//!anonField = append(anonField, false)

					// === Printing of fields
					posNewField = tr.getLine(v.Pos())

					if posNewField != posOldField {
						fieldLines += strings.Repeat(NL, posNewField - posOldField)
						fieldLines += strings.Repeat(TAB, tr.tabLevel + 1)
					} else {
						fieldLines += SP
					}

					if !isPointer {
						fieldLines += fmt.Sprintf("this.%s=%s;", name, name)
					} else {
						fieldLines += fmt.Sprintf("this.%s={p:%s};", name, name)
					}

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
/*			tr.WriteString(fmt.Sprintf("function %s(%s)%s{%sthis._z=%q;%s}",
				tSpec.Name, fieldNames, SP,
				SP, fieldsInit, fieldLines))
*/
			// Store the name of new type with its values initialized
			tr.typeZero[tr.funcId][tr.blockId][tSpec.Name.Name] = fieldsInit

			tr.line += posNewField - firstPos // update the global position

		default:
			tr.addLine(tSpec.Pos())
			tr.WriteString(fmt.Sprintf("function %s(t)%s{%sthis.t=t;%s}",
				tSpec.Name, SP, SP, SP))
		}

		if tr.hasError {
			continue
		}
		if isGlobal {
			tr.addIfExported(tSpec.Name)
		}
	}
}

// === Utility
//

// Writes variables for both declarations and assignments.
func (tr *transform) writeVar(names interface{}, values []ast.Expr, type_ interface{}, operator token.Token, isGlobal, isMultipleLine bool) {
	var sign string
	var isNewVar, isBitClear bool

	tr.isVar = true
	defer func() { tr.isVar = false }()

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
	var idxValidNames []int // index of variables which are not in blank
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

	// Check if there is any variable to use; and it is exported
	for i, v := range _names {
		if v != BLANK {
			idxValidNames = append(idxValidNames, i)

			if isGlobal {
				tr.addIfExported(v)
			}
		}
	}
	if len(idxValidNames) == 0 {
		return
	}

	if values != nil {
		// === Function
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
			if len(idxValidNames) == 1 {
				i := idxValidNames[0]
				tr.WriteString(fmt.Sprintf("%s[%d];",
					_names[i] + SP + sign + SP + fun, i))
				return
			}

			// multiple variables
			str := fmt.Sprintf("_%s", SP+sign+SP+fun)

			for _, i := range idxValidNames {
				str += fmt.Sprintf(",%s_[%d]", SP+_names[i]+SP+sign+SP, i)
			}

			tr.WriteString(str + ";")
			return
		}
	}

_noFunc:
	expr := tr.newExpression(nil)
	typeIs := otherType
	isFuncLit := false
	isZeroValue := false
	isFirst := true
	value := ""

	if values == nil { // initialization explicit
		value, typeIs = tr.zeroValue(true, type_)
		isZeroValue = true
	}

	for _, i := range idxValidNames {
		name := _names[i]
		nameExpr := ""

		tr.lastVarName = name

		// === Name
		if isFirst {
			nameExpr += name
			isFirst = false
		} else {
			 nameExpr += "," + SP + name
		}

		if !isNewVar {
			nameExpr += tagPointer(false, 'P', tr.funcId, tr.blockId, name)
		}

		// === Value
		if isZeroValue {
			if typeIs == sliceType {
				tr.slices[tr.funcId][tr.blockId][name] = void
			}
		} else {
			var valueOfValidName ast.Expr

			// _, ok = m[k]
			if len(values) == 1 && i == 1 {
				valueOfValidName = values[0]
			} else {
				valueOfValidName = values[i]
			}

			// If the expression is an anonymous function, then, at transforming,
			// it is written in the main buffer.
			expr = tr.newExpression(name)
			expr.isValue = true

			if _, ok := valueOfValidName.(*ast.FuncLit); !ok {
				expr.transform(valueOfValidName)
				exprStr := expr.String()

				if isBitClear {
					exprStr = "~(" + exprStr + ")"
				}
				value = exprStr

				_, typeIs = tr.zeroValue(false, type_)

				if expr.isAddress {
					tr.addr[tr.funcId][tr.blockId][name] = true
					if !isNewVar {
						nameExpr += ADDR
					}
				} /*else {
					tr.addr[tr.funcId][tr.blockId][name] = false
				}*/

				// == Map: v, ok := m[k]
				if len(values) == 1 && tr.isType(mapType, expr.mapName) {
					value = value[:len(value)-3] // remove '[0]'

					if len(idxValidNames) == 1 {
						tr.WriteString(fmt.Sprintf("%s%s%s[%d];",
							_names[idxValidNames[0]],
							SP + sign + SP,
							value, idxValidNames[0]))
					} else {
						tr.WriteString(fmt.Sprintf("_%s,%s_[%d],%s_[%d];",
							SP + sign + SP + value,
							SP + _names[0] + SP + sign + SP, 0,
							SP + _names[1] + SP + sign + SP, 1))
					}

					return
				}
				// ==
			} else {
				isFuncLit = true

				tr.WriteString(nameExpr)
				expr.transform(valueOfValidName)
			}

			// Check if new variables assigned to another ones are slices or maps.
			if isNewVar && expr.isIdent {
				if tr.isType(sliceType, value) {
					tr.slices[tr.funcId][tr.blockId][name] = void
				}
				if tr.isType(mapType, value) {
					tr.maps[tr.funcId][tr.blockId][name] = void
				}
			}
		}

		if isNewVar {
			typeIsPointer := false
			if typeIs == pointerType {
				typeIsPointer = true
			}

			tr.vars[tr.funcId][tr.blockId][name] = typeIsPointer

			// The value could be a pointer so this new variable has to be it.
			if tr.vars[tr.funcId][tr.blockId][value] {
				tr.vars[tr.funcId][tr.blockId][name] = true
			}

			// Could be addressed ahead
			if value != "" && !expr.isPointer && !expr.isAddress && !typeIsPointer {
				value = tagPointer(isZeroValue, 'L', tr.funcId, tr.blockId, name) +
					value +
					tagPointer(isZeroValue, 'R', tr.funcId, tr.blockId, name)
			}
		}

		if !isFuncLit {
			tr.WriteString(nameExpr)

			if expr.isSlice {
				if tr.isType(sliceType, expr.name) {
					tr.WriteString(".fromSlice(" + value + ")")
				} else {
					tr.WriteString(".fromArray(" + value + ")")
				}
			} else if value != "" {
				tr.WriteString(SP + sign + SP + value)
			}
		}
	}

	if !isFirst && !expr.skipSemicolon && !tr.skipSemicolon {
		tr.WriteString(";")
	}
	if tr.skipSemicolon {
		tr.skipSemicolon = false
	}
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

// === Zero value
//

type dataType uint8

const (
	otherType dataType = iota
	mapType
	pointerType
	sliceType
)

// Returns the zero value of the value type if "init", and a boolean indicating
// if it is a pointer.
func (tr *transform) zeroValue(init bool, typ interface{}) (value string, dt dataType) {
	var ident *ast.Ident

	switch t := typ.(type) {
	case nil:
		return

	case *ast.MapType:
		return "", mapType

	case *ast.ArrayType:
		if t.Len != nil {
			tr.skipSemicolon = true
			return tr.getExpression(t).String(), otherType
		}
		return fmt.Sprintf("new g.S(undefined,%s0,%s0)", SP, SP), sliceType

	case *ast.InterfaceType: // nil
		return "undefined", otherType

	case *ast.Ident:
		ident = t
	case *ast.StarExpr:
		tr.initIsPointer = true
		return tr.zeroValue(init, t.X)
	default:
		panic(fmt.Sprintf("zeroValue(): unexpected type: %T", typ))
	}

	if !init {
		if tr.initIsPointer {
			return "", pointerType
		}
		return
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
		value = fmt.Sprintf("new %s(%s)", value, tr.zeroOfType(value))
	}

	if tr.initIsPointer {
		value = "{p:undefined}"
		dt = pointerType
		tr.initIsPointer = false
	}
	return
}

// Returns the zero value of a map.
func (tr *transform) zeroOfMap(m *ast.MapType) string {
	if mapT, ok := m.Value.(*ast.MapType); ok { // nested map
		return tr.zeroOfMap(mapT)
	}
	v, _ := tr.zeroValue(true, m.Value)
	return v
}

// Returns the zero value of a custom type.
func (tr *transform) zeroOfType(name string) string {
	// In the actual function
	if tr.funcId != 0 {
		for block := tr.blockId; block >= 1; block-- {
			if _, ok := tr.typeZero[tr.funcId][block][name]; ok {
				return tr.typeZero[tr.funcId][block][name]
			}
		}
	}

	// Finally, search in the global variables (funcId = 0).
	for block := tr.blockId; block >= 0; block-- { // block until 0
		if _, ok := tr.typeZero[0][block][name]; ok {
			return tr.typeZero[0][block][name]
		}
	}
	//fmt.Printf("Function %d, block %d, name %s\n", tr.funcId, tr.blockId, name)
	panic("zeroOfType: type not found: " + name)
}

// === Checking
//

// Checks if a variable name is of a specific data type.
func (tr *transform) isType(t dataType, name string) bool {
	if name == "" {
		return false
	}

	name = strings.SplitN(name, "<<", 2)[0] // could have a tag

	for funcId := tr.funcId; funcId >= 0; funcId-- {
		for blockId := tr.blockId; blockId >= 0; blockId-- {
			if _, ok := tr.vars[funcId][blockId][name]; ok { // variable found
				switch t {
				case sliceType:
					if _, ok := tr.slices[funcId][blockId][name]; ok {
						return true
					}
				case mapType:
					if _, ok := tr.maps[funcId][blockId][name]; ok {
						return true
					}
				}

				return false
			}
		}
	}
	return false
}
