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
	"fmt"
	"go/ast"
	"go/token"
)

func (tr *transform) getStatement(stmt ast.Stmt) {
	switch typ := stmt.(type) {
	default:
		panic(fmt.Sprintf("[getStatement] unimplemented: %T", stmt))

	// http://golang.org/pkg/go/ast/#AssignStmt || godoc go/ast AssignStmt
	//  Lhs    []Expr
	//  TokPos token.Pos   // position of Tok
	//  Tok    token.Token // assignment token, DEFINE
	//  Rhs    []Expr
	case *ast.AssignStmt:
		switch typ.Tok {
		case token.DEFINE, token.ASSIGN:
		default:
			panic(fmt.Sprintf("[getStatement:AssignStmt] unimplemented: %T", typ.Tok))
		}

		isFirst := true
		tr.dst.WriteString("var ")

		for i, v := range typ.Lhs {
			lIdent := v.(*ast.Ident).Name

			if lIdent == "_" {
				continue
			}

			rIdent := typ.Rhs[i].(*ast.BasicLit).Value

			if isFirst {
				isFirst = false
			} else {
				tr.dst.WriteString(", ")
			}

			tr.dst.WriteString(lIdent + "=" + rIdent)
		}
		tr.dst.WriteString(";")

	// http://golang.org/pkg/go/ast/#IfStmt || godoc go/ast IfStmt
	//  If   token.Pos // position of "if" keyword
	//  Init Stmt      // initialization statement; or nil
	//  Cond Expr      // condition
	//  Body *BlockStmt
	//  Else Stmt // else branch; or nil
	case *ast.IfStmt:
		

	// http://golang.org/pkg/go/ast/#ReturnStmt || godoc go/ast ReturnStmt
	//  Return  token.Pos // position of "return" keyword
	//  Results []Expr    // result expressions; or nil
	case *ast.ReturnStmt:
		
	}
}

//
// === Utility


