// Copyright 2011  The "GoScript" Authors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package goscript

import (
	//"bytes"
	"fmt"
	"go/ast"
	"go/token"
)

/*// Represents a statement.
type statement struct {
	*bytes.Buffer // sintaxis translated
}

// Initializes a new type of "statement".
func newStatement() *statement {
	return &statement{
		new(bytes.Buffer),
	}
}

// Returns the Go statement in JavaScript.
func getStatement(stmt ast.Stmt) string {
	s := newStatement()

	s.transform(stmt)
	return s.String()
}

// Transforms the Go statement.
func (s *statement) transform(stmt ast.Stmt) {*/
func (tr *transform) getStatement(stmt ast.Stmt) {
	switch typ := stmt.(type) {

	// http://golang.org/pkg/go/ast/#AssignStmt || godoc go/ast AssignStmt
	//  Lhs    []Expr
	//  TokPos token.Pos   // position of Tok
	//  Tok    token.Token // assignment token, DEFINE
	//  Rhs    []Expr
	case *ast.AssignStmt:
		var isNew bool

		switch typ.Tok {
		case token.DEFINE:
			isNew = true
		case token.ASSIGN:
		default:
			panic(fmt.Sprintf("token unimplemented: %T", typ.Tok))
		}

		if isNew {
			tr.WriteString("var ")
		}

		isFirst := true
		for i, v := range typ.Lhs {
			lIdent := getExpression("", v)

			if lIdent == "_" {
				continue
			}

			rIdent := getExpression("", typ.Rhs[i])

			if isFirst {
				isFirst = false
			} else {
				tr.WriteString("," + SP)
			}

			tr.WriteString(lIdent + SP + "=" + SP + rIdent)
		}
		tr.WriteString(";")

	// http://golang.org/pkg/go/ast/#BlockStmt || godoc go/ast BlockStmt
	//  Lbrace token.Pos // position of "{"
	//  List   []Stmt
	//  Rbrace token.Pos // position of "}"
	case *ast.BlockStmt:
		tr.WriteString("{")

		for _, v := range typ.List {
			tr.addLine(v.Pos())
			tr.WriteString(TAB)
			tr.getStatement(v)
		}

		tr.addLine(typ.Rbrace)
		tr.WriteString("}")

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
		tr.addLine(typ.If)

		if typ.Init != nil {
			tr.getStatement(typ.Init)
			tr.WriteString(SP)
		}

		tr.WriteString(fmt.Sprintf("if%s(%s)%s", SP, getExpression("", typ.Cond), SP))
		tr.getStatement(typ.Body)

		if typ.Else != nil {
			tr.WriteString(SP + "else ")
			tr.getStatement(typ.Else)
		}

	// http://golang.org/doc/go_spec.html#Return_statements
	// https://developer.mozilla.org/en/JavaScript/Reference/Statements/return
	//
	// http://golang.org/pkg/go/ast/#ReturnStmt || godoc go/ast ReturnStmt
	//  Return  token.Pos // position of "return" keyword
	//  Results []Expr    // result expressions; or nil
	case *ast.ReturnStmt:
		ret := "return"

		if typ.Results == nil {
			tr.WriteString(ret + ";")
			break
		}
		if len(typ.Results) != 1 {
			tr.addError("%s: return multiple values", tr.fset.Position(typ.Return))
			break
		}

		tr.WriteString(ret + " " + getExpression("", typ.Results[0]) + ";")

	default:
		panic(fmt.Sprintf("unimplemented: %T", stmt))
	}
}
