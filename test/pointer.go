package main

// Global declaration of a pointer
var i int
var hello string
var p *int

func init() {
	p = &i             // p points to i (p stores the address of i)
	helloPtr := &hello // pointer variable of type *string and it points hello
}

func declaration() {
	var i int
	var hello string
	var p *int

	p = &i
	helloPtr := &hello
}

func showAddress() {
	var (
		i     int       = 9
		hello string    = "Hello world"
		pi    float32   = 3.14
		b     bool      = true
	)

	println("Hexadecimal address of 'i' is:", &i)
	println("Hexadecimal address of 'hello' is:", &hello)
	println("Hexadecimal address of 'pi' is:", &pi)
	println("Hexadecimal address of 'b' is:", &b)
}

func access() {
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

func main() {
	println("\n== showAddress()\n")
	showAddress()

	println("\n== access()\n")
	access()

	println("\n== access_2()\n")
	access_2()

	println("\n== allocation()\n")
	allocation()
}
