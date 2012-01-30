package main

import "math"

type Rectangle struct {
	width, height float64
}

func area(r Rectangle) float64 {
	return r.width * r.height
}

func noMethod() {
	r1 := Rectangle{12, 2}
	r2 := Rectangle{9, 4}
	println("Area of r1 is:", area(r1))
	println("Area of r2 is:", area(r2))
}

// * * *

func (r Rectangle) area() float64 {
	return r.width * r.height
}

type Circle struct {
	radius float64
}

func (c Circle) area() float64 {
	return c.radius * c.radius * math.Pi
}

func method() {
	r1 := Rectangle{12, 2}
	r2 := Rectangle{9, 4}
	c1 := Circle{10}
	c2 := Circle{25}

	println("Area of r1 is:", r1.area())
	println("Area of r2 is:", r2.area())
	println("Area of c1 is:", c1.area())
	println("Area of c2 is:", c2.area())
}

// * * *

type SliceOfints []int
type AgesByNames map[string]int

func (s SliceOfints) sum() int {
	sum := 0
	for _, value := range s {
		sum += value
	}
	return sum
}

func (people AgesByNames) older() string {
	a := 0
	n := ""
	for key, value := range people {
		if value > a {
			a = value
			n = key
		}
	}
	return n
}

func withNamedType() {
	s := SliceOfints{1, 2, 3, 4, 5}
	folks := AgesByNames{
		"Bob":   36,
		"Mike":  44,
		"Jane":  30,
		"Popey": 100,
	}

	println("The sum of ints in the slice s is:", s.sum())
	println("The older in the map folks is:", folks.older())
}

// * * *

func main() {
	println("\n== noMethod()\n")
	noMethod()
	println("\n== method()\n")
	method()
	println("\n== withNamedType()\n")
	withNamedType()
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
