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
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	HEADER = "/* Generated by GoScript <github.com/kless/GoScript> */"
	BLANK  = "_"  // blank identifier
	EMPTY  = `""` // empty string
)

const (
	// To be able to minimize code
	NL  = "<<NL>>" // new line
	SP  = "<<SP>>" // space
	TAB = "<<TAB>>"

	ADDR = "<<&>>" // to mark assignments to addresses
	IOTA = "<<iota>>"
	NIL  = "<<nil>>"
)

var void struct{} // A struct without any elements occupies no space at all.

var (
	Bootstrap  bool // is transforming the gojs's package?
	MaxMessage = 10 // maximum number of errors and warnings to show.
)

// Represents information about code being transformed to JavaScript.
type transform struct {
	line     int // actual line
	hasError bool

	fset          *token.FileSet
	*bytes.Buffer // sintaxis translated to JS
	*dataStmt     // extra data for a statement

	err      []error  // errors
	warn     []string // warnings
	exported []string // declarations to be exported

	//slice map[string]string // for range; key: function name, value: slice name
	//function string // actual function

	// == Variables defined in each block, for each function.
	// {Function Id: {Block id: {Name:
	vars   map[int]map[int]map[string]bool // is pointer?
	addr   map[int]map[int]map[string]bool // variable was assigned to an address?
	maps   map[int]map[int]map[string]struct{}
	slices map[int]map[int]map[string]struct{}

	// Zero value for custom types.
	zeroType map[int]map[int]map[string]string
}

func newTransform() *transform {
	tr := &transform{
		0,
		false,

		token.NewFileSet(),
		new(bytes.Buffer),
		&dataStmt{},

		make([]error, 0, MaxMessage),
		make([]string, 0, MaxMessage),
		make([]string, 0),

		//make(map[string]string),
		//"",

		make(map[int]map[int]map[string]bool),
		make(map[int]map[int]map[string]bool),
		make(map[int]map[int]map[string]struct{}),
		make(map[int]map[int]map[string]struct{}),
		make(map[int]map[int]map[string]string),
	}

	// == Global variables
	// Ones related to local variables are set in:
	// file func: *transform.getFunc()
	// file stmt: *transform.getStatement() (case: *ast.BlockStmt)

	// funcId = 0
	tr.vars[0] = make(map[int]map[string]bool)
	tr.addr[0] = make(map[int]map[string]bool)
	tr.maps[0] = make(map[int]map[string]struct{})
	tr.slices[0] = make(map[int]map[string]struct{})
	tr.zeroType[0] = make(map[int]map[string]string)

	// blockId = 0
	tr.vars[0][0] = make(map[string]bool)
	tr.addr[0][0] = make(map[string]bool)
	tr.maps[0][0] = make(map[string]struct{})
	tr.slices[0][0] = make(map[string]struct{})
	tr.zeroType[0][0] = make(map[string]string)

	return tr
}

// Returns the line number.
func (tr *transform) getLine(pos token.Pos) int {
	// -1 because it was inserted a line (the header)
	return tr.fset.Position(pos).Line - 1
}

// Appends new lines according to the position.
// Returns a boolean to indicate if have been added.
func (tr *transform) addLine(pos token.Pos) bool {
	var s string

	new := tr.getLine(pos)
	dif := new - tr.line

	if dif == 0 {
		return false
	}

	for i := 0; i < dif; i++ {
		s += NL
	}

	tr.WriteString(s)
	tr.line = new
	return true
}

// Appends an error.
func (tr *transform) addError(value interface{}, a ...interface{}) {
	if len(tr.err) == MaxMessage {
		return
	}

	switch typ := value.(type) {
	case string:
		tr.err = append(tr.err, fmt.Errorf(typ, a...))
	case error:
		tr.err = append(tr.err, typ)
	default:
		panic("wrong type")
	}

	if !tr.hasError {
		tr.hasError = true
	}
}

// Appends a warning message.
func (tr *transform) addWarning(format string, a ...interface{}) {
	if len(tr.warn) == MaxMessage {
		return
	}
	tr.warn = append(tr.warn, fmt.Sprintf(format, a...))
}

// Appends the declaration name if it is exported.
func (tr *transform) addIfExported(iName interface{}) {
	var name = ""

	switch typ := iName.(type) {
	case *ast.Ident:
		name = typ.Name
	case string:
		name = typ
	}

	if ast.IsExported(name) {
		tr.exported = append(tr.exported, name)
	}
}

// * * *

