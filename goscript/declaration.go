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
	"go/ast"
)

// Maximum number of expressions to get.
// The expressions are the values after of "=".
const MAX_EXPRESSION = 10

// Constants
//
// http://golang.org/doc/go_spec.html#Constant_declarations
// https://developer.mozilla.org/en/JavaScript/Reference/Statements/const
func (tr *transform) getConst(spec []ast.Spec) {
	iotas := make([]int, MAX_EXPRESSION)
	lastValues := make([]string, MAX_EXPRESSION)

	// http://golang.org/pkg/go/ast/#ValueSpec || godoc go/ast ValueSpec
	//  Names   []*Ident      // value names (len(Names) > 0)
	//  Type    Expr          // value type; or nil
	//  Values  []Expr        // initial values; or nil
	for _, s := range spec {
		vSpec := s.(*ast.ValueSpec)

		if len(vSpec.Values) > MAX_EXPRESSION {
			panic("length of 'iotas' is lesser than 'vSpec.Values'")
		}

		// Checking
		if err := newCheck(tr.fset).Type(vSpec.Type); err != nil {
			tr.err = append(tr.err, err)
			continue
		}

		names, skipName := getName(vSpec)

		// === Values
		// http://golang.org/pkg/go/ast/#Expr || godoc go/ast Expr
		//  type Expr interface
		values := make([]string, 0)

		if len(vSpec.Values) != 0 {
			for i, v := range vSpec.Values {
				var expr string

				// Checking
				if err := newCheck(tr.fset).Type(v); err != nil {
					tr.err = append(tr.err, err)
					continue
				}

				val := newValue(names[i])
				val.getValue(v)

				if val.useIota {
					expr = fmt.Sprintf(val.String(), iotas[i])
					iotas[i]++
				} else {
					expr = val.String()
				}

				if !skipName[i] {
					values = append(values, expr)
					lastValues[i] = val.String()
				}
			}
		} else { // get last value of iota
			for i := 0; i < len(names); i++ {
				expr := fmt.Sprintf(lastValues[i], iotas[0])
				values = append(values, expr)
			}
			iotas[0]++
		}

		// Skip write buffer, if any error
		if tr.err != nil {
			continue
		}

		// === Write
		// TODO: calculate expression using "exp/types"
		isFirst := true
		for i, v := range names {
			if skipName[i] {
				continue
			}

			if isFirst {
				isFirst = false
				tr.bConst.WriteString(fmt.Sprintf("const %s = %s", v, values[i]))
			} else {
				tr.bConst.WriteString(fmt.Sprintf(", %s = %s", v, values[i]))
			}
		}

		// It is possible that there is only a blank identifier
		if isFirst {
			continue
		}

		tr.bConst.WriteString(";\n")
	}
}

// Variables
//
// http://golang.org/doc/go_spec.html#Variable_declarations
// https://developer.mozilla.org/en/JavaScript/Reference/Statements/var
// https://developer.mozilla.org/en/JavaScript/Reference/Statements/let
//
// TODO: use let for local variables
func (tr *transform) getVar(spec []ast.Spec) {
	// http://golang.org/pkg/go/ast/#ValueSpec || godoc go/ast ValueSpec
	for _, s := range spec {
		vSpec := s.(*ast.ValueSpec)

		// Checking
		if err := newCheck(tr.fset).Type(vSpec.Type); err != nil {
			tr.err = append(tr.err, err)
			continue
		}

		names, skipName := getName(vSpec)

		// === Values
		// http://golang.org/pkg/go/ast/#Expr || godoc go/ast Expr
		values := make([]string, 0)

		for i, v := range vSpec.Values {
			// Checking
			if err := newCheck(tr.fset).Type(v); err != nil {
				tr.err = append(tr.err, err)
				continue
			}

			// Skip when it is not a function
			if skipName[i] {
				if _, ok := v.(*ast.CallExpr); !ok {
					continue
				}
			}

			val := newValue(names[i])
			val.getValue(v)

			if !skipName[i] {
				values = append(values, val.String())
			}
		}

		if tr.err != nil {
			continue
		}

		// === Write
		// TODO: calculate expression using "exp/types"
		isFirst := true
		for i, n := range names {
			if skipName[i] {
				continue
			}

			if isFirst {
				isFirst = false
				tr.bVar.WriteString("var " + n)
			} else {
				tr.bVar.WriteString(", " + n)
			}

			if len(values) != 0 {
				tr.bVar.WriteString(" = " + values[i])
			}
		}

		last := tr.bVar.Bytes()[tr.bVar.Len()-1] // last character

		if last != '}' && last != ';' {
			tr.bVar.WriteString(";\n")
		} else {
			tr.bVar.WriteString("\n")
		}
	}
}

