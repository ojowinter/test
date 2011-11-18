// Copyright 2011  The "GotoJS" Authors
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

package tojs

import (
	"bufio"
	_ "fmt"
	"go/parser"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"os"
)

// Checks that the Go sintaxis is correct.
func checkSintaxis(filename string) error {
	_, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	return nil
}

// Compiles a Go source file to JavaScript.
func Compile(filename string) error {
	var s scanner.Scanner

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	bufConst := bufio.NewWriter(os.Stdout)
	//bufVar := bufio.NewWriter(os.Stdout)

	fset := token.NewFileSet()
	fileSet := fset.AddFile(filename, fset.Base(), len(file)) // register file
	s.Init(fileSet, file, nil, scanner.InsertSemis)

	//_, tok, lit := s.Scan()
	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}

		switch tok {
		case token.PACKAGE:
			_, _, lit = s.Scan()
			println(lit)

		case token.CONST:
			_const(bufConst, s)
			break
		}
	}

	bufConst.Flush()
	return nil
}

// Constants
//
// http://golang.org/doc/go_spec.html#Constant_declarations
// https://developer.mozilla.org/en/JavaScript/Reference/Statements/const
func _const(buf *bufio.Writer, s scanner.Scanner) {
	isNew := true // to know when write the keyword "const"

	for {
		_, tok, lit := s.Scan()
		if tok.IsKeyword() || tok == token.RPAREN {
			break
		}

		// There are multiple statements
		if tok == token.LPAREN { // (
			continue
		}

		if tok.IsLiteral() {
			if isType(tok, lit) {
				continue
			}

			if isNew && tok == token.IDENT {
				buf.WriteString("const ")
				isNew = false
			}

			buf.WriteString(lit)
			continue
		}

		// === IsOperator
		switch tok {
		case token.COMMA:
			buf.WriteString(lit + " ")
		case token.ASSIGN:
			buf.WriteString(" " + lit + " ")
		case token.SEMICOLON: // End of statement
			buf.WriteString(";\n")
			isNew = true
		default:
			buf.WriteString(lit)
		}

	}
}
