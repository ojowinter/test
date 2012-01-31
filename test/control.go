package main

import "fmt"

func testIf() {
	x := 5
	code := ""

	// Simple
	if x > 10 {
		code = "Error"
	} else {
		code = "OK"
	}
	println("[" + code + "] simple")

	// Leading initial short
	if x := 12; x > 10 {
		code = "OK"
	} else {
		code = "Error"
	}
	println("[" + code + "] with statement")

	// Multiple if/else
	i := 7

	if i == 3 {
		code = "Error"
	} else if i < 3 {
		code = "Error"
	} else {
		code = "OK"
	}
	println("[" + code + "] multiple")
}

func testSwitch() {
	i := 10
	code := ""

	// Simple
	switch i {
	default:
		code = "Error"
	case 1:
		code = "Error"
	case 2, 3, 4:
		code = "Error"
	case 10:
		code = "OK"
	}
	println("[" + code + "] simple")

	// Without expression
	switch i = 5; {
	case i < 10:
		code = "OK"
	case i > 10, i < 0:
		code = "Error"
	case i == 10:
		code = "Error"
	default:
		code = "Error"
	}
	println("[" + code + "] with statement")

	switch {
	case i == 5:
		code = "OK"
	}
	println("[" + code + "] without expression")

	// With fallthrough
	switch i {
	case 4:
		code = "Error"
		fallthrough
	case 5:
		code = "Error"
		fallthrough
	case 6:
		code = "Error"
		fallthrough
	case 7:
		code = "OK"
	case 8:
		code = "Error"
	default:
		code = "Error"
	}
	println("[" + code + "] with fallthrough")
}

func testFor() {
	sum := 0

	// Simple
	for i := 0; i < 10; i++ {
		sum += i
	}
	// Checking
	code := ""
	if sum == 45 {
		code = "OK"
	} else {
		code = "Error"
	}
	println("[" + code + "] simple")
	//==

	// Expression1 and expression3 are omitted here
	sum = 1
	for ; sum < 1000; {
		sum += sum
	}
	// Checking
	if sum == 1024 {
		code = "OK"
	} else {
		code = "Error"
	}
	println("[" + code + "] 2 expressions omitted")
	//==

	// Expression1 and expression3 are omitted here, and semicolons gone
	sum = 1
	for sum < 1000 {
		sum += sum
	}
	// Checking
	if sum == 1024 {
		code = "OK"
	} else {
		code = "Error"
	}
	println("[" + code + "] 2 expressions omitted, no semicolons")
	//==

	// Infinite loop (limited to show the output), no semicolons at all
	i := 0
	s := ""
	for {
		i++
		if i == 3 {
			s = fmt.Sprintf("%d", i)
			break
		}
	}
	// Checking
	if s == "3" {
		code = "OK"
	} else {
		code = "Error"
	}
	println("[" + code + "] infinite loop")
	//==

	// break
	s = ""
	for i := 10; i > 0; i-- {
		if i < 5 {
			break
		}
		s += fmt.Sprintf("%d ", i)
	}
	// Checking
	if s == "10 9 8 7 6 5 " {
		println("[OK] break")
	} else {
		fmt.Printf("[Error] value in break: %s\n", s)
	}
	//==

	// continue
	s = ""
	for i := 10; i > 0; i-- {
		if i == 5 {
			continue
		}
		s += fmt.Sprintf("%d ", i)
	}
	// Checking
	if s == "10 9 8 7 6 4 3 2 1 " {
		println("[OK] continue")
	} else {
		fmt.Printf("[Error] value in continue: %s\n", s)
	}
	//==
}

func testRange() {
	s := []int{2, 3, 5}

	for i, v := range s {
		println("key:", i, "value:", v)
	}
}

func main() {
	println("\n== testIf")
	testIf()
	println("\n== testSwitch")
	testSwitch()
	println("\n== testFor")
	testFor()
/*	println("\n== testRange")
	testRange()*/
}

/*

== testRange()

key: 0 value: 2
key: 1 value: 3
key: 2 value: 5
*/
