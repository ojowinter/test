package test

import (
	"fmt"
	"math"
)

// Function in the same line.
func hello() { print("Hello world!") }

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
		return A+B, A*B
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

func emptyReturn(f float64) (squareroot float64, ok bool) {
	if f > 0 {
		squareroot, ok = math.Sqrt(f), true
	}
	return // Omitting the output named variables, but keeping the "return".
}

func emptyReturn2(n int) (ok bool) {
	if n > 0 {
		ok = true
	}
	return
}

func parameterByValue() {
	// Returns 1 plus its input parameter
	var add = func(v int) int {
		v = v + 1
		return v
	}

	x := 3
	fmt.Println("x = ", x) // "x = 3"

	x1 := add(x)
	fmt.Println("x+1 = ", x1) // "x+1 = 4"
	fmt.Println("x = ", x)    // "x = 3"
}

func parameterByReference() {
	add := func(v *int) int { // notice that we give it a pointer to an int
		*v = *v + 1 // we dereference and change the value pointed by a
		return *v
	}

	x := 3
	fmt.Println("x = ", x) // "x = 3"

	x1 := add(&x) // by passing the adress of x to it
	fmt.Println("x+1 = ", x1) // "x+1 = 4"
	fmt.Println("x = ", x)    // "x = 4"

	x1 = add(&x)
	fmt.Println("x+1 = ", x1) // "x+1 = 5"
	fmt.Println("x = ", x)    // "x = 5"
}

func byReference2() {
	add := func(v *int, i int) { *v += i }

	value := 6
	incr := 1

	add(&value, incr)
	fmt.Println(value) // 7

	add(&value, incr)
	fmt.Println(value) // 8
}

func byReference3() {
	x := 3
	y := &x

	*y++
	println(x) // 4

	*y++
	println(x) // 5
}

func byReference4() {
	x := 3
	f := func(){
		x = 4
	}
	y := &x

	f()
	println(*y) // 4
}

/*
function byReference4() {
	var x = [3];
	var f = function() {
		x[0] = 4;
	};
	var y = x;

	f();
	console.log(y[0] + "\n");
}*/
