package main

import "fmt"

// Global declaration of a pointer
var i int
var hello string
var p *int

func init() {
	p = &i             // p points to i (p stores the address of i)
	helloPtr := &hello // pointer variable of type *string and it points hello
	println("helloPtr:", helloPtr)
}

func declaration() {
	var i int
	var hello string
	var p *int

	p = &i
	helloPtr := &hello
	println("p: ", p, "\nhelloPtr:", helloPtr)
}

func showAddress() {
	var (
		i     int     = 9
		hello string  = "Hello world"
		pi    float32 = 3.14
		b     bool    = true
	)

	println("Hexadecimal address of 'i' is:", &i)
	println("Hexadecimal address of 'hello' is:", &hello)
	println("Hexadecimal address of 'pi' is:", &pi)
	println("Hexadecimal address of 'b' is:", &b)
}

func access_1() {
	hello := "Hello, mina-san!"

	var helloPtr *string
	helloPtr = &hello

	i := 6
	iPtr := &i

	// Checking
	if hello == "Hello, mina-san!" && *helloPtr == "Hello, mina-san!" {
		println("[OK] string")
	} else {
		fmt.Printf("[Error] The string \"hello\" is: %v\n", hello)
		fmt.Printf("\tThe string pointed to by \"helloPtr\" is: %v\n", *helloPtr)
	}

	if i == 6 && *iPtr == 6 {
		println("[OK] int")
	} else {
		fmt.Printf("[Error] The value of \"i\" is: %v\n", i)
		fmt.Printf("\tThe value pointed to by \"iPtr\" is: %v\n", *iPtr)
	}
}

func access_2() {
	x := 3
	y := &x

	*y++
	// Checking
	if x == 4 {
		println("[OK]")
	} else {
		fmt.Println("[Error] x is:", x)
	}
	//==

	*y++
	// Checking
	if x == 5 {
		println("[OK]")
	} else {
		fmt.Println("[Error] x is:", x)
	}
}

func allocation() {
	sum := 0
	var doubleSum *int // a pointer to int
	for i := 0; i < 10; i++ {
		sum += i
	}

	doubleSum = new(int) // allocate memory for an int and make doubleSum point to it
	*doubleSum = sum * 2 // use the allocated memory, by dereferencing doubleSum

	// Checking
	if sum == 45 && *doubleSum == 90 {
		println("[OK]")
	} else {
		fmt.Printf("[Error] The sum of numbers from 0 to 10 is: %v\n", sum)
		fmt.Printf("\tThe double of this sum is: %v\n", *doubleSum)
	}
}

func parameterByValue() {
	// Returns 1 plus its input parameter
	var add = func(v int) int {
		v = v + 1
		return v
	}

	x := 3
	x1 := add(x)

	// Checking
	if x1 == 4 && x == 3 {
		println("[OK]")
	} else {
		fmt.Printf("[Error] x+1 = %v\n", x1)
		fmt.Printf("\tx = %v\n", x)
	}
}

func byReference_1() {
	add := func(v *int) int { // notice that we give it a pointer to an int
		*v = *v + 1 // we dereference and change the value pointed by a
		return *v
	}

	x := 3

	x1 := add(&x)             // by passing the adress of x to it
	// Checking
	if x1 == 4 && x == 4 {
		println("[OK]")
	} else {
		fmt.Printf("[Error] x+1 = %v\n", x1)
		fmt.Printf("\tx = %v\n", x)
	}
	//==

	x1 = add(&x)
	// Checking
	if x1 == 5 && x == 5 {
		println("[OK]")
	} else {
		fmt.Printf("[Error] x+1 = %v\n", x1)
		fmt.Printf("\tx = %v\n", x)
	}
	//==
}

func byReference_2() {
	add := func(v *int, i int) { *v += i }

	value := 6
	incr := 1

	add(&value, incr)
	// Checking
	if value == 7 {
		println("[OK]")
	} else {
		fmt.Printf("[Error] value = %v\n", value)
	}
	//==

	add(&value, incr)
	// Checking
	if value == 8 {
		println("[OK]")
	} else {
		fmt.Printf("[Error] value = %v\n", value)
	}
	//==
}

func byReference_3() {
	x := 3
	f := func() {
		x = 4
	}
	y := &x

	f()
	if *y == 4 {
		println("[OK]")
	} else {
		fmt.Println("[Error] y = ", *y)
	}
}

func main() {
	println("\n== declaration")
	declaration()
	println("\n== showAddress")
	showAddress()
	println("\n== access_1")
	access_1()
	println("\n== access_2")
	access_2()
	println("\n== allocation")
	allocation()
	println("\n== parameterByValue")
	parameterByValue()
	println("\n== byReference_1")
	byReference_1()
	println("\n== byReference_2")
	byReference_2()
	println("\n== byReference_3")
	byReference_3()
}

/*
helloPtr: 0x428038

== declaration()

p:  0x7feaabc2af5c 
helloPtr: 0x7feaabc2af68

== showAddress()

Hexadecimal address of 'i' is: 0x7feaabc2af68
Hexadecimal address of 'hello' is: 0x7feaabc2af70
Hexadecimal address of 'pi' is: 0x7feaabc2af6c
Hexadecimal address of 'b' is: 0x7feaabc2af67

*/
