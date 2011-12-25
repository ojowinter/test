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
	"go/ast"
	"strings"
)

// Valid functions to be transformed since they have a similar function in JS.
var ImportAndFunc = map[string][]string{
	"fmt": []string{"Print", "Println"},
}

// Similar functions in JavaScript.
var transformFunc = map[string]string{
	"print":       "alert",
	"println":     "alert",
	"fmt.Print":   "alert",
	"fmt.Println": "alert",
}

// Returns the equivalent function in JavaScript.
func (tr *transform) GetFuncJS(importName, funcName *ast.Ident, args []ast.Expr) (string, error) {
	var jsArgs, importStr string

	if importName != nil {
		importStr = importName.Name + "."

		if !isValidFunc(importName, funcName) {
			return "", fmt.Errorf("%s.%s: function from core library", importName, funcName)
		}
	}

	switch funcName.Name {
	case "print", "Print":
		jsArgs = tr.getPrintArgs(args, false)
	case "println", "Println":
		jsArgs = tr.getPrintArgs(args, true)
	}

	jsFunc := transformFunc[importStr+funcName.Name]
	return fmt.Sprintf("%s(%s);", jsFunc, jsArgs), nil
}

// Returns arguments to print.
func (tr *transform) getPrintArgs(args []ast.Expr, addLine bool) string {
	var jsArgs string
	lenArgs := len(args) - 1

	// Appends a character.
	add := func(s, char string) string {
		if strings.HasSuffix(s, "\"") {
			s = s[:len(s)-1] + char + "\""
		} else {
			s += SP + "+" + SP + "\"" + char + "\""
		}
		return s
	}

	for i, v := range args {
		expr := tr.getExpression(v)

		if i != 0 {
			jsArgs += SP + "+" + SP + expr
		} else {
			jsArgs = expr
		}

		if addLine {
			if i == lenArgs {
				jsArgs = add(jsArgs, "\\n")
			} else {
				jsArgs = add(jsArgs, " ")
			}
		}
	}

	return jsArgs
}

// * * *

// Checks if the function can be transformed.
func isValidFunc(importName, funcName *ast.Ident) bool {
	for _, f := range ImportAndFunc[importName.Name] {
		if f == funcName.Name {
			return true
		}
	}
	return false
}
