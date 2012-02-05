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
)

// Functions
//
// http://golang.org/doc/go_spec.html#Function_declarations
// https://developer.mozilla.org/en/JavaScript/Reference/Statements/function
func (tr *transform) getFunc(decl *ast.FuncDecl) {
	// godoc go/ast FuncDecl
	//  Doc  *CommentGroup // associated documentation; or nil
	//  Recv *FieldList    // receiver (methods); or nil (functions)
	//  Name *Ident        // function/method name
	//  Type *FuncType     // position of Func keyword, parameters and results
	//  Body *BlockStmt    // function body; or nil (forward declaration)

	// Check empty functions
	if len(decl.Body.List) == 0 {
		return
	}

	isFuncInit := false // function init()

	// === Initialization to save variables created on this function
	if decl.Name != nil { // discard literal functions
		tr.funcTotal++
		tr.funcId = tr.funcTotal
		tr.blockId = 0

		tr.vars[tr.funcId] = make(map[int]map[string]bool)
		tr.addr[tr.funcId] = make(map[int]map[string]bool)
		tr.typeZero[tr.funcId] = make(map[int]map[string]string)
		tr.mapKeys[tr.funcId] = make(map[int]map[string]map[string]struct{})
		tr.mapZero[tr.funcId] = make(map[int]map[string]map[int]string)
	}
	// ===

	tr.addLine(decl.Pos())
	tr.addIfExported(decl.Name)

	if decl.Name.Name != "init" {
		tr.writeFunc(decl.Recv, decl.Name, decl.Type)
	} else {
		isFuncInit = true
		tr.WriteString("(function()" + SP)
	}

	tr.getStatement(decl.Body)

	if isFuncInit {
		tr.WriteString("());")
	}

	// At exiting of the function, it returns at the global scope.
	if decl.Name != nil {
		tr.funcId = 0
		tr.blockId = 0
	}
	if decl.Recv != nil {
		tr.recvVar = ""
	}
}

// godoc go/ast FuncType
//  Func    token.Pos  // position of "func" keyword
//  Params  *FieldList // (incoming) parameters; or nil
//  Results *FieldList // (outgoing) results; or nil

// godoc go/ast FieldList
//  Opening token.Pos // position of opening parenthesis/brace, if any
//  List    []*Field  // field list; or nil
//  Closing token.Pos // position of closing parenthesis/brace, if any

// godoc go/ast Field
//  Doc     *CommentGroup // associated documentation; or nil
//  Names   []*Ident      // field/method/parameter names; or nil if anonymous field
//  Type    Expr          // field/method/parameter type
//  Tag     *BasicLit     // field tag; or nil
//  Comment *CommentGroup // line comments; or nil

// Writes the function declaration.
func (tr *transform) writeFunc(recv *ast.FieldList, name *ast.Ident, typ *ast.FuncType) {
	if recv != nil { // method
		field := recv.List[0]
		tr.recvVar = field.Names[0].Name

		tr.WriteString(fmt.Sprintf("%s.prototype.%s%s=%sfunction(%s)%s",
			field.Type, name,SP, SP, joinParams(typ), SP))
	} else if name != nil {
		tr.WriteString(fmt.Sprintf("function %s(%s)%s", name, joinParams(typ), SP))
	} else { // Literal function
		tr.WriteString(fmt.Sprintf("%s=%sfunction(%s)%s", SP, SP, joinParams(typ), SP))
	}

	// Return multiple values
	declResults, declReturn := tr.joinResults(typ)

	if declResults != "" {
		tr.WriteString("{" + SP + declResults)
		tr.skipLbrace = true
		tr.results = declReturn
	} else {
		tr.results = ""
	}
}

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
func (tr *transform) joinResults(f *ast.FuncType) (decl, ret string) {
	isFirst := true
	isMultiple := false

	if f.Results == nil {
		return
	}

	for _, list := range f.Results.List {
		if list.Names == nil {
			continue
		}

		value, _ := tr.zeroValue(true, list.Type)

		for _, v := range list.Names {
			if !isFirst {
				decl += "," + SP
				ret += "," + SP
				isMultiple = true
			} else {
				isFirst = false
			}

			decl += fmt.Sprintf("%s=%s", v.Name+SP, SP+value)
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
