// Copyright 2011  The "GoJscript" Authors
//
// Use of this source code is governed by the BSD 2-Clause License
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

package gojs

import (
	"fmt"
	"regexp"
	"strings"
)

/*
## Pointers

To identify variables that could be addressed ahead, it is used the map:

	{ number of function: {number of block: variable name} }

In the generated code, it is added a tag before and after of each new variable
but pointer. The tag uses the schema `{{side:funcId:blockId:varName}}`

	*side:* *L* or *R* if the tag is on the left or on the right of the variable
	*funcId:* identifier of function. '0' is for global declarations
	*blockId:* number of block inner of that function. Start with '1'
	*varName:* variable's name

so, the variables addressed can be boxed (placed between brackets).

It is also added the tag `{{A:funcId:blockId:varName}}` after of the variable
name in the assignment of variables.
*/

// To remove tags related to pointers
var reTagPointer = regexp.MustCompile(`<<[LRA]:\d+:\d+:[^>]+>>`)

// Returns a tag to identify pointers.
func tagPointer(typ rune, funcId, blockId int, name string) string {
	if typ != 'L' && typ != 'R' && typ != 'A' {
		panic("invalid identifier for pointer: " + string(typ))
	}

	return fmt.Sprintf("<<%s:%d:%d:%s>>", string(typ), funcId, blockId, name)
}

// Checks if a variable name is in the list of pointers by addressing
func (tr *transform) isPointer(str string) bool {
	// Check from the last block until the first one.
	for i := tr.blockLevel; i >= 0; i-- {
		for _, name := range tr.addressed[tr.funcLevel][i] {
			if name == str { // It is already marked
				return true
			}
		}
	}
	return false
}

// Appends a variable name to the list of pointers.
func (tr *transform) addPointer(str string) {
	if tr.isPointer(str) {
		return
	}

	// Search the point where the variable was declared.
	for i := tr.blockLevel; i >= 0; i-- {
		for _, name := range tr.vars[tr.funcLevel][i] {
			if name == str {
				tr.addressed[tr.funcLevel][i] = append(tr.addressed[tr.funcLevel][i], name)
				return
			}
		}
	}

	// Finally, search in the global variables (funcId = 0).
	for i := tr.blockLevel; i >= 0; i-- {
		for _, name := range tr.vars[0][i] {
			if name == str {
				tr.addressed[0][i] = append(tr.addressed[0][i], name)
				return
			}
		}
	}

	panic("unreachable")
}

// Replaces brackets in tags for variables addressed.
func (tr *transform) replaceBrackets(str *string) {
//println("Pointers")
	for funcId, v := range tr.addressed {
//fmt.Println(funcId, v)
		for blockId, vars := range v {
			for _, varName := range vars {
				lBrack := tagPointer('L', funcId, blockId, varName)
				rBrack := tagPointer('R', funcId, blockId, varName)

				*str = strings.Replace(*str, lBrack, "[", 1)
				*str = strings.Replace(*str, rBrack, "]", 1)
			}
		}
	}
}

// Replaces '[0]' in tags for assignment in variables addressed.
func (tr *transform) replaceAssign(str *string) {
//println("Assignment")
	for funcId, v := range tr.assigned {
//fmt.Println(funcId, v)
		for blockId, vars := range v {
			for _, varName := range vars {
				assigned := tagPointer('A', funcId, blockId, varName)
				*str = strings.Replace(*str, assigned, "[0]", 1)
			}
		}
	}
}
