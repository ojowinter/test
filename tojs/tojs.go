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

package tojs

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	//"io/ioutil"
	"os"
	//"path"
	//"strings"
)

// Compiles a Go source file into JavaScript.
// Writes the output in "filename" but with extension ".js".
func Compile(filename string) error {
	var hasError bool

	/* Parse several files
	fset := token.NewFileSet()
	parse.ParseFile(fset, "a.go", nil, 0)
	parse.ParseFile(fset, "b.go", nil, 0)
	*/

	// If Go sintaxis is incorrect then there will be an error.
	node, err := parser.ParseFile(token.NewFileSet(), filename, nil, 0) //parser.ParseComments)
	if err != nil {
		return err
	}

	// === To get the position
	/*file, _ := os.Open(filename)
	info, err := file.Stat()
	if err != nil {
		return err
	}
	defer file.Close()

	fset := token.NewFileSet()
	nodeFile := fset.AddFile(filename, fset.Base(), int(info.Size))
	//fmt.Println(nodeFile.Position(97).String())*/

	// === Buffers
	bufConst := new(bytes.Buffer)
	bufType := new(bytes.Buffer)
	bufVar := new(bytes.Buffer)

	bufConst.WriteString("// Generated by GoScript <github.com/kless/GoScript>\n")

	for _, decl := range node.Decls {
		//fmt.Printf("%T : %v\n", decl,decl) // TODO: to delete when it's finished

		switch decl.(type) {

		// http://golang.org/pkg/go/ast/#GenDecl || godoc go/ast GenDecl
		//  Tok    token.Token   // IMPORT, CONST, TYPE, VAR
		//  Specs  []Spec
		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)

			switch genDecl.Tok {
			case token.CONST:
				getConst(bufConst, genDecl.Specs)
			case token.TYPE:
				getType(bufType, genDecl.Specs)
			case token.VAR:
				if err := getVar(bufVar, genDecl.Specs); err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err.Error())
					hasError = true
				}
			}
		}
	}

	if hasError {
		return errors.New("error: not supported in JavaScript")
	}

	// Write all buffers in the first one
	bufConst.Write(bufType.Bytes())
	bufConst.Write(bufVar.Bytes())
	fmt.Print(bufConst.String()) // TODO: delete

	//jsFile := strings.Replace(filename, path.Ext(filename), ".js", 1)
	//return ioutil.WriteFile(jsFile, bufConst.Bytes(), 0664)
	return nil
}
