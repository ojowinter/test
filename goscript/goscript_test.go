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
	"testing"
)

func TestConst(t *testing.T)  { compile("const.go", t) }
func TestVar(t *testing.T)    { compile("var.go", t) }
func TestStruct(t *testing.T) { compile("struct.go", t) }

func TestError(t *testing.T) {
	err := Compile("../test/error.go")
	if err == nil {
		t.Fatal("expected error")
	}

	fmt.Println(err.Error())
}

// * * *

func compile(filename string, t *testing.T) {
	if err := Compile("../test/" + filename); err != nil {
		t.Fatalf("expected parse file, got\n%s", err)
	}
}
