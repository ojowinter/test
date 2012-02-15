// Copyright 2012  The "GoScript" Authors
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
	"go/ast"
	"strings"
)

// Checks if a variable name is a map.
func (tr *transform) isMap(name string) bool {
	if name == "" {
		return false
	}
	name = strings.SplitN(name, "<<", 2)[0] // could have a tag

	for funcId := tr.funcId; funcId >= 0; funcId-- {
		for blockId := tr.blockId; blockId >= 0; blockId-- {
			if _, ok := tr.vars[funcId][blockId][name]; ok { // variable found
				if _, okm := tr.maps[funcId][blockId][name]; okm { // map checking
					return true
				}
				return false
			}
		}
	}
	return false
}

// Returns the zero value of a map.
func (tr *transform) zeroOfMap(m *ast.MapType) string {
	if mapT, ok := m.Value.(*ast.MapType); ok { // nested map
		return tr.zeroOfMap(mapT)
	}
	v, _ := tr.zeroValue(true, m.Value)
	return v
}
