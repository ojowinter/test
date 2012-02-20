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

// Adds public names from "exported" to the map "pkg".
func Export(pkg map[interface{}]interface{}, exported []interface{}) {
	for _, v := range exported {
		pkg.v = v
	}
}

// == Slice
//

// S represents a slice.
type S struct {
	f   []interface{} // slice field
	cap int           // capacity
}

// == Map
//

// M represents a map.
// The compiler adds the appropriate zero value for the map (which it is work out
// from the map type).
type M struct {
	f    map[interface{}]interface{} // map field
	zero interface{}                 // zero value for the map
}

// Gets the value for the key "k".
// If looking some key up in M's map gets you "nil" ("undefined" in JS),
// then return a copy of the zero value.
func (m M) get(k interface{}) (interface{}, bool) {
	v := m.f

	// Allow multi-dimensional index (separated by commas)
	for i := 0; i < len(arguments); i++ {
		v = v[arguments[i]]
	}

	if v == nil {
		return m.zero, false
	}
	return v, true
}
