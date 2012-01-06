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
	tabLevel int // tabulation level
	lenCase  int // number of "case" statements
	iCase    int // index in "case" statements

	wasFallthrough bool // the last statement was "fallthrough"?
	wasReturn      bool // the last statement was "return"?
	skipLbrace     bool // left brace

	isSwitch   bool
	switchInit string
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
		// Can not be indicated the variable's type in the assignment.
		tr.writeValues(typ.Lhs, typ.Rhs, nil, typ.Tok, false)

	// http://golang.org/doc/go_spec.html#Blocks
	// https://developer.mozilla.org/en/JavaScript/Reference/Statements/block
	//
	// http://golang.org/pkg/go/ast/#BlockStmt || godoc go/ast BlockStmt
	//  Lbrace token.Pos // position of "{"
	//  List   []Stmt
	//  Rbrace token.Pos // position of "}"
	case *ast.BlockStmt:
		if !tr.skipLbrace {
			tr.WriteString("{")
		} else {
			tr.skipLbrace = false
		}

		for i, v := range typ.List {
			skipTab := false

			// Don't insert tabulation in both "case", "label" clauses
			switch v.(type) {
			case *ast.CaseClause, *ast.LabeledStmt:
				skipTab = true
			default:
				tr.tabLevel++
			}

			if tr.addLine(v.Pos()) {
				tr.WriteString(strings.Repeat(TAB, tr.tabLevel))
			} else if i == 0 {
				tr.WriteString(SP)
			}
			tr.getStatement(v)

			if !skipTab {
				tr.tabLevel--
			}
		}

		if tr.addLine(typ.Rbrace) {
			tr.WriteString(strings.Repeat(TAB, tr.tabLevel))
		} else {
			tr.WriteString(SP)
		}
		tr.WriteString("}")

	// http://golang.org/pkg/go/ast/#BranchStmt || godoc go/ast BranchStmt
	//  TokPos token.Pos   // position of Tok
	//  Tok    token.Token // keyword token (BREAK, CONTINUE, GOTO, FALLTHROUGH)
	//  Label  *Ident      // label name; or nil
	case *ast.BranchStmt:
		/*label := ";"
		if typ.Label != nil {
			label = SP + typ.Label.Name + ";"
		}*/

		tr.addLine(typ.TokPos)

		switch typ.Tok {
		// http://golang.org/doc/go_spec.html#Break_statements
		// https://developer.mozilla.org/en/JavaScript/Reference/Statements/break
		case token.BREAK:
			tr.WriteString("break;")
		// http://golang.org/doc/go_spec.html#Continue_statements
		// https://developer.mozilla.org/en/JavaScript/Reference/Statements/continue
		case token.CONTINUE:
			tr.WriteString("continue;")
		// http://golang.org/doc/go_spec.html#Goto_statements
		// http://golang.org/doc/go_spec.html#Fallthrough_statements
		case token.FALLTHROUGH:
			tr.wasFallthrough = true
		case token.GOTO: // not used since "label" is not transformed
			tr.addError("%s: goto statement", tr.fset.Position(typ.TokPos))
		}

	// http://golang.org/pkg/go/ast/#CaseClause || godoc go/ast CaseClause
	//  Case  token.Pos // position of "case" or "default" keyword
	//  List  []Expr    // list of expressions or types; nil means default case
	//  Colon token.Pos // position of ":"
	//  Body  []Stmt    // statement list; or nil
	case *ast.CaseClause:
		// To check the last statements
		tr.wasReturn = false
		tr.wasFallthrough = false

		tr.iCase++
		tr.addLine(typ.Case)

		if typ.List != nil {
			for i, expr := range typ.List {
				if i != 0 {
					tr.WriteString(SP)
				}
				tr.WriteString(fmt.Sprintf("case %s:", tr.getExpression(expr).String()))
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

		if !tr.wasFallthrough && !tr.wasReturn && tr.iCase != tr.lenCase {
			tr.WriteString(SP + "break;")
		}

	// http://golang.org/pkg/go/ast/#DeclStmt || godoc go/ast DeclStmt
	//  Decl Decl
	case *ast.DeclStmt:
		switch decl := typ.Decl.(type) {
		case *ast.GenDecl:
			switch decl.Tok {
			case token.VAR:
				tr.getVar(decl.Specs, false)
			case token.CONST:
				tr.getConst(decl.Specs, false)
			case token.TYPE:
				tr.getType(decl.Specs, false)
			default:
				panic("unreachable")
			}
		default:
			panic(fmt.Sprintf("unimplemented: %T", decl))
		}

	// http://golang.org/pkg/go/ast/#ExprStmt || godoc go/ast ExprStmt
	//  X Expr // expression
	case *ast.ExprStmt:
		tr.WriteString(tr.getExpression(typ.X).String() + ";")

	// http://golang.org/doc/go_spec.html#For_statements
	// https://developer.mozilla.org/en/JavaScript/Reference/Statements/for
	//
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

		if typ.Cond != nil {
			tr.WriteString(SP)
			tr.WriteString(tr.getExpression(typ.Cond).String())
		}
		tr.WriteString(";")

		if typ.Post != nil {
			tr.WriteString(SP)
			tr.getStatement(typ.Post)
		}

		tr.WriteString(")" + SP)
		tr.getStatement(typ.Body)

	// http://golang.org/doc/go_spec.html#Go_statements
	//
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

		tr.WriteString(fmt.Sprintf("if%s(%s)%s", SP, tr.getExpression(typ.Cond).String(), SP))
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
		tr.WriteString(tr.getExpression(typ.X).String() + typ.Tok.String())

	// http://golang.org/doc/go_spec.html#For_statements
	// https://developer.mozilla.org/en/JavaScript/Reference/Statements/for...in
	//
	// http://golang.org/pkg/go/ast/#RangeStmt || godoc go/ast RangeStmt
	//  For        token.Pos   // position of "for" keyword
	//  Key, Value Expr        // Value may be nil
	//  TokPos     token.Pos   // position of Tok
	//  Tok        token.Token // ASSIGN, DEFINE
	//  X          Expr        // value to range over
	//  Body       *BlockStmt
	case *ast.RangeStmt:
		expr := tr.getExpression(typ.X).String()
		key := tr.getExpression(typ.Key).String()
		value := ""

		if typ.Value != nil {
			value = tr.getExpression(typ.Value).String()

			if typ.Tok == token.DEFINE {
				tr.WriteString(fmt.Sprintf("var %s;%s", value, SP))
			}
		}

		tr.WriteString(fmt.Sprintf("for%s(%s in %s)%s", SP, key, expr, SP))

		if typ.Value != nil {
			tr.WriteString(fmt.Sprintf("{%s=%s[%s];", SP+value+SP, SP+expr, key))
			tr.skipLbrace = true
		}

		tr.getStatement(typ.Body)

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

		// Multiple values
		if len(typ.Results) != 1 {
			results := ""
			for i, v := range typ.Results {
				if i != 0 {
					results += "," + SP
				}
				results += tr.getExpression(v).String()
			}

			tr.WriteString("return [" + results + "];")
		} else {
			tr.WriteString("return " + tr.getExpression(typ.Results[0]).String() + ";")
		}

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
		tr.isSwitch = true
		tr.lenCase = len(typ.Body.List)
		tr.iCase = 0

		if typ.Init != nil {
			tr.getStatement(typ.Init) // use isSwitch
			tr.WriteString(SP)
		}

		if typ.Tag != nil {
			tag = tr.getExpression(typ.Tag).String()
		} else if tr.switchInit != "" {
			tag = tr.switchInit
			tr.switchInit = ""
		} else {
			tag = "true"
		}

		tr.WriteString(fmt.Sprintf("switch%s(%s)%s", SP, tag, SP))
		tr.isSwitch = false
		tr.getStatement(typ.Body)

	// === Not supported

	// http://golang.org/doc/go_spec.html#Defer_statements
	//
	// http://golang.org/pkg/go/ast/#DeferStmt || godoc go/ast DeferStmt
	//  Defer token.Pos // position of "defer" keyword
	//  Call  *CallExpr
	case *ast.DeferStmt:
		tr.addError("%s: defer statement", tr.fset.Position(typ.Defer))

	// http://golang.org/doc/go_spec.html#Labeled_statements
	// https://developer.mozilla.org/en/JavaScript/Reference/Statements/label
	//
	// http://golang.org/pkg/go/ast/#LabeledStmt || godoc go/ast LabeledStmt
	//  Label *Ident
	//  Colon token.Pos // position of ":"
	//  Stmt  Stmt
	case *ast.LabeledStmt:
		tr.addError("%s: use of label", tr.fset.Position(typ.Pos()))

	default:
		panic(fmt.Sprintf("unimplemented: %T", stmt))
	}
}
