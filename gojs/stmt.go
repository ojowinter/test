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
	"strings"
)

// Represents data for a statement.
type dataStmt struct {
	tabLevel  int  // tabulation level
	lenCase   int  // number of "case" statements
	iCase     int  // index in "case" statements
	wasReturn bool // the last statement was "return"?
}

// Transforms the Go statement.
func (tr *transform) getStatement(stmt ast.Stmt) {
	switch typ := stmt.(type) {

	// http://golang.org/doc/go_spec.html#Arithmetic_operators
	// https://developer.mozilla.org/en/JavaScript/Reference/Operators/Assignment_Operators
	//
	// http://golang.org/pkg/go/ast/#AssignStmt || godoc go/ast AssignStmt
	//  Lhs    []Expr
	//  TokPos token.Pos   // position of Tok
	//  Tok    token.Token // assignment token, DEFINE
	//  Rhs    []Expr
	case *ast.AssignStmt:
		var assign string
		var isNew, needChange bool

		switch typ.Tok {
		case token.DEFINE:
			assign = "="
			isNew = true
		case token.ASSIGN,
			token.ADD_ASSIGN, token.SUB_ASSIGN, token.MUL_ASSIGN, token.QUO_ASSIGN,
			token.REM_ASSIGN,
			token.AND_ASSIGN, token.OR_ASSIGN, token.XOR_ASSIGN, token.SHL_ASSIGN,
			token.SHR_ASSIGN:

			assign = typ.Tok.String()
		case token.AND_NOT_ASSIGN:
			needChange = true

		default:
			panic(fmt.Sprintf("token unimplemented: %s", typ.Tok))
		}

		if isNew {
			tr.WriteString("var ")
		}

		isFirst := true
		for i, v := range typ.Lhs {
			lIdent := getExpression(v)

			if lIdent == "_" {
				continue
			}
			rIdent := getExpression(typ.Rhs[i])

			if isFirst {
				isFirst = false
			} else {
				tr.WriteString("," + SP)
			}

			tr.WriteString(lIdent)
			// Skip empty strings
			if rIdent != EMPTY {
				if !needChange {
					tr.WriteString(SP + assign + SP + rIdent)
				} else {
					if typ.Tok == token.AND_NOT_ASSIGN {
						tr.WriteString(SP + "=" + SP + lIdent + SP + "&" + SP +
							rIdent + SP + "^" + SP + lIdent)
					}
				}
			}
		}
		tr.WriteString(";")

	// http://golang.org/pkg/go/ast/#BlockStmt || godoc go/ast BlockStmt
	//  Lbrace token.Pos // position of "{"
	//  List   []Stmt
	//  Rbrace token.Pos // position of "}"
	case *ast.BlockStmt:
		tr.WriteString("{")

		for _, v := range typ.List {
			isCase := false

			// Don't insert tabulation in "case" clauses
			if _, ok := v.(*ast.CaseClause); ok {
				isCase = true
			} else {
				tr.tabLevel++
			}

			tr.addLine(v.Pos())
			tr.WriteString(strings.Repeat(TAB, tr.tabLevel))
			tr.getStatement(v)

			if !isCase {
				tr.tabLevel--
			}
		}

		tr.addLine(typ.Rbrace)
		tr.WriteString(strings.Repeat(TAB, tr.tabLevel) + "}")

	// http://golang.org/pkg/go/ast/#CaseClause || godoc go/ast CaseClause
	//  Case  token.Pos // position of "case" or "default" keyword
	//  List  []Expr    // list of expressions or types; nil means default case
	//  Colon token.Pos // position of ":"
	//  Body  []Stmt    // statement list; or nil
	case *ast.CaseClause:
		tr.wasReturn = false // to check the last statement
		tr.iCase++
		tr.addLine(typ.Case)

		if typ.List != nil {
			for i, expr := range typ.List {
				if i != 0 {
					tr.WriteString(SP)
				}
				tr.WriteString(fmt.Sprintf("case %s:", getExpression(expr)))
			}
		} else {
			tr.WriteString("default:")

			if tr.iCase != tr.lenCase {
				tr.addWarning("%s: 'default' clause above 'case' clause in switch statement",
					tr.fset.Position(typ.Pos()))
			}
		}

		if typ.Body != nil {
			for _, v := range typ.Body {
				if ok := tr.addLine(v.Pos()); ok {
					tr.WriteString(strings.Repeat(TAB, tr.tabLevel+1))
				} else {
					tr.WriteString(SP)
				}
				tr.getStatement(v)
			}
		}

		// Skip last "case" statement
		if !tr.wasReturn && tr.iCase != tr.lenCase {
			tr.WriteString(SP + "break;")
		}

	// http://golang.org/pkg/go/ast/#ExprStmt || godoc go/ast ExprStmt
	//  X Expr // expression
	/*case *ast.ExprStmt:
	tr.WriteString(getExpression(typ.X))*/

	// http://golang.org/pkg/go/ast/#ForStmt || godoc go/ast ForStmt
	//  For  token.Pos // position of "for" keyword
	//  Init Stmt      // initialization statement; or nil
	//  Cond Expr      // condition; or nil
	//  Post Stmt      // post iteration statement; or nil
	//  Body *BlockStmt
	case *ast.ForStmt:
		tr.WriteString("for" + SP + "(")

		if typ.Init != nil {
			tr.getStatement(typ.Init)
		} else {
			tr.WriteString(";")
		}
		tr.WriteString(SP)

		if typ.Cond != nil {
			tr.WriteString(getExpression(typ.Cond))
		}
		tr.WriteString(";" + SP)

		if typ.Post != nil {
			tr.getStatement(typ.Post)
		}

		tr.WriteString(")" + SP)
		tr.getStatement(typ.Body)

	// http://golang.org/pkg/go/ast/#GoStmt || godoc go/ast GoStmt
	//  Go   token.Pos // position of "go" keyword
	//  Call *CallExpr
	case *ast.GoStmt:
		tr.addError("%s: goroutine", tr.fset.Position(typ.Go))

	// http://golang.org/doc/go_spec.html#If_statements
	// https://developer.mozilla.org/en/JavaScript/Reference/Statements/if...else
	//
	// http://golang.org/pkg/go/ast/#IfStmt || godoc go/ast IfStmt
	//  If   token.Pos // position of "if" keyword
	//  Init Stmt      // initialization statement; or nil
	//  Cond Expr      // condition
	//  Body *BlockStmt
	//  Else Stmt // else branch; or nil
	case *ast.IfStmt:
		if typ.Init != nil {
			tr.getStatement(typ.Init)
			tr.WriteString(SP)
		}

		tr.WriteString(fmt.Sprintf("if%s(%s)%s", SP, getExpression(typ.Cond), SP))
		tr.getStatement(typ.Body)

		if typ.Else != nil {
			tr.WriteString(SP + "else ")
			tr.getStatement(typ.Else)
		}

	// http://golang.org/pkg/go/ast/#IncDecStmt || godoc go/ast IncDecStmt
	//  X      Expr
	//  TokPos token.Pos   // position of Tok
	//  Tok    token.Token // INC or DEC
	case *ast.IncDecStmt:
		tr.WriteString(getExpression(typ.X) + typ.Tok.String())

	// http://golang.org/pkg/go/ast/#RangeStmt || godoc go/ast RangeStmt
	//  For        token.Pos   // position of "for" keyword
	//  Key, Value Expr        // Value may be nil
	//  TokPos     token.Pos   // position of Tok
	//  Tok        token.Token // ASSIGN, DEFINE
	//  X          Expr        // value to range over
	//  Body       *BlockStmt
	/*case *ast.RangeStmt:
		key := getExpression(typ.Key)

		value := ""
		if typ.Value != nil {
			value = getExpression(typ.Value)
		}

		expr := getExpression(typ.X)

		switch t := typ.X.(type) {
		case *ast.ArrayType: // string
			init := key + "SP" + "=" + SP + "0"  // initialization

			if typ.Tok == token.DEFINE {
				init = "var " + init

				if typ.Value != nil {
					init += "," + SP + value
				}
			}

			tr.WriteString(fmt.Sprintf("for%s(%s;%s;%s%s)", SP, init, SP, SP, expr))

		case *ast.Ident:
			fmt.Printf("%T : %v\n", t.Obj,t.Obj.Data)
		default:
			fmt.Printf("%T\n", t)
		}*/

	// http://golang.org/doc/go_spec.html#Return_statements
	// https://developer.mozilla.org/en/JavaScript/Reference/Statements/return
	//
	// http://golang.org/pkg/go/ast/#ReturnStmt || godoc go/ast ReturnStmt
	//  Return  token.Pos // position of "return" keyword
	//  Results []Expr    // result expressions; or nil
	case *ast.ReturnStmt:
		tr.wasReturn = true

		if typ.Results == nil {
			tr.WriteString("return;")
			break
		}

		if len(typ.Results) != 1 {
			tr.addError("%s: return multiple values", tr.fset.Position(typ.Return))
			break
		}
		tr.WriteString("return " + getExpression(typ.Results[0]) + ";")

	// http://golang.org/doc/go_spec.html#Switch_statements
	// https://developer.mozilla.org/en/JavaScript/Reference/Statements/switch
	//
	// http://golang.org/pkg/go/ast/#SwitchStmt || godoc go/ast SwitchStmt
	//  Switch token.Pos  // position of "switch" keyword
	//  Init   Stmt       // initialization statement; or nil
	//  Tag    Expr       // tag expression; or nil
	//  Body   *BlockStmt // CaseClauses only
	case *ast.SwitchStmt:
		tag := ""
		tr.lenCase = len(typ.Body.List)
		tr.iCase = 0

		if typ.Init != nil {
			tr.getStatement(typ.Init)
			tr.WriteString(SP)
		}
		if typ.Tag != nil {
			tag = getExpression(typ.Tag)
		} else {
			tag = "1" // true
		}

		tr.WriteString(fmt.Sprintf("switch%s(%s)%s", SP, tag, SP))
		tr.getStatement(typ.Body)

	default:
		panic(fmt.Sprintf("unimplemented: %T", stmt))
	}
}
