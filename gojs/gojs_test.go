// Copyright 2011  The "GoJscript" Authors
//
// Use of this source code is governed by the BSD 2-Clause License
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

package gojs

import "testing"

const DIR = "../test/"

func init() {
	MaxMessage = 100 // to show all errors
}

func TestConst(t *testing.T)   { compile("const.go", t) }
func TestVar(t *testing.T)     { compile("var.go", t) }
func TestType(t *testing.T)    { compile("type.go", t) }
func TestFunc(t *testing.T)    { compile("func.go", t) }
func TestControl(t *testing.T) { compile("control.go", t) }
//func TestOp(t *testing.T)      { compile("operator.go", t) }

// == Errors
//
// os: import from core library
// ../test/error_decl.go:13:10: complex128 type
// ../test/error_decl.go:14:10: complex128 type
// ../test/error_decl.go:15:10: complex128 type
// ../test/error_decl.go:16:10: complex128 type
//MORE ERRORS
func ExampleCompile_decl() { Compile(DIR + "error_decl.go") }

// == Errors
//
// ../test/error_stmt.go:6:13: channel type
// ../test/error_stmt.go:8:2: goroutine
// ../test/error_stmt.go:9:2: defer statement
// ../test/error_stmt.go:11:2: built-in function panic()
// ../test/error_stmt.go:12:2: built-in function recover()
// ../test/error_stmt.go:18:1: use of label
func ExampleCompile_stmt () { Compile(DIR + "error_stmt.go") }

// * * *

func compile(filename string, t *testing.T) {
	if err := Compile(DIR + filename); err != nil {
		t.Fatal("expected parse file")
	}
}
