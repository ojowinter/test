package main

func testIf() {
	x := 5

	// Simple
	if x > 10 {
		println("x =", x, "is greater than 10")
	} else {
		println("x =", x, "is less than 10")
	}

	// Leading initial short
	if x := 12; x > 10 {
		println("x =", x, "is greater than 10")
	} else {
		println("x =", x, "is less than 10")
	}

	// Multiple if/else
	i := 7

	if i == 3 {
		println("i =", i, "is equal to 3")
	} else if i < 3 {
		println("i =", i, "is less than 3")
	} else {
		println("i =", i, "is greater than 3")
	}
}

func testSwitch() string {
	i := 10

	// Simple
	switch i {
	case 1:
		println("i =", i, "is equal to 1")
	case 2, 3, 4:
		println("i =", i, "is equal to 2, 3 or 4")
	case 10:
		println("i =", i, "is equal to 10")
	default:
		println("All I know is that i is an integer")
	}

	// Without expression
	switch i = 5; {
	case i < 10:
		println("i =", i, "is less than 10")
	case i > 10, i < 0:
		println("i =", i, "is either bigger than 10 or less than 0")
	case i == 10:
		println("i =", i, "is equal to 10")
	default:
		println("This won't be printed anyway")
	}

	switch {
	case i == 5:
		println("i is 5")
	}

	// With fallthrough
	switch i {
	case 4:
		println("was <= 4")
		fallthrough
	case 5:
		println("was <= 5")
		fallthrough
	case 6:
		println("was <= 6")
		fallthrough
	case 7:
		println("was <= 7")
		fallthrough
	case 8:
		println("was <= 8")
		//fallthrough
	default:
		println("default case")
	}

	// With return
	switch i {
	default:
	case 1, 3, 5, 7, 9:
		return "odd"
	case 2, 4, 6, 8:
		return "even"
	}

	return ""
}

func testFor() {
	sum := 0

	// Simple
	for i := 0; i < 10; i++ {
		sum += i
	}
	println("sum is equal to", sum)

	// Expression1 and expression3 are omitted here
	sum = 1
	for sum < 1000 {
		sum += sum
	}
	println("sum is equal to", sum)

	// Expression1 and expression3 are omitted here, and semicolons gone
	sum = 1
	for sum < 1000 {
		sum += sum
	}
	println("sum is equal to", sum)

	// Infinite loop (limited to show the output), no semicolons at all
	i := 0
	for {
		println("I loop for ever!")
		i++
		if i == 3 {
			break
		}
	}

	// break
	print("break on 5: ")
	for i := 10; i > 0; i-- {
		if i < 5 {
			break
		}
		print(i, " ")
	}

	// continue
	print("\nskip 5: ")
	for i := 10; i > 0; i-- {
		if i == 5 {
			continue
		}
		print(i, " ")
	}
	println()
}

func testRange() {
	s := []int{2, 3, 5}

	for i, v := range s {
		println("key:", i, "value:", v)
	}
}

func main() {
	println("\n== testIf()\n")
	testIf()
	println("\n== testSwitch()\n")
	println(testSwitch())
	println("\n== testFor()\n")
	testFor()
	println("\n== testRange()\n")
	testRange()
}

/*
== testIf()

x = 5 is less than 10
x = 12 is greater than 10
i = 7 is greater than 3

== testSwitch()

i = 10 is equal to 10
i = 5 is less than 10
i is 5
was <= 5
was <= 6
was <= 7
was <= 8
odd

== testFor()

sum is equal to 45
sum is equal to 1024
sum is equal to 1024
I loop for ever!
I loop for ever!
I loop for ever!
break on 5: 10 9 8 7 6 5 
skip 5: 10 9 8 7 6 4 3 2 1 

== testRange()

key: 0 value: 2
key: 1 value: 3
key: 2 value: 5
*/
