package main

import (
	"fmt"
	"math"
)

func singleLine() { print("Hello world!") }

var x = 10

func init() {
	x = 13
}

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
	max_xz := max(x, z) // calling max(x, z)

	fmt.Printf("max(%d, %d) = %d\n", x, y, max_xy)
	fmt.Printf("max(%d, %d) = %d\n", x, z, max_xz)
	fmt.Printf("max(%d, %d) = %d\n", y, z, max(y, z)) // just call it here
}

func twoOuputValues() {
	// Returns A+B and A*B in a single shot
	SumAndProduct := func(A, B int) (int, int) {
		return A + B, A * B
	}

	x := 3
	y := 4
	xPLUSy, xTIMESy := SumAndProduct(x, y)

	fmt.Printf("%d + %d = %d\n", x, y, xPLUSy)  // 3 + 4 = 7
	fmt.Printf("%d * %d = %d\n", x, y, xTIMESy) // 3 * 4 = 12
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

	for i := -2.0; i <= 10; i++ {
		sqroot, ok := MySqrt(i)
		if ok {
			fmt.Printf("The square root of %f is %f\n", i, sqroot)
		} else {
			fmt.Printf("Sorry, no square root for %f\n", i)
		}
	}
}

func testReturn_1() {
	MySqrt := func(f float64) (squareroot float64, ok bool) {
		if f > 0 {
			squareroot, ok = math.Sqrt(f), true
		}
		return // Omitting the output named variables, but keeping the "return".
	}

	_, check := MySqrt(5)
	fmt.Println(check) // true
}

func testReturn_2(n int) (ok bool) {
	if n > 0 {
		ok = true
	}
	return
}

func testPanic() {
	panic("unreachable")
	panic(fmt.Sprintf("not implemented: %s", "foo"))
}

func main() {
	println("\n== singleLine()\n")
	singleLine()

	println("\n== simpleFunc()\n")
	simpleFunc()

	println("\n== twoOuputValues()\n")
	twoOuputValues()

	println("\n== resultVariable()\n")
	resultVariable()

	println("\n== testReturn_1()\n")
	testReturn_1()

	println("\n== testReturn_2(-1)\n")
	println(testReturn_2(-1))

	println("\n== testPanic()\n")
	testPanic()
}

/*
== singleLine()

Hello world!
== simpleFunc()

max(3, 4) = 4
max(3, 5) = 5
max(4, 5) = 5

== twoOuputValues()

3 + 4 = 7
3 * 4 = 12

== resultVariable()

Sorry, no square root for -2.000000
Sorry, no square root for -1.000000
Sorry, no square root for 0.000000
The square root of 1.000000 is 1.000000
The square root of 2.000000 is 1.414214
The square root of 3.000000 is 1.732051
The square root of 4.000000 is 2.000000
The square root of 5.000000 is 2.236068
The square root of 6.000000 is 2.449490
The square root of 7.000000 is 2.645751
The square root of 8.000000 is 2.828427
The square root of 9.000000 is 3.000000
The square root of 10.000000 is 3.162278

== testReturn_1()

true

== testReturn_2(-1)

false

== testPanic()

panic: unreachable

goroutine 1 [running]:
main.testPanic()
	/home/neo/Público/Proyectos/Go/GoScript/test/func.go:91 +0x47
main.main()
	/home/neo/Público/Proyectos/Go/GoScript/test/func.go:115 +0x108
*/
