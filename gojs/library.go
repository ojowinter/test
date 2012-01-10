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

// Functions that can be transformed since JavaScript has an equivalent one.
var Function = map[string]string{
	"panic": "throw new Error",

	"print":   "console.log",
	"println": "console.log",

	"fmt.Print":   "alert",
	"fmt.Println": "alert",
	"fmt.Printf":  "alert",
	"fmt.Sprintf": "",

	"math.Abs":   "Math.abs",
	"math.Acos":  "Math.acos",
	"math.Asin":  "Math.asin",
	"math.Atan":  "Math.atan",
	"math.Atan2": "Math.atan2",
	"math.Ceil":  "Math.ceil",
	"math.Cos":   "Math.cos",
	"math.Exp":   "Math.exp",
	"math.Floor": "Math.floor",
	"math.Log":   "Math.log",
	"math.Max":   "Math.max",
	"math.Min":   "Math.min",
	"math.Pow":   "Math.pow",
	"math.Sin":   "Math.sin",
	"math.Sqrt":  "Math.sqrt",
	"math.Tan":   "Math.tan",
	// https://developer.mozilla.org/en/JavaScript/Reference/Global_Objects/Math/round
	//"math.":      "Math.round",

	"rand.Float32": "Math.random",
	"rand.Float64": "Math.random",
}

// Returns the arguments of a Go function, formatted for JS.
func (tr *transform) GetArgs(funcName string, args []ast.Expr) string {
	var jsArgs string

	switch funcName {
	case "print", "fmt.Print":
		jsArgs = tr.joinArgsPrint(args, false)
	case "println", "fmt.Println":
		jsArgs = tr.joinArgsPrint(args, true)
	case "fmt.Printf", "fmt.Sprintf":
		jsArgs = tr.joinArgsPrintf(args)
	default:
		for i, v := range args {
			if i != 0 {
				jsArgs += "," + SP
			}
			jsArgs += tr.getExpression(v).String()
		}
	}

	return jsArgs
}

//
// === Utility

func (tr *transform) checkLib(selector *ast.SelectorExpr) (goName, jsName string, err error) {
	var ok bool

	goName = selector.X.(*ast.Ident).Name + "." + selector.Sel.Name

	jsName, ok = Function[goName]
	if !ok {
		jsName, ok = Constant[goName]
	}

	if !ok {
		return "", "", fmt.Errorf("%s: %q not supported in JS",
			tr.fset.Position(selector.Sel.Pos()), goName)
	}

	return goName, jsName, nil
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
		expr := tr.getExpression(v).String()

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
	reVerb      = regexp.MustCompile(`%[+\-# 0]?[bcdefgopqstvxEGTUX]`)
	reVerbWidth = regexp.MustCompile(`%[0-9]+[.]?[0-9]*[bcdefgoqxEGUXsqxX]`)
)

const VERB = "{{v}}"

// Returns arguments of Printf.
func (tr *transform) joinArgsPrintf(args []ast.Expr) string {
	result := ""

	// === Format
	format := tr.getExpression(args[0]).String()

	format = strings.Replace(format, "%%", "%", -1) // literal percent sign
	format = reVerb.ReplaceAllString(format, VERB)

	if reVerbWidth.MatchString(format) {
		format = reVerbWidth.ReplaceAllString(format, VERB)
	}
	// ===

	values := strings.Split(format, VERB)

	for i, v := range args[1:] {
		if i != 0 {
			result += fmt.Sprintf("%s+%s", SP, SP+`"`)
		}
		result += fmt.Sprintf("%s+%s", values[i]+`"`+SP, SP+tr.getExpression(v).String())
	}
	result += fmt.Sprintf("%s+%s", SP, SP+`"`+values[len(values)-1])

	return result
}
