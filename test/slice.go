package main

import "fmt"

func shortHand() {
	// Declare an array of 10 bytes (ASCII characters). Remember: byte is uint8
	var array = [10]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'}
	// Declare a and b as slice of bytes
	var a_slice, b_slice []byte

	println("=== Slicing")
	a_slice = array[4:8]
	println(string(a_slice))
	a_slice = array[6:7]
	println(string(a_slice))

	println("\n=== Shorthands")
	a_slice = array[:3]
	println(string(a_slice))
	a_slice = array[5:]
	println(string(a_slice))
	a_slice = array[:]
	println(string(a_slice))

	println("\n=== Slice of a slice")
	a_slice = array[3:7]
	println(string(a_slice))
	b_slice = a_slice[1:3]
	println(string(b_slice))
	b_slice = a_slice[:3]
	println(string(b_slice))
	b_slice = a_slice[:]
	println(string(b_slice))
}

// Return the biggest value in a slice of ints.
func Max(slice []int) int { // The input parameter is a slice of ints.
	max := slice[0] // The first element is the max for now.
	for index := 1; index < len(slice); index++ {
		if slice[index] > max { // We found a bigger value in our slice.
			max = slice[index]
		}
	}
	return max
}

func useFunc() {
	// Declare three arrays of different sizes, to test the function Max.
	A1 := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	A2 := [4]int{1, 2, 3, 4}
	A3 := [1]int{1}

	// Declare a slice of ints.
	var slice []int

	slice = A1[:] // Take all A1 elements.
	println("The biggest value of A1 is", Max(slice))
	slice = A2[:] // Take all A2 elements.
	println("The biggest value of A2 is", Max(slice))
	slice = A3[:] // Take all A3 elements.
	println("The biggest value of A3 is", Max(slice))
}

func PrintByteSlice(name string, slice []byte) {
	s := fmt.Sprintf("%s is : [", name)
	for index := 0; index < len(slice)-1; index++ {
		s += fmt.Sprintf("%q,", slice[index])
	}
	s += fmt.Sprintf("%q]", slice[len(slice)-1])

	println(s)
}

func reference() {
	// Declare an array of 10 bytes.
	A := [10]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'}

	//declare some slices.
	slice1 := A[3:7]     // Slice1 == {'d', 'e', 'f', 'g'}
	slice2 := A[5:]      // Slice2 == {'f', 'g', 'h', 'i', 'j'}
	slice3 := slice1[:2] // Slice3 == {'d', 'e'}

	// Let's print the current content of A and the slices.
	println("=== First content of A and the slices")
	PrintByteSlice("A", A[:])
	PrintByteSlice("slice1", slice1)
	PrintByteSlice("slice2", slice2)
	PrintByteSlice("slice3", slice3)

	// Let's change the 'e' in A to 'E'.
	A[4] = 'E'
	println("\n=== Content of A and the slices, after changing 'e' to 'E' in array A")
	PrintByteSlice("A", A[:])
	PrintByteSlice("slice1", slice1)
	PrintByteSlice("slice2", slice2)
	PrintByteSlice("slice3", slice3)

	// Let's change the 'g' in slice2 to 'G'.
	slice2[1] = 'G' // Remember that 1 is actually the 2nd element in slice2!
	println("\n=== Content of A and the slices, after changing 'g' to 'G' in slice2")
	PrintByteSlice("A", A[:])
	PrintByteSlice("slice1", slice1)
	PrintByteSlice("slice2", slice2)
	PrintByteSlice("slice3", slice3)
}

func main() {
	println("\n== shortHand()\n")
	shortHand()

	println("\n== useFunc()\n")
	useFunc()

	println("\n== reference()\n")
	reference()
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

== useFunc()

The biggest value of A1 is 9
The biggest value of A2 is 4
The biggest value of A3 is 1

== reference()

=== First content of A and the slices
A is : ['a','b','c','d','e','f','g','h','i','j']
slice1 is : ['d','e','f','g']
slice2 is : ['f','g','h','i','j']
slice3 is : ['d','e']

=== Content of A and the slices, after changing 'e' to 'E' in array A
A is : ['a','b','c','d','E','f','g','h','i','j']
slice1 is : ['d','E','f','g']
slice2 is : ['f','g','h','i','j']
slice3 is : ['d','E']

=== Content of A and the slices, after changing 'g' to 'G' in slice2
A is : ['a','b','c','d','E','f','G','h','i','j']
slice1 is : ['d','E','f','G']
slice2 is : ['f','G','h','i','j']
slice3 is : ['d','E']

*/
