package main

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

	println("The string \"hello\" is:", hello)                      // Hello, mina-san!
	println("The string pointed to by \"helloPtr\" is:", *helloPtr) // idem
	println("The value of \"i\" is:", i)                            // 6
	println("The value pointed to by \"iPtr\" is:", *iPtr)          // idem
}

func access_2() {
	x := 3
	y := &x

	*y++
	println("x is:", x) // 4

	*y++
	println("x is:", x) // 5
}

func allocation() {
	sum := 0
	var doubleSum *int // a pointer to int
	for i := 0; i < 10; i++ {
		sum += i
	}

	doubleSum = new(int) // allocate memory for an int and make doubleSum point to it
	*doubleSum = sum * 2 // use the allocated memory, by dereferencing doubleSum

	println("The sum of numbers from 0 to 10 is:", sum) // 45
	println("The double of this sum is:", *doubleSum)   // 90
}

func parameterByValue() {
	// Returns 1 plus its input parameter
	var add = func(v int) int {
		v = v + 1
		return v
	}

	x := 3
	println("x =", x) // "x = 3"

	x1 := add(x)
	println("x+1 =", x1) // "x+1 = 4"
	println("x =", x)    // "x = 3"
}

func byReference_1() {
	add := func(v *int) int { // notice that we give it a pointer to an int
		*v = *v + 1 // we dereference and change the value pointed by a
		return *v
	}

	x := 3
	println("x =", x) // "x = 3"

	x1 := add(&x)             // by passing the adress of x to it
	println("x+1 =", x1) // "x+1 = 4"
	println("x =", x)    // "x = 4"

	x1 = add(&x)
	println("x+1 =", x1) // "x+1 = 5"
	println("x =", x)    // "x = 5"
}

func byReference_2() {
	add := func(v *int, i int) { *v += i }

	value := 6
	incr := 1

	add(&value, incr)
	println(value) // 7

	add(&value, incr)
	println(value) // 8
}

func byReference_3() {
	x := 3
	f := func() {
		x = 4
	}
	y := &x

	f()
	println(*y) // 4
}

func main() {
	println("\n== declaration()\n")
	declaration()
	println("\n== showAddress()\n")
	showAddress()
	println("\n== access_1()\n")
	access_1()
	println("\n== access_2()\n")
	access_2()
	println("\n== allocation()\n")
	allocation()
	println("\n== parameterByValue()\n")
	parameterByValue()
	println("\n== byReference_1()\n")
	byReference_1()
	println("\n== byReference_2()\n")
	byReference_2()
	println("\n== byReference_3()\n")
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

== access_1()

The string "hello" is: Hello, mina-san!
The string pointed to by "helloPtr" is: Hello, mina-san!
The value of "i" is: 6
The value pointed to by "iPtr" is: 6

== access_2()

x is: 4
x is: 5

== allocation()

The sum of numbers from 0 to 10 is: 45
The double of this sum is: 90

== parameterByValue()

x = 3
x+1 = 4
x = 3

== byReference_1()

x = 3
x+1 = 4
x = 4
x+1 = 5
x = 5

== byReference_2()

7
8

== byReference_3()

4
*/
