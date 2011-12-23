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
)

// Valid functions to be transformed since they have a similar function in JS.
var ImportAndFunc = map[string][]string {
	"fmt": []string{"Print", "Println"},
}

// Similar functions in JavaScript.
var transformFunc = map[string]string {
	"print":       "alert",
	"println":     "alert",
	"fmt.Print":   "alert",
	"fmt.Println": "alert",
}

// Checks if the function can be transformed.
func isValidFunc(importPath, funcName *ast.Ident) bool {
	for _, f := range ImportAndFunc[importPath.Name] {
		if f == funcName.Name {
			return true
		}
	}
	return false
}

// Returns the equivalent function in JavaScript.
func getFuncJS(importPath, funcName *ast.Ident, args []ast.Expr) (string, error) {
	if !isValidFunc(importPath, funcName) {
		return "", fmt.Errorf("%s.%s: function from core library", importPath, funcName)
	}

	jsFunc := transformFunc[importPath.Name + "." + funcName.Name]
	arg1 := getExpression(args[0])

	switch funcName.Name {
	case "println", "Println":
		arg1 = arg1[:len(arg1)-1] + "\\n\""
	}

	return fmt.Sprintf("%s(%s);", jsFunc, arg1), nil
}
