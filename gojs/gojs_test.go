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

func TestConst(t *testing.T) { compile("const.go", t) }
func TestVar(t *testing.T)   { compile("var.go", t) }
func TestType(t *testing.T)  { compile("type.go", t) }
func TestFunc(t *testing.T)  { compile("func.go", t) }
func TestOp(t *testing.T)  { compile("operator.go", t) }

func TestErrorDecl(t *testing.T) { compileErr("error_decl.go", t) }
func TestErrorStmt(t *testing.T) { compileErr("error_stmt.go", t) }

// * * *

func compile(filename string, t *testing.T) {
	if err := Compile("../test/" + filename); err != nil {
		t.Fatalf("expected parse file, got\n%s", err)
	}
}

func compileErr(filename string, t *testing.T) {
	err := Compile("../test/" + filename)
	if err == nil {
		t.Fatal("expected error")
	}
}
