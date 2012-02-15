package main

import (
	"fmt"
	"math"
)

var x = 10

func init() {
	x = 13
}

func testInit() {
	code := ""
	if x == 13 {
		code = "OK"
	} else {
		code = "Error"
	}
	println("[" + code + "]")
}

func singleLine() { _ = "Hello world!"; println("[OK]") }

func simpleFunc() {
	// Returns the maximum between two int a, and b.
	var max = func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	x := 3
	y := 4
	z := 5

	max_xy := max(x, y) // calling max(x, y)
	// Checking
	if max_xy == 4 {
		println("[OK] x,y")
	} else {
		fmt.Printf("[Error] max(%d, %d) = %d\n", x, y, max_xy)
	}
	//==

	max_xz := max(x, z) // calling max(x, z)
	// Checking
	if max_xz == 5 {
		println("[OK] x,z")
	} else {
		fmt.Printf("[Error] max(%d, %d) = %d\n", x, z, max_xz)
	}
	//==

	// Checking
	if max(y, z) == 5 { // just call it here
		println("[OK] y,z")
	} else {
		fmt.Printf("[Error] max(%d, %d) = %d\n", y, z, max(y, z))
	}
	//==
}

func twoOuputValues() {
	// Returns A+B and A*B in a single shot
	SumAndProduct := func(A, B int) (int, int) {
		return A + B, A * B
	}

	x := 3
	y := 4
	xPLUSy, xTIMESy := SumAndProduct(x, y)

	// Checking
	if xPLUSy == 7 && xTIMESy == 12 {
		println("[OK]")
	} else {
		fmt.Printf("[Error] %d + %d = %d\t", x, y, xPLUSy)
		fmt.Printf("%d * %d = %d\n", x, y, xTIMESy)
	}
	//==
}

func resultVariable() {
	// A function that returns a bool that is set to true of Sqrt is possible
	// and false when not. And the actual square root of a float64
	MySqrt := func(f float64) (s float64, ok bool) {
		if f > 0 {
			s, ok = math.Sqrt(f), true
		}
		return s, ok
	}

	results := map[float64]float64 {
		1: 1,
		2: 1.4142135623730951,
		3: 1.7320508075688772,
		4: 2,
		5: 2.23606797749979,
		6: 2.449489742783178,
		7: 2.6457513110645907,
		8: 2.8284271247461903,
		9: 3,
		10: 3.1622776601683795,
	}

	err := false
	for i := -2.0; i <= 10; i++ {
		sqroot, ok := MySqrt(i)
		if ok {
			if sqroot != results[i] {
				fmt.Printf("[Error] The square root of %v is %v\n", i, sqroot)
				err = true
			}
		} else {
			if i != -2.0 && i != -1.0 && i != 0 {
				fmt.Printf("[Error] The square root for %v should not be run\n", i)
				err = true
			}
		}
	}
	if !err {
		println("[OK]")
	}
}

func testReturn() {
	MySqrt := func(f float64) (squareroot float64, ok bool) {
		if f > 0 {
			squareroot, ok = math.Sqrt(f), true
		}
		return // Omitting the output named variables, but keeping the "return".
	}

	_, check := MySqrt(5)

	// Checking
	code := ""
	if check {
		code = "OK"
	} else {
		code = "Error"
	}
	println("[" + code + "]")

	if _, ok := MySqrt(0); !ok {
		code = "OK"
	} else {
		code = "Error"
	}
	println("[" + code + "]")
	//==
}

func testPanic() {
	panic("unreachable")
	panic(fmt.Sprintf("not implemented: %s", "foo"))
}

func main() {
	println("\n== testInit")
	testInit()
	println("\n== singleLine")
	singleLine()
	println("\n== simpleFunc")
	simpleFunc()
	println("\n== twoOuputValues")
	twoOuputValues()
	println("\n== resultVariable")
	resultVariable()
	println("\n== testReturn")
	testReturn()
	println("\n== testPanic")
	testPanic()
}