// Types
//
// http://golang.org/doc/go_spec.html#Type_declarations
func (tr *transform) getType(spec []ast.Spec) {
	// Format fields
	format := func(fields []string) (args, allFields string) {
		for i, f := range fields {
			if i == 0 {
				args = f
			} else {
				args += ", " + f
			}

			allFields += fmt.Sprintf("\n\tthis.%s = %s;", f, f)
		}
		return
	}

	// http://golang.org/pkg/go/ast/#TypeSpec || godoc go/ast TypeSpec
	//  Name    *Ident        // type name
	//  Type    Expr          // *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
	for _, s := range spec {
		tSpec := s.(*ast.TypeSpec)
		fields := make([]string, 0) // names of fields
		//!anonField := make([]bool, 0) // anonymous field

		// Checking
		if err := newCheck(tr.fset).Type(tSpec.Type); err != nil {
			tr.err = append(tr.err, err)
			continue
		}

		switch typ := tSpec.Type.(type) {
		default:
			panic(fmt.Sprintf("[getType] unimplemented: %T", typ))

		case *ast.Ident:

		// http://golang.org/pkg/go/ast/#StructType || godoc go/ast StructType
		//  Struct     token.Pos  // position of "struct" keyword
		//  Fields     *FieldList // list of field declarations
		//  Incomplete bool       // true if (source) fields are missing in the Fields list
		case *ast.StructType:
			if typ.Incomplete {
				panic("[getType:StructType] list of fields incomplete ???")
			}

			// http://golang.org/pkg/go/ast/#FieldList || godoc go/ast FieldList
			//  List    []*Field  // field list; or nil
			for _, field := range typ.Fields.List {
				if _, ok := field.Type.(*ast.FuncType); ok {
					tr.err = append(tr.err,
						fmt.Errorf("%s: function type in struct",
							tr.fset.Position(field.Pos())),
					)
					continue
				}

				// http://golang.org/pkg/go/ast/#Field || godoc go/ast Field
				//  Names   []*Ident      // field/method/parameter names; or nil if anonymous field
				//  Type    Expr          // field/method/parameter type
				//  Tag     *BasicLit     // field tag; or nil

				// Checking
				if err := newCheck(tr.fset).Type(field.Type); err != nil {
					tr.err = append(tr.err, err)
					continue
				}
				if field.Names == nil {
					tr.err = append(tr.err,
						fmt.Errorf("%s: anonymous field in struct",
							tr.fset.Position(field.Pos())),
					)
					continue
				}

				for _, n := range field.Names {
					name := n.Name

					if name == "_" {
						continue
					}

					fields = append(fields, name)
					//!anonField = append(anonField, false)
				}
			}
		}

		if tr.err != nil {
			continue
		}

		// === Write
		name := tSpec.Name.Name
		args, allFields := format(fields)

		tr.bType.WriteString(fmt.Sprintf("function %s(%s) {", name, args))

		if len(allFields) != 0 {
			tr.bType.WriteString(allFields)
			tr.bType.WriteString("\n}\n")
		} else {
			tr.bType.WriteString("}\n") //! empty struct
		}
	}
}

//
// === Utility

// Gets the identifiers.
//
// http://golang.org/pkg/go/ast/#Ident || godoc go/ast Ident
//  Name    string    // identifier name
func getName(spec *ast.ValueSpec) (names []string, skipName []bool) {
	skipName = make([]bool, len(spec.Names)) // for blank identifiers "_"

	for i, v := range spec.Names {
		if v.Name == "_" {
			skipName[i] = true
			continue
		}
		names = append(names, v.Name)
	}

	return
}
