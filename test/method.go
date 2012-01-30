package main

import "math"

type Rectangle struct {
	width, height float64
}

func (r Rectangle) area() float64 {
	return r.width * r.height
}

type Circle struct {
	radius float64
}

func (c Circle) area() float64 {
	return c.radius * c.radius * math.Pi
}

func noMethod() {
	area := func(r Rectangle) float64 {
		return r.width * r.height
	}

	r1 := Rectangle{12, 2}
	r2 := Rectangle{9, 4}
	println("Area of r1 is:", area(r1))
	println("Area of r2 is:", area(r2))
}

func test_1() {
	r1 := Rectangle{12, 2}
	r2 := Rectangle{9, 4}
	c1 := Circle{10}
	c2 := Circle{25}

	//Now look how we call our methods.
	println("Area of r1 is:", r1.area())
	println("Area of r2 is:", r2.area())
	println("Area of c1 is:", c1.area())
	println("Area of c2 is:", c2.area())
}

func main() {
	println("\n== noMethod()\n")
	noMethod()
	println("\n== test_1()\n")
	test_1()
}

/*
== noMethod()

Area of r1 is: +2.400000e+001
Area of r2 is: +3.600000e+001

== test_1()

Area of r1 is: +2.400000e+001
Area of r2 is: +3.600000e+001
Area of c1 is: +3.141593e+002
Area of c2 is: +1.963495e+003

*/
