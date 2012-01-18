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

	{number of function: {number of block: {variable name: is pointer?} }}

In the generated code, it is added a tag before and after of each new variable
but pointer. The tag uses the schema `{{side:funcId:blockId:varName}}`

	*side:* *L* or *R* if the tag is on the left or on the right of the variable
	*funcId:* identifier of function. '0' is for global declarations
	*blockId:* number of block inner of that function. Start with '1'
	*varName:* variable's name

so, the variables addressed can be boxed (placed between brackets).

It is also added the tag `{{P:funcId:blockId:varName}}` after of each variable
name.
*/

// To remove tags related to pointers
var reTagPointer = regexp.MustCompile(`<<[LRP]:\d+:\d+:[^>]+>>`)

// Returns a tag to identify pointers.
func tagPointer(typ rune, funcId, blockId int, name string) string {
	/*if typ != 'L' && typ != 'R' && typ != 'P' {
		panic("invalid identifier for pointer: " + string(typ))
	}*/

	return fmt.Sprintf("<<%s:%d:%d:%s>>", string(typ), funcId, blockId, name)
}

// Search the point where the variable was declared for tag it as pointer.
func (tr *transform) addPointer(name string) {
	// In the actual function
	if tr.funcId != 0 {
		for block := tr.blockId; block >= 1; block-- {
			if _, ok := tr.vars[tr.funcId][block][name]; ok {
				tr.vars[tr.funcId][block][name] = true
				return
			}
		}
	}

	// Finally, search in the global variables (funcId = 0).
	for block := tr.blockId; block >= 1; block-- {
		if _, ok := tr.vars[0][block][name]; ok {
			tr.vars[0][block][name] = true
			return
		}
	}
	//fmt.Printf("Function %d, block %d, name %s\n", tr.funcId, tr.blockId, name)
	panic("addPointer: unreachable")
}

// Replaces tags related to variables addressed.
func (tr *transform) replacePointers(str *string) {
	// Replaces tags in variables that access to pointers.
	toPointer := func(funcId, startBlock, endBlock int, varName string) {
		for block := startBlock; block <= endBlock; block++ {
			// Check if there is a variable named like the pointer in another block.
			if block != startBlock {
				if _, ok := tr.vars[funcId][block][varName]; ok {
					break
				}
			}
			pointer := tagPointer('P', funcId, block, varName)
			*str = strings.Replace(*str, pointer, "[0]", -1)
		}
	}

	for funcId, blocks := range tr.vars {
		for blockId := 1; blockId <= len(blocks); blockId++ {
			for name, isPointer := range tr.vars[funcId][blockId] {
				if isPointer {
					toPointer(funcId, blockId, len(blocks), name)

					// Replace brackets around variables addressed.
					lBrack := tagPointer('L', funcId, blockId, name)
					rBrack := tagPointer('R', funcId, blockId, name)

					*str = strings.Replace(*str, lBrack, "[", 1)
					*str = strings.Replace(*str, rBrack, "]", 1)
				}
			}
		}
	}
}
