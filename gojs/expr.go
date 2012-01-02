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
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

// Represents an expression.
type expression struct {
	tr            *transform
	*bytes.Buffer // sintaxis translated

	ident      string // variable's name
	funcName   string // function's name
	useIota    bool
	isNegative bool
	isAddress  bool
	isPointer  bool
	isFunc     bool // anonymous function

	lenArray int      // store length of array; to use in case of ellipsis [...]
	valArray []string // store the last values of an array
}

// Initializes a new type of "expression".
func (tr *transform) newExpression(ident *ast.Ident) *expression {
	var id string

	if ident != nil {
		id = ident.Name
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
		0,
		make([]string, 0),
	}
}

// Returns the Go expression in JavaScript.
func (tr *transform) getExpression(expr ast.Expr) string {
	e := tr.newExpression(nil)
	e.transform(expr)

	return e.String()
}

// Returns the values of an array formatted like "[i0][i1]..."
func (e *expression) printArray() string {
	a := ""

	for i := 0; i < len(e.valArray); i++ {
		vArray := "i" + strconv.Itoa(i)
		a = fmt.Sprintf("%s[%s]", a, vArray)
	}
	return a
}

// Transforms the Go expression.
func (e *expression) transform(expr ast.Expr) {
	switch typ := expr.(type) {

	// http://golang.org/pkg/go/ast/#ArrayType || godoc go/ast ArrayType
	//  Len    Expr      // Ellipsis node for [...]T array types, nil for slice types
	//  Elt    Expr      // element type
	case *ast.ArrayType:
		if typ.Len == nil { // slice
			break
		}

		if len(e.valArray) == 0 {
			e.WriteString("new Array(")

			if e.lenArray != 0 { // ellipsis
				e.WriteString(strconv.Itoa(e.lenArray))
			} else {
				e.transform(typ.Len)
			}

			e.WriteString(")")
		} else {
			iArray := len(e.valArray) - 1        // index of array
			vArray := "i" + strconv.Itoa(iArray) // variable's name for the loop

			e.WriteString(fmt.Sprintf(
				";%sfor%s(var %s=0;%s<%s;%s++){%s=new Array(",
				SP, SP, vArray, SP+vArray, e.valArray[iArray], SP+vArray,
				SP+e.ident+e.printArray()))
			e.transform(typ.Len)
			e.WriteString(")")
		}

		if _, ok := typ.Elt.(*ast.ArrayType); ok {
			e.transform(typ.Elt)
		} else if len(e.valArray) > 1 {
			e.WriteString(";" + SP + strings.Repeat("}", len(e.valArray)-1))
		}

	// http://golang.org/pkg/go/ast/#BasicLit || godoc go/ast BasicLit
	//  Kind     token.Token // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING
	//  Value    string      // literal string
	case *ast.BasicLit:
		e.WriteString(typ.Value)

		// === Add the value
		sign := ""

		if e.isNegative {
			sign = "-"
		}
		e.valArray = append(e.valArray, sign+typ.Value)

	// http://golang.org/doc/go_spec.html#Comparison_operators
	// https://developer.mozilla.org/en/JavaScript/Reference/Operators/Comparison_Operators
	//
	// http://golang.org/pkg/go/ast/#BinaryExpr || godoc go/ast BinaryExpr
	//  X     Expr        // left operand
	//  Op    token.Token // operator
	//  Y     Expr        // right operand
	case *ast.BinaryExpr:
		op := typ.Op.String()

		switch typ.Op {
		case token.EQL:
			op += "="
		}

		e.transform(typ.X)
		e.WriteString(SP + op + SP)
		e.transform(typ.Y)

	// http://golang.org/pkg/go/ast/#CallExpr || godoc go/ast CallExpr
	//  Fun      Expr      // function expression
	//  Args     []Expr    // function arguments; or nil
	case *ast.CallExpr:
		// === Library
		if call, ok := typ.Fun.(*ast.SelectorExpr); ok {
			e.transform(call)
			e.WriteString(fmt.Sprintf("(%s)", e.tr.GetArgs(e.funcName, typ.Args)))
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

		// === Built-in functions
		call := typ.Fun.(*ast.Ident).Name

		switch call {
		case "make":
			switch argType := typ.Args[0].(type) {
			default:
				panic(fmt.Sprintf("call of 'make' unimplemented: %T", argType))

			// For slice
			case *ast.ArrayType:
				e.WriteString("new Array(")
				e.transform(typ.Args[len(typ.Args)-1]) // capacity
				e.WriteString(")")

			// The second argument (in Args), if any, is the capacity which
			// is not useful in JS since it is dynamic.
			case *ast.MapType:
				e.WriteString("{};") // or "new Object()"
			}

		case "new":
			switch argType := typ.Args[0].(type) {
			default:
				panic(fmt.Sprintf("call of 'new' unimplemented: %T", argType))

			case *ast.ArrayType:
				for _, arg := range typ.Args {
					e.transform(arg)
				}
			}

		// Conversion
		case "uint", "uint8", "uint16", "uint32",
			"int", "int8", "int16", "int32",
			"float32", "float64", "byte", "rune", "string":
			e.transform(typ.Args[0])

		case "print", "println":
			e.WriteString(fmt.Sprintf("%s(%s)",
				Function[call], e.tr.GetArgs(call, typ.Args)))

		// Not supported
		case "panic", "recover":
			e.tr.addError("%s: built-in function %s()",
				e.tr.fset.Position(typ.Fun.Pos()), call)

		// Not implemented - http://golang.org/pkg/builtin/
		case "append", "cap", "close", "copy", "delete", "len", "uintptr":
			panic(fmt.Sprintf("built-in call unimplemented: %s", call))

		// Defined functions
		default:
			args := ""

			for i, v := range typ.Args {
				if i != 0 {
					args += "," + SP
				}
				args += e.tr.getExpression(v)
			}

			e.WriteString(fmt.Sprintf("%s(%s)", call, args))
		}

	// http://golang.org/pkg/go/ast/#CompositeLit || godoc go/ast CompositeLit
	//  Type   Expr      // literal type; or nil
	//  Elts   []Expr    // list of composite elements; or nil
	case *ast.CompositeLit:
		switch compoType := typ.Type.(type) {
		default:
			panic(fmt.Sprintf("'CompositeLit' unimplemented: %s", compoType))

		case *ast.ArrayType:
			e.lenArray = len(typ.Elts) // for ellipsis
			e.transform(typ.Type)

			// For arrays initialized
			if len(typ.Elts) != 0 {
				if compoType.Len == nil {
					e.WriteString("[")
				} else {
					e.WriteString(fmt.Sprintf(";%s=%s[", SP+e.ident+SP, SP))
				}

				for i, el := range typ.Elts {
					if i != 0 {
						e.WriteString("," + SP)
					}
					e.transform(el)
				}
				e.WriteString("]")
			}

		// http://golang.org/pkg/go/ast/#MapType || godoc go/ast MapType
		//  Key   Expr
		//  Value Expr
		case *ast.MapType:
			lenElts := len(typ.Elts) - 1
			e.WriteString("{")

			for i, el := range typ.Elts {
				e.transform(el)

				if i != lenElts {
					e.WriteString("," + SP)
				}
			}
			e.WriteString("};")
		}

	// http://golang.org/pkg/go/ast/#Ellipsis || godoc go/ast Ellipsis
	//  Elt      Expr      // ellipsis element type (parameter lists only); or nil
	//case *ast.Ellipsis:

	// http://golang.org/doc/go_spec.html#Function_literals
	// https://developer.mozilla.org/en/JavaScript/Reference/Functions_and_function_scope#Function_constructor_vs._function_declaration_vs._function_expression
	// http://golang.org/pkg/go/ast/#FuncLit || godoc go/ast FuncLit
	//
	//  Type *FuncType  // function type
	//  Body *BlockStmt // function body
	case *ast.FuncLit:
		e.isFunc = true
		e.transform(typ.Type)
		e.tr.getStatement(typ.Body)

	// http://golang.org/pkg/go/ast/#FuncType || godoc go/ast FuncType
	//  Func    token.Pos  // position of "func" keyword
	//  Params  *FieldList // (incoming) parameters; or nil
	//  Results *FieldList // (outgoing) results; or nil
	case *ast.FuncType:
		e.tr.WriteString(fmt.Sprintf("function(%s)%s", joinParams(typ), SP))

	// http://golang.org/pkg/go/ast/#Ident || godoc go/ast Ident
	//  Name    string    // identifier name
	case *ast.Ident:
		name := typ.Name

		if name == "iota" {
			e.WriteString(IOTA)
			e.useIota = true
			break
		}
		// Undefined value in array / slice
		if name == "_" && len(e.valArray) != 0 {
			break
		}

		// https://developer.mozilla.org/en/JavaScript/Reference/Global_Objects/undefined
		if name == "nil" {
			name = "undefined"
		} else if e.isPointer { // `*x` => `x[0]`
			name += "[0]"
		} else if e.isAddress { // `&x` => `x=[x]`
			var pointer *[]string

			if e.isFunc {
				pointer = &e.tr.funcPointer
			} else {
				pointer = &e.tr.globPointer
			}

			found := false
			for _, v := range *pointer {
				if v == name {
					found = true
					break
				}
			}

			if !found {
				*pointer = append(*pointer, name)
				name = fmt.Sprintf("%s=[%s]", name, name)
			}
		}

		e.WriteString(name)

	// http://golang.org/pkg/go/ast/#KeyValueExpr || godoc go/ast KeyValueExpr
	//  Key   Expr
	//  Value Expr
	case *ast.KeyValueExpr:
		e.transform(typ.Key)
		e.WriteString(":" + SP)
		e.transform(typ.Value)

	// http://golang.org/pkg/go/ast/#ParenExpr || godoc go/ast ParenExpr
	//  X      Expr      // parenthesized expression
	case *ast.ParenExpr:
		e.transform(typ.X)

	// http://golang.org/pkg/go/ast/#SelectorExpr || godoc go/ast SelectorExpr
	//   X   Expr   // expression
	//   Sel *Ident // field selector
	//
	// 'X' have the import name, and 'Sel' the constant of function name.
	case *ast.SelectorExpr:
		goName, jsName, err := e.tr.checkLib(typ)
		if err != nil {
			e.tr.addError(err)
			break
		}

		e.funcName = goName
		e.WriteString(jsName)

	// http://golang.org/pkg/go/ast/#StructType || godoc go/ast StructType
	//  Struct     token.Pos  // position of "struct" keyword
	//  Fields     *FieldList // list of field declarations
	//  Incomplete bool       // true if (source) fields are missing in the Fields list
	/*case *ast.StructType:*/

	// http://golang.org/pkg/go/ast/#StarExpr || godoc go/ast StarExpr
	//  X    Expr      // operand
	case *ast.StarExpr:
		e.isPointer = true
		e.transform(typ.X)

	// http://golang.org/pkg/go/ast/#UnaryExpr || godoc go/ast UnaryExpr
	//  Op    token.Token // operator
	//  X     Expr        // operand
	case *ast.UnaryExpr:
		writeOp := true
		op := typ.Op.String()

		switch typ.Op {
		case token.SUB:
			e.isNegative = true
		// Bitwise complement
		case token.XOR:
			op = "~"
		// Address operator
		case token.AND:
			e.isAddress = true
			writeOp = false
		}

		if writeOp {
			e.WriteString(op)
		}
		e.transform(typ.X)

	default:
		panic(fmt.Sprintf("unimplemented: %T", expr))
	}
}
