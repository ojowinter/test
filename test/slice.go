package main

//import "fmt"

func shortHand() {
	// Declare an array of 10 bytes (ASCII characters). Remember: byte is uint8
	var array = [10]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'}
	// Declare a and b as slice of bytes
	var a_slice, b_slice []byte

	println("=== Slicing")
	a_slice = array[4:8]; println(string(a_slice))
	a_slice = array[6:7]; println(string(a_slice))

	println("\n=== Shorthands")
	a_slice = array[:3]; println(string(a_slice))
	a_slice = array[5:]; println(string(a_slice))
	a_slice = array[:]; println(string(a_slice))

	println("\n=== Slice of a slice")
	a_slice = array[3:7]; println(string(a_slice))
	b_slice = a_slice[1:3]; println(string(b_slice))
	b_slice = a_slice[:3]; println(string(b_slice))
	b_slice = a_slice[:]; println(string(b_slice))
}

func main() {
	println("\n== shortHand()\n")
	shortHand()

}

/*
== shortHand()

=== Slicing
efgh
g

=== Shorthands
abc
fghij
abcdefghij

=== Slice of a slice
defg
ef
def
defg


*/
