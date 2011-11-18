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

import "testing"

var files = []string{"const.go"}

func TestSintaxis(t *testing.T) {
	for _, v := range files {
		if err := checkSintaxis(file(v)); err != nil {
			t.Fatalf("expected correct sintaxis, got\n%s", err)
		}
	}
}

func TestConst(t *testing.T) {
	if err := Compile(file("const.go")); err != nil {
		t.Fatalf("expected scanning file, got\n%s", err)
	}
}

// * * *

func file(name string) string {
	return "../test/"+name
}
