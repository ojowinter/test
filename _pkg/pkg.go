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

//	base *interface{}  // element type
// S represents a slice.
type S struct {
	f   []interface{} // slice
	len int
	cap int
}

// Sets the slice from an array.
func (s S) fromArray(i interface{}, low, high int) {
	s.f = i.slice(low, high)
	s.len = high - low
	s.cap = len(i) - s.len
}

// Sets the slice from another slice.
func (s S) fromSlice(i []interface{}, low, high int) {
	s.f = i.f.slice(low, high)
	s.len = high - low
	s.cap = len(i.f) - s.len
}

// Initializes the slice.
func (s S) make(len, cap int) {
	s.f = nil
	s.len = len

	if cap != nil {
		s.cap = cap
	} else {
		s.cap = len
	}
}

// Appends an element to the slice.
func (s S) append(elt interface{}) {
	if s.len == s.cap {
		s.cap = s.len * 2
	}
	s.len++
}

// Returns the slice like a string.
func (s S) String() string {
	str := s.f.toString()
	spl := str.split(",")
	return spl.join("")

//	str := strings.Split(string(s.f.toString), ",")
//	return strings.Join(str, "")
}

/*func (s S) get() []interface{} {
	
}*/

/*// Returns the length of the slice.
func (s S) len() int {
	if s.f != nil {
		return len(s.f)
	}
	return 0
}

// Returns the capacity of the slice.
func (s S) cap() int {
	if s.f != nil {
		return s.base.p.length - s.f.length
	}
	return 0
}*/

// == Map
//

// M represents a map.
// The compiler adds the appropriate zero value for the map (which it is work out
// from the map type).
type M struct {
	f    map[interface{}]interface{} // map
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
