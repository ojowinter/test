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

import "testing"

const (
	DIR_PKG  = "../_pkg/"
	DIR_TEST = "../_test/"
)

func init() {
	MaxMessage = 100 // to show all errors

	// In the test files, it is used functions "fmt.Print" to show errors.
	for _, v := range []string{"", "ln", "f"} {
		Function["fmt.Print" + v] = "alert"
	}
}

func TestConst(t *testing.T)  { compile('t', "decl_const.go", t) }
func TestVar(t *testing.T)    { compile('t', "decl_var.go", t) }
func TestStruct(t *testing.T) { compile('t', "decl_struct.go", t) }
//func TestOp(t *testing.T)     { compile('t', "operator.go", t) }

func TestPointer(t *testing.T) { compile('t', "pointer.go", t) }
func TestFunc(t *testing.T)    { compile('t', "func.go", t) }
func TestCompo(t *testing.T)   { compile('t', "composite.go", t) }
func TestSlice(t *testing.T)   { compile('t', "slice.go", t) }
func TestMap(t *testing.T)     { compile('t', "map.go", t) }
func TestMethod(t *testing.T)  { compile('t', "method.go", t) }

// == Warnings
//
// ../test/control.go:44:2: 'default' clause above 'case' clause in switch statement
func Example_control() { Compile(DIR_TEST + "control.go") }

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
func Example_decl() { Compile(DIR_TEST + "error_decl.go") }

// == Errors
//
// ../test/error_stmt.go:6:13: channel type
// ../test/error_stmt.go:8:2: goroutine
// ../test/error_stmt.go:9:2: defer directive
// ../test/error_stmt.go:12:2: built-in function recover()
// ../test/error_stmt.go:18:1: use of label
// ../test/error_stmt.go:23:3: goto directive
func Example_stmt() { Compile(DIR_TEST + "error_stmt.go") }

// === Library
//

func TestPkg(t *testing.T) { compile('p', "pkg.go", t) }

// * * *

func compile(kind rune, filename string, t *testing.T) {
	dir := ""

	if kind == 't' {
		dir = DIR_TEST
	} else if kind == 'p' {
		dir = DIR_PKG
	} else {
		panic("Wrong kind")
	}

	if err := Compile(dir + filename); err != nil {
		t.Fatalf("expected parse file: %s", err)
	}
}