// Compiles a Go source file into JavaScript.
// Writes the output in "filename" but with extension ".js".
func Compile(filename string) error {
	trans := newTransform()
	pkgName := ""

	/* Parse several files
	parse.ParseFile(fset, "a.go", nil, 0)
	parse.ParseFile(fset, "b.go", nil, 0)
	*/

	// godoc go/ast File
	//  Doc        *CommentGroup   // associated documentation; or nil
	//  Package    token.Pos       // position of "package" keyword
	//  Name       *Ident          // package name
	//  Decls      []Decl          // top-level declarations; or nil
	//  Scope      *Scope          // package scope (this file only)
	//  Imports    []*ImportSpec   // imports in this file
	//  Unresolved []*Ident        // unresolved identifiers in this file
	//  Comments   []*CommentGroup // list of all comments in the source file

	node, err := parser.ParseFile(trans.fset, filename, nil, 0) //parser.ParseComments)
	if err != nil {
		return err
	}

	trans.WriteString(HEADER)

	// Package name
	pkgName = trans.getExpression(node.Name).String()

	if pkgName != "main" {
		trans.addLine(node.Package)
		trans.WriteString(fmt.Sprintf("var %s=%s{};%s(function()%s{",
			pkgName+SP, SP, SP, SP))
	}

	for _, decl := range node.Decls {
		switch decl.(type) {
		case *ast.FuncDecl:
			trans.getFunc(decl.(*ast.FuncDecl))

		// godoc go/ast GenDecl
		//  Tok    token.Token   // IMPORT, CONST, TYPE, VAR
		//  Specs  []Spec
		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)

			switch genDecl.Tok {
			case token.IMPORT:
				trans.getImport(genDecl.Specs)
			case token.CONST:
				trans.getConst(genDecl.TokPos, genDecl.Specs, true)
			case token.TYPE:
				trans.getType(genDecl.Specs, true)
			case token.VAR:
				trans.getVar(genDecl.Specs, true)
			}

		default:
			panic(fmt.Sprintf("unimplemented: %T", decl))
		}
	}

	// Any error?
	if trans.hasError {
		fmt.Fprintln(os.Stderr, " == Errors\n")

		for _, err := range trans.err {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		if len(trans.err) == MaxMessage {
			fmt.Fprintln(os.Stderr, "\n Too many errors")
		}

		return errors.New("") // to indicate that there was any error
	}

	// Export declarations in packages
	if pkgName != "main" {
		if len(trans.exported) != 0 {
			for i, v := range trans.exported {
				if i == 0 {
					trans.WriteString(NL+NL)
				}

				if !Bootstrap {
					if i == 0 {
						trans.WriteString(fmt.Sprintf("g.Export(%s,%s[%s",
							pkgName, SP, v))
					} else {
						trans.WriteString("," + SP + v)
					}
				} else {
					trans.WriteString(fmt.Sprintf("%s.%s=%s;%s",
						pkgName, v+SP, SP+v, NL))
				}
			}
			if !Bootstrap {
				trans.WriteString("]);")
			}
		} else {
			trans.WriteString(NL)
		}

		trans.WriteString(NL + "})();")
	}
	trans.WriteString(NL)

	// === Write
	baseFilename := strings.Replace(filename, path.Ext(filename), "", 1)
	str := trans.String()

	// Variables addressed
	trans.replacePointers(&str)

	// Version to debug
	deb := strings.Replace(str, NL, "\n", -1)
	deb = strings.Replace(deb, TAB, "\t", -1)
	deb = strings.Replace(deb, SP, " ", -1)

	if err := ioutil.WriteFile(baseFilename+".js", []byte(deb), 0664); err != nil {
		return err
	}
/*
	// Minimized version
	min := strings.Replace(str, NL, "", -1)
	min = strings.Replace(min, TAB, "", -1)
	min = strings.Replace(min, SP, "", -1)

	if err := ioutil.WriteFile(baseFilename + ".min.js", []byte(min), 0664); err != nil {
		return err
	}*/

	// Print warnings
	if len(trans.warn) != 0 {
		fmt.Fprintln(os.Stderr, " == Warnings\n")

		for _, v := range trans.warn {
			fmt.Fprintln(os.Stderr, v)
		}
		if len(trans.warn) == MaxMessage {
			fmt.Fprintln(os.Stderr, "\n Too many warnings")
		}
	}
/*
for k, v := range trans.slices {
	fmt.Println(k, v)
}
*/
	return nil
}
