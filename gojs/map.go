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

type opMap uint8

const (
	_ opMap = iota
	MAP_KEYS
	MAP_ZERO
)

// Finds if a variable name is a map.
func (tr *transform) findMap(name string) bool {
	for funcId := tr.funcId; funcId >= 0; funcId-- {
		for blockId := tr.blockId; blockId >= 0; blockId-- {
			if _, ok := tr.vars[funcId][blockId][name]; ok { // variable found
				if _, okm := tr.mapKeys[funcId][blockId][name]; okm { // map checking
					return true
				}
				return false
			}
		}
	}
	return false
}
/*
func (tr *transform) 

	if _, ok := tr.mapKeys[tr.funcId][tr.blockId][name]; !ok {
		tr.mapKeys[tr.funcId][tr.blockId][name] = make(map[string]struct{})
	}


	if _, ok := e.tr.mapKeys[e.tr.funcId][e.tr.blockId][e.tr.lastVarName]; !ok {
		e.tr.mapKeys[e.tr.funcId][e.tr.blockId][e.tr.lastVarName] = make(map[string]struct{})
	} else {
*/
