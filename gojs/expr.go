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
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

// Represents an expression.
type expression struct {
	tr            *transform
	*bytes.Buffer // sintaxis translated

	varName  string // variable name
	funcName string // function name

	//isFunc      bool // anonymous function
	isAddress   bool
	isEllipsis  bool
	isInitArray bool // initialization of array?
	isPointer   bool

	useIota       bool
	skipSemicolon bool
	hasError      bool

	// To handle comparisons
	isBasicLit     bool
	returnBasicLit bool

	lenArray []string // the lengths of an array
}

// Initializes a new type of "expression".
func (tr *transform) newExpression(iVar interface{}) *expression {
	var id string

	if iVar != nil {
		switch tVar := iVar.(type) {
		case *ast.Ident:
			id = tVar.Name
		case string:
			id = tVar
		}
	}

	return &expression{
		tr,
		new(bytes.Buffer),
		id,
		"",
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		make([]string, 0),
	}
}

// Transforms the Go expression.
func (e *expression) transform(expr ast.Expr) {
	switch typ := expr.(type) {

	// godoc go/ast ArrayType
	//  Lbrack token.Pos // position of "["
	//  Len    Expr      // Ellipsis node for [...]T array types, nil for slice types
	//  Elt    Expr      // element type
	case *ast.ArrayType:
		// Type checking
		if _, ok := typ.Elt.(*ast.Ident); ok {
			if e.tr.getExpression(typ.Elt).hasError {
				return
			}
		}
		if typ.Len == nil { // slice
			break
		}
		if _, ok := typ.Len.(*ast.Ellipsis); ok {
			e.isEllipsis = true
			break
		}

		if len(e.lenArray) != 0 {
			e.writeLoop()
			e.WriteString(fmt.Sprintf("{%s%s=", SP+e.varName, e.printArray()))
		}
		e.WriteString("[]")
		e.addLenArray(typ.Len)

		switch t := typ.Elt.(type) {
		case *ast.ArrayType: // multi-dimensional array
			e.transform(typ.Elt)
		case *ast.Ident, *ast.StarExpr: // the type is initialized
			init, isPointer := e.tr.initValue(true, typ.Elt)
			if isPointer { // remove '[', ']'
				init = init[1:len(init)-1]
			}

			e.writeLoop()
			e.WriteString(fmt.Sprintf("{%s=%s;%s}",
				SP+e.tr.lastVarName+e.printArray(), init, SP))

			if len(e.lenArray) > 1 {
				e.WriteString(strings.Repeat("}", len(e.lenArray)-1))
			}
		default:
			panic(fmt.Sprintf("*expression.transform: type unimplemented: %T", t))
		}
		e.skipSemicolon = true

	// godoc go/ast BasicLit
	//  Kind     token.Token // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING
	//  Value    string      // literal string
	case *ast.BasicLit:
		e.WriteString(typ.Value)
		e.isBasicLit = true

	// http://golang.org/doc/go_spec.html#Comparison_operators
	// https://developer.mozilla.org/en/JavaScript/Reference/Operators/Comparison_Operators
	//
	// godoc go/ast BinaryExpr
	//  X     Expr        // left operand
	//  Op    token.Token // operator
	//  Y     Expr        // right operand
	case *ast.BinaryExpr:
		op := typ.Op.String()
		isOpEq := false

		switch typ.Op {
		case token.EQL:
			op += "="
			isOpEq = true
		}

		if e.tr.isConst {
			e.transform(typ.X)
			e.WriteString(SP + op + SP)
			e.transform(typ.Y)
			break
		}

		stringify := false
		x := e.tr.getExpression(typ.X)
		y := e.tr.getExpression(typ.Y)

		// JavaScript only compares basic literals.
		if isOpEq && !x.isBasicLit && !x.returnBasicLit && !y.isBasicLit && !y.returnBasicLit {
			stringify = true
		}

		if stringify {
			e.WriteString("JSON.stringify(" + x.String() + ")")
		} else {
			e.WriteString(x.String())
		}

		e.WriteString(SP + op + SP)

		if stringify {
			e.WriteString("JSON.stringify(" + y.String() + ")")
		} else {
			e.WriteString(y.String())
		}

	// godoc go/ast CallExpr
	//  Fun      Expr      // function expression
	//  Args     []Expr    // function arguments; or nil
	case *ast.CallExpr:
		// === Library
		if call, ok := typ.Fun.(*ast.SelectorExpr); ok {
			e.transform(call)

			str := fmt.Sprintf("%s", e.tr.GetArgs(e.funcName, typ.Args))
			if e.funcName != "fmt.Sprintf" {
				str = "(" + str + ")"
			}

			e.WriteString(str)
			break
		}

		// === Conversion: []byte()
		if call, ok := typ.Fun.(*ast.ArrayType); ok {
			if call.Elt.(*ast.Ident).Name == "byte" {
				e.transform(typ.Args[0])
			} else {
				panic(fmt.Sprintf("call of conversion unimplemented: []%T()", call))
			}
			break
		}

		// === Built-in functions - golang.org/pkg/builtin/
		call := typ.Fun.(*ast.Ident).Name

		switch call {
		case "make":
			// Type checking
			if e.tr.getExpression(typ.Args[0]).hasError {
				return
			}

			switch argType := typ.Args[0].(type) {
			// For slice
			case *ast.ArrayType:
				length := e.tr.getExpression(typ.Args[1]).String()
				// TODO: add capacity
				//capacity := 0

				/*if len(typ.Args) == 3 {
					capacity = e.tr.getExpression(typ.Args[-1]).String()
				}*/

				init, _ := e.tr.initValue(true, argType.Elt)

				e.WriteString("[]")
				e.lenArray = append(e.lenArray, length)
				e.writeLoop()
				e.WriteString(fmt.Sprintf("{%s=%s;%s}",
					SP+e.tr.lastVarName+e.printArray(), init, SP))
				e.skipSemicolon = true

			case *ast.MapType:
				e.WriteString("{}") // or "new Object()"

			case *ast.ChanType:
				e.transform(typ.Fun)

			default:
				panic(fmt.Sprintf("call of 'make' unimplemented: %T", argType))
			}

		case "new":
			switch argType := typ.Args[0].(type) {
			case *ast.ArrayType:
				for _, arg := range typ.Args {
					e.transform(arg)
				}

			case *ast.Ident:
				value, _ := e.tr.initValue(true, argType)
				e.WriteString(value)

			default:
				panic(fmt.Sprintf("call of 'new' unimplemented: %T", argType))
			}

		// Conversion
		case "uint", "uint8", "uint16", "uint32",
			"int", "int8", "int16", "int32",
			"float32", "float64", "byte", "rune", "string":
			e.transform(typ.Args[0])
			e.returnBasicLit = true

		case "print", "println":
			e.WriteString(fmt.Sprintf("console.log(%s)", e.tr.GetArgs(call, typ.Args)))

		case "len":
			e.WriteString(fmt.Sprintf("%s.length", e.tr.getExpression(typ.Args[0])))
			e.returnBasicLit = true

		case "cap":
			e.WriteString("'cap'")

		case "panic":
			e.WriteString(fmt.Sprintf("throw new Error(%s)",
				e.tr.getExpression(typ.Args[0])))

		// === Not supported
		case "recover", "complex":
			e.tr.addError("%s: built-in function %s()",
				e.tr.fset.Position(typ.Fun.Pos()), call)
			e.tr.hasError = true
			return
		case "int64", "uint64":
			e.tr.addError("%s: conversion of type %s",
				e.tr.fset.Position(typ.Fun.Pos()), call)
			e.tr.hasError = true
			return

		// === Not implemented
		case "append", "close", "copy", "delete", "uintptr":
			panic(fmt.Sprintf("built-in call unimplemented: %s", call))

		// Defined functions
		default:
			args := ""

			for i, v := range typ.Args {
				if i != 0 {
					args += "," + SP
				}
				args += e.tr.getExpression(v).String()
			}

			e.WriteString(fmt.Sprintf("%s(%s)", call, args))
		}

	// godoc go/ast ChanType
	//  Begin token.Pos // position of "chan" keyword or "<-" (whichever comes first)
	//  Dir   ChanDir   // channel direction
	//  Value Expr      // value type
	case *ast.ChanType:
		e.tr.addError("%s: channel type", e.tr.fset.Position(typ.Pos()))
		e.tr.hasError = true
		return

	// godoc go/ast CompositeLit
	//  Type   Expr      // literal type; or nil
	//  Lbrace token.Pos // position of "{"
	//  Elts   []Expr    // list of composite elements; or nil
	//  Rbrace token.Pos // position of "}"
	case *ast.CompositeLit:
		switch compoType := typ.Type.(type) {
		case *ast.ArrayType:
			if !e.isInitArray {
				e.transform(typ.Type)
			}

			if e.isEllipsis {
				e.WriteString("[")
				e.writeElts(typ.Elts, typ.Lbrace, typ.Rbrace)
				e.WriteString("]")
				break
			}

			// For arrays initialized
			if len(typ.Elts) != 0 {
				if !e.isInitArray && compoType.Len != nil {
					e.WriteString(fmt.Sprintf("%s=%s", SP+e.varName+SP, SP))
					e.isInitArray = true
				}
				e.WriteString("[")
				e.writeElts(typ.Elts, typ.Lbrace, typ.Rbrace)
				e.WriteString("]")

				e.skipSemicolon = false
			}

		case *ast.Ident: // Custom types
			useField := false
			e.WriteString("new " + typ.Type.(*ast.Ident).Name)

			if len(typ.Elts) != 0 {
				// Specify the fields
				if _, ok := typ.Elts[0].(*ast.KeyValueExpr); ok {
					typeName := e.tr.lastVarName
					useField = true

					e.WriteString("()")

					for i, v := range typ.Elts {
						kv := v.(*ast.KeyValueExpr)

						if i != 0 {
							e.WriteString(",")
						} else {
							e.WriteString(";")
						}

						e.WriteString(fmt.Sprintf("%s.%s=%s",
							SP + typeName,
							e.tr.getExpression(kv.Key).String() + SP,
							SP + e.tr.getExpression(kv.Value).String(),
						))
					}
				}
			}
			if !useField {
				e.WriteString("(")
				e.writeElts(typ.Elts, typ.Lbrace, typ.Rbrace)
				e.WriteString(")")
			}

		case *ast.MapType:
			// Type checking
			if e.tr.getExpression(typ.Type).hasError {
				return
			}

			e.WriteString("{")
			e.writeElts(typ.Elts, typ.Lbrace, typ.Rbrace)
			e.WriteString("}")

		case nil:
			e.WriteString("[")
			e.writeElts(typ.Elts, typ.Lbrace, typ.Rbrace)
			e.WriteString("]")

		default:
			panic(fmt.Sprintf("'CompositeLit' unimplemented: %T", compoType))
		}

	// godoc go/ast Ellipsis
	//  Ellipsis token.Pos // position of "..."
	//  Elt      Expr      // ellipsis element type (parameter lists only); or nil
	//case *ast.Ellipsis:

	// http://golang.org/doc/go_spec.html#Function_literals
	// https://developer.mozilla.org/en/JavaScript/Reference/Functions_and_function_scope#Function_constructor_vs._function_declaration_vs._function_expression
	// godoc go/ast FuncLit
	//
	//  Type *FuncType  // function type
	//  Body *BlockStmt // function body
	case *ast.FuncLit:
		e.transform(typ.Type)
		e.tr.getStatement(typ.Body)

	// godoc go/ast FuncType
	//  Func    token.Pos  // position of "func" keyword
	//  Params  *FieldList // (incoming) parameters; or nil
	//  Results *FieldList // (outgoing) results; or nil
	case *ast.FuncType:
		//e.isFunc = true
		e.tr.writeFunc(nil, nil, typ)

	// godoc go/ast Ident
	//  Name    string    // identifier name
	case *ast.Ident:
		name := typ.Name

		switch name {
		case "iota":
			e.WriteString(IOTA)
			e.useIota = true

		// Undefined value in array / slice
		case "_":
			if len(e.lenArray) == 0 {
				e.WriteString(name)
			}
		// https://developer.mozilla.org/en/JavaScript/Reference/Global_Objects/undefined
		case "nil":
			e.WriteString("undefined")
			e.isBasicLit = true

		// Not supported
		case "int64", "uint64", "complex64", "complex128":
			e.tr.addError("%s: %s type", e.tr.fset.Position(typ.Pos()), name)
			e.tr.hasError = true
		// Not implemented
		case "uintptr":
			e.tr.addError("%s: unimplemented type %q", e.tr.fset.Position(typ.Pos()), name)
			e.tr.hasError = true

		default:
			if e.isPointer { // `*x` => `x[0]`
				name += "[0]"
			} else if e.isAddress { // `&x` => `x`
				e.tr.addPointer(name)
			}

			e.WriteString(name)
		}

	// godoc go/ast IndexExpr
	// Represents an expression followed by an index.
	//  X      Expr      // expression
	//  Lbrack token.Pos // position of "["
	//  Index  Expr      // index expression
	//  Rbrack token.Pos // position of "]"
	case *ast.IndexExpr:
		e.transform(typ.X)
		e.WriteString("[")
		e.transform(typ.Index)
		e.WriteString("]")

	// godoc go/ast InterfaceType
	//  Interface  token.Pos  // position of "interface" keyword
	//  Methods    *FieldList // list of methods
	//  Incomplete bool       // true if (source) methods are missing in the Methods list
	case *ast.InterfaceType: // TODO: review

	// godoc go/ast KeyValueExpr
	//  Key   Expr
	//  Colon token.Pos // position of ":"
	//  Value Expr
	case *ast.KeyValueExpr:
		e.transform(typ.Key)
		e.WriteString(":" + SP)
		e.transform(typ.Value)

	// godoc go/ast MapType
	//  Map   token.Pos // position of "map" keyword
	//  Key   Expr
	//  Value Expr
	case *ast.MapType:
		// Type checking
		e.tr.getExpression(typ.Key)
		e.tr.getExpression(typ.Value)

	// godoc go/ast ParenExpr
	//  Lparen token.Pos // position of "("
	//  X      Expr      // parenthesized expression
	//  Rparen token.Pos // position of ")"
	case *ast.ParenExpr:
		e.transform(typ.X)

	// godoc go/ast SelectorExpr
	//   X   Expr   // expression
	//   Sel *Ident // field selector
	case *ast.SelectorExpr:
		isPkg := false
		x := ""

		switch t := typ.X.(type) {
		case *ast.Ident:
			x = t.Name
			if e.tr.recvVar != "" && x == e.tr.recvVar {
				x = "this"
			}
		case *ast.IndexExpr:
			e.transform(t)
			e.WriteString("." + typ.Sel.Name)
			return
		default:
			panic(fmt.Sprintf("'SelectorExpr': unimplemented: %T", t))
		}

		goName := x + "." + typ.Sel.Name

		// Check is the selector is a package
		for _, v := range validImport {
			if v == x {
				isPkg = true
				break
			}
		}

		// Check if it can be transformed to its equivalent in JavaScript.
		if isPkg {
			jsName, ok := Function[goName]
			if !ok {
				jsName, ok = Constant[goName]
			}

			if !ok {
				e.tr.addError(fmt.Errorf("%s: %q not supported in JS",
					e.tr.fset.Position(typ.Sel.Pos()), goName))
				e.tr.hasError = true
				break
			}

			e.funcName = goName
			e.WriteString(jsName)
		} else {
			/*if _, ok := e.tr.types[x]; !ok {
				panic("selector: " + x)
			}*/

			e.WriteString(goName)
		}

	// godoc go/ast SliceExpr
	//  X      Expr      // expression
	//  Lbrack token.Pos // position of "["
	//  Low    Expr      // begin of slice range; or nil
	//  High   Expr      // end of slice range; or nil
	//  Rbrack token.Pos // position of "]"
	case *ast.SliceExpr:
		slice := "0"

		if typ.Low != nil {
			slice = typ.Low.(*ast.BasicLit).Value
		}
		if typ.High != nil {
			slice += "," + typ.High.(*ast.BasicLit).Value
		}

		e.WriteString(fmt.Sprintf("%s.slice(%s)", typ.X.(*ast.Ident), slice))

	// godoc go/ast StructType
	//  Struct     token.Pos  // position of "struct" keyword
	//  Fields     *FieldList // list of field declarations
	//  Incomplete bool       // true if (source) fields are missing in the Fields list
	case *ast.StructType:

	// godoc go/ast StarExpr
	//  Star token.Pos // position of "*"
	//  X    Expr      // operand
	case *ast.StarExpr:
		e.isPointer = true
		e.transform(typ.X)

	// godoc go/ast UnaryExpr
	//  OpPos token.Pos   // position of Op
	//  Op    token.Token // operator
	//  X     Expr        // operand
	case *ast.UnaryExpr:
		writeOp := true
		op := typ.Op.String()

		switch typ.Op {
		// Bitwise complement
		case token.XOR:
			op = "~"
		// Address operator
		case token.AND:
			e.isAddress = true
			writeOp = false
		case token.ARROW:
			e.tr.addError("%s: channel operator", e.tr.fset.Position(typ.OpPos))
			e.tr.hasError = true
			return
		}

		if writeOp {
			e.WriteString(op)
		}
		e.transform(typ.X)

	// The type has not been indicated
	case nil:

	default:
		panic(fmt.Sprintf("unimplemented: %T", expr))
	}
}

