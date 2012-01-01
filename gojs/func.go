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
	"regexp"
	"strings"
)

// Functions which can be transformed.
// The empty values are to indicate that the package (in the key) have any
// function to be transformed.
var Function = map[string]string{
	"print":        "alert",
	"println":      "alert",
	"fmt":          "",
	"fmt.Print":    "alert",
	"fmt.Println":  "alert",
	"fmt.Printf":   "alert",
	"math":         "",
	"math.Abs":     "Math.abs",
	"math.Acos":    "Math.acos",
	"math.Asin":    "Math.asin",
	"math.Atan":    "Math.atan",
	"math.Atan2":   "Math.atan2",
	"math.Ceil":    "Math.ceil",
	"math.Cos":     "Math.cos",
	"math.Exp":     "Math.exp",
	"math.Floor":   "Math.floor",
	"math.Log":     "Math.log",
	"math.Max":     "Math.max",
	"math.Min":     "Math.min",
	"math.Pow":     "Math.pow",
	"rand.Float32": "Math.random",
	"rand.Float64": "Math.random",
	//"math.":        "Math.round", // https://developer.mozilla.org/en/JavaScript/Reference/Global_Objects/Math/round
	"math.Sin":  "Math.sin",
	"math.Sqrt": "Math.sqrt",
	"math.Tan":  "Math.tan",
}

// Constants to transform.
var Constant = map[string]string{
	"math.E":      "Math.E",
	"math.Ln2":    "Math.LN2",
	"math.Log2E":  "Math.LOG2E",
	"math.Ln10":   "Math.LN10",
	"math.Log10E": "Math.LOG10E",
	"math.Pi":     "Math.PI",
	"math.Sqrt2":  "Math.SQRT2",
}

// Returns the equivalent function in JavaScript.
func (tr *transform) GetFuncJS(importName, funcName *ast.Ident, args []ast.Expr) (string, error) {
	var jsArgs, importStr string

	if importName != nil {
		importStr = importName.Name + "."
	}

	jsFunc, ok := Function[importStr+funcName.Name]
	if !ok {
		return "", fmt.Errorf("%s.%s: function from core library", importName, funcName)
	}

	switch funcName.Name {
	case "print", "Print":
		jsArgs = tr.joinArgsPrint(args, false)
	case "println", "Println":
		jsArgs = tr.joinArgsPrint(args, true)
	case "Printf":
		jsArgs = tr.joinArgsPrintf(args)
	}

	return fmt.Sprintf("%s(%s);", jsFunc, jsArgs), nil
}

// Returns arguments of Print, Println.
func (tr *transform) joinArgsPrint(args []ast.Expr, addLine bool) string {
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

// Matches verbs for "fmt.Printf"
var reVerb = regexp.MustCompile(`%[+\-# 0]?[bcdefgopqstvxEGTUX]`)

// Returns arguments of Printf.
// ToDo: handle "%9.3f km", or "%.*s"
func (tr *transform) joinArgsPrintf(args []ast.Expr) string {
	value := tr.getExpression(args[0])
	result := ""

	// http://golang.org/pkg/fmt/
	value = strings.Replace(value, "%%", "%", -1)

	value = reVerb.ReplaceAllString(value, "{{v}}")
	values := strings.Split(value, "{{v}}")

	for i, v := range args[1:] {
		if i != 0 {
			result += fmt.Sprintf("%s+%s", SP, SP+`"`)
		}
		result += fmt.Sprintf("%s+%s", values[i]+`"`+SP, SP+tr.getExpression(v))
	}
	result += fmt.Sprintf("%s+%s", SP, SP+`"`+values[len(values)-1])

	return result
}
