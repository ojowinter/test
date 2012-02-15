package main

import ("fmt"; "math")

type Rectangle struct {
	width, height float64
}

func noMethod() {
	area := func(r Rectangle) float64 {
		return r.width * r.height
	}

	r1 := Rectangle{12, 2}

	// Checking
	if area(r1) == 24 && area(Rectangle{9, 4}) == 36 {
		println("[OK]")
	} else {
		fmt.Println("[Error] Area of r1 is:", area(r1))
		fmt.Println("\tArea of \"Rectangle{9, 4}\" is:", area(Rectangle{9, 4}))
	}
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

	// Checking
	if r1.area() == 24 && r2.area() == 36 {
		println("[OK] rectangle")
	} else {
		fmt.Println("[Error] Area of r1 is:", r1.area())
		fmt.Println("\tArea of r2 is:", r2.area())
	}

	if c1.area() == 314.1592653589793 && c2.area() == 1963.4954084936207 {
		println("[OK] circle")
	} else {
		fmt.Println("[Error] Area of c1 is:", c1.area())
		fmt.Println("\tArea of c2 is:", c2.area())
	}
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

	// Checking
	if s.sum() == 15 {
		println("[OK] sum")
	} else {
		fmt.Println("[Error] The sum of ints in the slice s is:", s.sum())
	}

	if folks.older() == "Popey" {
		println("[OK] older")
	} else {
		fmt.Println("[Error] The older in the map folks is:", folks.older())
	}
}

// * * *

func main() {
	println("\n== noMethod")
	noMethod()
	println("\n== method")
	method()
	println("\n== withNamedType")
	withNamedType()
}