//
// === Utility

// Appends a new length of array.
func (e *expression) addLenArray(expr ast.Expr) {
	e.lenArray = append(e.lenArray, e.tr.getExpression(expr).String())
}

// Returns the values of an array formatted like "[i][j]..."
func (e *expression) printArray() string {
	a := ""

	for i := 0; i < len(e.lenArray); i++ {
		vArray := string('i' + i)
		a += fmt.Sprintf("[%s]", vArray)
	}
	return a
}

// Writes the loop for the last length of the array.
func (e *expression) writeLoop() {
	iArray := len(e.lenArray) - 1  // index of array
	vArray := string('i' + iArray) // variable's name for the loop

	e.WriteString(fmt.Sprintf(";%sfor%s(var %s=0;%s<%s;%s++)",
		SP, SP, vArray, SP+vArray, e.lenArray[iArray], SP+vArray))
}

// Writes the list of composite elements.
func (e *expression) writeElts(elts []ast.Expr, Lbrace, Rbrace token.Pos) {
	firstPos := e.tr.getLine(Lbrace)
	posOldElt := firstPos
	posNewElt := 0

	for i, el := range elts {
		posNewElt = e.tr.getLine(el.Pos())

		if i != 0 {
			e.WriteString(",")
		}
		if posNewElt != posOldElt {
			e.WriteString(strings.Repeat(NL, posNewElt - posOldElt))
			e.WriteString(strings.Repeat(TAB, e.tr.tabLevel + 1))
		} else if i != 0 { // in the same line
			e.WriteString(SP)
		}

		e.transform(el)
		posOldElt = posNewElt
	}

	// The right brace
	posNewElt = e.tr.getLine(Rbrace)
	if posNewElt != posOldElt {
		e.WriteString(strings.Repeat(NL, posNewElt - posOldElt))
		e.WriteString(strings.Repeat(TAB, e.tr.tabLevel))
	}

	e.tr.line += posNewElt - firstPos // update the global position
}
