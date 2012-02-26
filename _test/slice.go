package main

import "fmt"

func valueNil() {
	var s []byte

	// Checking
	msg := "value"
	if s == nil {
		println("[OK]", msg)
	} else {
		fmt.Println("[Error]", msg)
	}

	msg = "length"
	if len(s) == 0 {
		println("[OK]", msg)
	} else {
		fmt.Println("[Error]", msg)
	}

	msg = "capacity"
	if cap(s) == 0 {
		println("[OK]", msg)
	} else {
		fmt.Println("[Error]", msg)
	}
	//==
}

func shortHand() {
	// Declare an array of 10 bytes (ASCII characters). Remember: byte is uint8
	var array = [10]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'}
	// Declare a and b as slice of bytes
	var a_slice, b_slice []byte

	msg := "slicing"

	a_slice = array[4:8]
	// Checking
	if string(a_slice) == "efgh" && len(a_slice) == 4 && cap(a_slice) == 6 {
		println("[OK]", msg)
	} else {
		fmt.Println("[Error]", msg)
	}
	//==

	a_slice = array[6:7]
	// Checking
	if string(a_slice) == "g" {
		println("[OK]")
	} else {
		fmt.Println("[Error]")
	}
	//==

	msg = "shorthand"

	a_slice = array[:3]
	// Checking
	if string(a_slice) == "abc" && len(a_slice) == 3 && cap(a_slice) == 10 {
		println("[OK]", msg)
	} else {
		fmt.Println("[Error]", msg)
	}
	//==

	a_slice = array[5:]
	// Checking
	if string(a_slice) == "fghij" {
		println("[OK]")
	} else {
		fmt.Println("[Error]")
	}
	//==

	a_slice = array[:]
	// Checking
	if string(a_slice) == "abcdefghij" {
		println("[OK]")
	} else {
		fmt.Println("[Error]")
	}
	//==

	msg = "slice of a slice"

	a_slice = array[3:7]
	// Checking
	if string(a_slice) == "defg" && len(a_slice) == 4 && cap(a_slice) == 7 {
		println("[OK]", msg)
	} else {
		fmt.Println("[Error]", msg)
	}
	//==

	b_slice = a_slice[1:3]
	// Checking
	if string(b_slice) == "ef" && len(b_slice) == 2 && cap(b_slice) == 6 {
		println("[OK]")
	} else {
		fmt.Println("[Error]")
	}
	//==

	b_slice = a_slice[:3]
	// Checking
	if string(b_slice) == "def" {
		println("[OK]")
	} else {
		fmt.Println("[Error]")
	}
	//==

	b_slice = a_slice[:]
	// Checking
	if string(b_slice) == "defg" {
		println("[OK]")
	} else {
		fmt.Println("[Error]")
	}
	//==
}

// * * *

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
	// Checking
	if Max(slice) == 9 {
		println("[OK] A1")
	} else {
		fmt.Println("[Error] A1")
	}
	//==

	slice = A2[:] // Take all A2 elements.
	// Checking
	if Max(slice) == 4 {
		println("[OK] A2")
	} else {
		fmt.Println("[Error] A2")
	}
	//==

	slice = A3[:] // Take all A3 elements.
	// Checking
	if Max(slice) == 1 {
		println("[OK] A3")
	} else {
		fmt.Println("[Error] A3")
	}
	//==
}

// * * *

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

// * * *

func resize() {
	var slice []byte

	// Let's allocate the underlying array:
	slice = make([]byte, 4, 5)
	// Checking
	if len(slice) == 4 && cap(slice) == 5 &&
		slice[0] == 0 && slice[1] == 0 && slice[2] == 0 && slice[3] == 0 {
		println("[OK] allocation")
	} else {
		fmt.Println("[Error] allocation")
	}
	//==
	println(fmt.Sprint(slice))

	// Let's change things:
	slice[1], slice[3] = 2, 3
	// Checking
	if slice[0] == 0 && slice[1] == 2 && slice[2] == 0 && slice[3] == 3 {
		println("[OK] change")
	} else {
		fmt.Println("[Error] change")
	}
	//==
	println(fmt.Sprint(slice))

	slice = make([]byte, 2)
	// Checking
	if len(slice) == 2 && cap(slice) == 2 &&
		slice[0] == 0 && slice[1] == 0 {
		println("[OK] resize")
	} else {
		fmt.Println("[Error] resize")
	}
	//==
	println(fmt.Sprint(slice))
}

// * * *

func main() {
	println("\n== valueNil")
	valueNil()
	println("\n== shortHand")
	shortHand()
	println("\n== useFunc")
	useFunc()
	println("\n== reference")
	reference()
	println("\n== resize")
	resize()
}

/*
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
