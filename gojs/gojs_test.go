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

func TestConst(t *testing.T) { compile("const.go", t) }
func TestVar(t *testing.T)   { compile("var.go", t) }
func TestType(t *testing.T)  { compile("type.go", t) }
func TestFunc(t *testing.T)  { compile("func.go", t) }
//func TestOp(t *testing.T)    { compile("operator.go", t) }

// == Warnings
//
// ../test/control.go:82:2: 'default' clause above 'case' clause in switch statement
func Example_control() { Compile(DIR + "control.go") }

// == Errors
//
// os: import from core library
// ../test/error_decl.go:13:10: complex128 type
// ../test/error_decl.go:14:10: complex128 type
// ../test/error_decl.go:15:10: complex128 type
// ../test/error_decl.go:16:10: complex128 type
// ../test/error_decl.go:18:6: built-in function complex()
// ../test/error_decl.go:23:17: complex64 type
// ../test/error_decl.go:24:13: complex64 type
// ../test/error_decl.go:25:11: complex64 type
// ../test/error_decl.go:30:14: complex128 type
// ../test/error_decl.go:31:9: complex128 type
// ../test/error_decl.go:32:12: complex128 type
// ../test/error_decl.go:37:16: complex64 type
// ../test/error_decl.go:38:23: complex64 type
// ../test/error_decl.go:39:11: complex64 type
// ../test/error_decl.go:48:12: channel type
// ../test/error_decl.go:49:12: channel type
// ../test/error_decl.go:50:7: channel operator
// ../test/error_decl.go:60:2: function type in struct
// ../test/error_decl.go:64:4: int64 type
// ../test/error_decl.go:65:2: anonymous field in struct
// ../test/error_decl.go:66:4: complex128 type
func Example_decl() { Compile(DIR + "error_decl.go") }

// == Errors
//
// ../test/error_stmt.go:6:13: channel type
// ../test/error_stmt.go:8:2: goroutine
// ../test/error_stmt.go:9:2: defer directive
// ../test/error_stmt.go:12:2: built-in function recover()
// ../test/error_stmt.go:18:1: use of label
// ../test/error_stmt.go:23:3: goto directive
func Example_stmt() { Compile(DIR + "error_stmt.go") }

// === Helpers
//

//func TestHelper(t *testing.T) { compile("../helper/helper.go", t) }

// * * *

func compile(filename string, t *testing.T) {
	if err := Compile(DIR + filename); err != nil {
		t.Fatal("expected parse file")
	}
}
