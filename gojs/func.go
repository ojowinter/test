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

var validImport = []string{"fmt", "math", "rand"}

// Functions which can be transformed.
// The empty values are to indicate that the package (in the key) have any
// function to be transformed.
var Function = map[string]string{
	"print":   "alert",
	"println": "alert",

	"fmt.Print":   "alert",
	"fmt.Println": "alert",
	"fmt.Printf":  "alert",

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
	"math.Sin":     "Math.sin",
	"math.Sqrt":    "Math.sqrt",
	"math.Tan":     "Math.tan",
	// https://developer.mozilla.org/en/JavaScript/Reference/Global_Objects/Math/round
	//"math.":        "Math.round",
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
// http://golang.org/pkg/fmt/
var (
	reVerb       = regexp.MustCompile(`%[+\-# 0]?[bcdefgopqstvxEGTUX]`)
	reVerbNumber = regexp.MustCompile(`%[0-9]+[.]?[0-9]*[bcdefgoqxEGUX]`)
	reVerbString = regexp.MustCompile(`%[0-9]+[.]?[0-9]*[sqxX]`)
)

const VERB = "{{v}}"

// Returns arguments of Printf.
func (tr *transform) joinArgsPrintf(args []ast.Expr) string {
	result := ""
	value := tr.getExpression(args[0])

	value = strings.Replace(value, "%%", "%", -1) // literal percent sign
	value = reVerb.ReplaceAllString(value, VERB)

	if reVerbNumber.MatchString(value) {
		value = reVerbNumber.ReplaceAllString(value, VERB)
	}
	if reVerbString.MatchString(value) {
		value = reVerbString.ReplaceAllString(value, VERB)
	}

	values := strings.Split(value, VERB)

	for i, v := range args[1:] {
		if i != 0 {
			result += fmt.Sprintf("%s+%s", SP, SP+`"`)
		}
		result += fmt.Sprintf("%s+%s", values[i]+`"`+SP, SP+tr.getExpression(v))
	}
	result += fmt.Sprintf("%s+%s", SP, SP+`"`+values[len(values)-1])

	return result
}
