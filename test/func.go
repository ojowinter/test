package test

import (
	"fmt"
	"math"
)

// Function in the same line.
func hello() { print("Hello world!") }

func testSimpleFunc() {
	// Returns the maximum between two int a, and b.
	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	x := 3
	y := 4
	z := 5

	max_xy := max(x, y) //calling max(x, y)
	max_xz := max(x, z) //calling max(x, z)

	fmt.Printf("max(%d, %d) = %d\n", x, y, max_xy)
	fmt.Printf("max(%d, %d) = %d\n", x, z, max_xz)
	fmt.Printf("max(%d, %d) = %d\n", y, z, max(y, z)) //just call it here
}

//func testResultVar() {
// A function that returns a bool that is set to true of Sqrt is possible
// and false when not. And the actual square root of a float64
/*	MySqrt := func(f float64) (s float64, ok bool) {
		if f > 0 {
			s, ok = math.Sqrt(f), true
		} else {
			s, ok = 0, false
		}
		return s, ok
	}

	for i := -2.0; i <= 10; i++ {
		possible, sqroot := MySqrt(i)
		if possible {
			fmt.Printf("The square root of %f is %f\n", i, sqroot)
		} else {
			fmt.Printf("Sorry, no square root for %f\n", i)
		}
	}
}*/

func testByValue() {
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

func testByReference() {
	add := func(v *int) int { // notice that we give it a pointer to an int
		*v = *v + 1 // we dereference and change the value pointed by a
		return *v
	}

	x := 3
	fmt.Println("x = ", x) // "x = 3"

	x1 := add(&x) // by passing the adress of x to it
	fmt.Println("x+1 = ", x1) // "x+1 = 4"

	fmt.Println("x = ", x)    // "x = 4"
}

func testByReference2() {
	add := func(v *int, i int) { *v += i }

	value := 6
	incr := 1

	add(&value, incr)
	fmt.Println(value) // 7
}
