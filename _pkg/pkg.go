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

// Handle Go features.

package g

// Adds public names to the map named in "pkg".
func Export(pkg map[interface{}]interface{}, exported []interface{}) {
	for _, v := range exported {
		pkg.v = v
	}
}

// == Map
//

// Represents a map M as an object with two fields, "map" and "zero".
// The compiler put the appropriate zero value for the map (which it is work out
// from the map type).
type M struct {
	m map[interface{}]interface{} // map
	z interface{}                 // zero value for the map
}

// Gets the value for the key "k".
// If looking some key up in M's map gets you "nil" ("undefined" in JS),
// then return a copy of the zero value.
func (m M) get(k interface{}) (interface{}, bool) {
	v := m.m

	// Allow multi-dimensional index (separated by commas)
	for i := 0; i < len(arguments); i++ {
		v = v[arguments[i]]
	}

	if v == nil {
		return m.z, false
	}
	return v, true
}
