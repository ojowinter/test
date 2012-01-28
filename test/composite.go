package main

import "fmt"

// We declare our new type
type person struct {
	name string
	age  int
}

// Return the older person of p1 and p2, and the difference in their ages.
func Older(p1, p2 person) (person, int) {
	if p1.age > p2.age { // Compare p1 and p2's ages
		return p1, p1.age - p2.age
	}
	return p2, p2.age - p1.age
}

func testStruct() {
	var tom person

	tom.name, tom.age = "Tom", 18

	// Look how to declare and initialize easily.
	bob := person{age: 25, name: "Bob"} //specify the fields and their values
	paul := person{"Paul", 43}          //specify values of fields in their order

	tb_Older, tb_diff := Older(tom, bob)
	tp_Older, tp_diff := Older(tom, paul)
	bp_Older, bp_diff := Older(bob, paul)

	fmt.Printf("Of %s and %s, %s is older by %d years\n",
		tom.name, bob.name, tb_Older.name, tb_diff)

	fmt.Printf("Of %s and %s, %s is older by %d years\n",
		tom.name, paul.name, tp_Older.name, tp_diff)

	fmt.Printf("Of %s and %s, %s is older by %d years\n",
		bob.name, paul.name, bp_Older.name, bp_diff)
}

// Return the older person in a group of 10 persons.
func Older10(people [10]person) person {
	older := people[0] // The first one is the older for now.

	// Loop through the array and check if we could find an older person.
	for index := 1; index < 10; index++ { // We skipped the first element here.
		if people[index].age > older.age { // Current's persons age vs olderest so far.
			older = people[index] // If people[index] is older, replace the value of older.
		}
	}
	return older
}

func testArray() {
	// Declare an example array variable of 10 person called 'array'.
	var array [10]person

	// Initialize some of the elements of the array, the others are by default
	// set to person{"", 0}
	array[1] = person{"Paul", 23}
	array[2] = person{"Jim", 24}
	array[3] = person{"Sam", 84}
	array[4] = person{"Rob", 54}
	array[8] = person{"Karl", 19}

	older := Older10(array) // Call the function by passing it our array.

	println("The older of the group is:", older.name)
}

func initializeArray() {
	// Declare and initialize an array A of 10 person.
	array1 := [10]person{
		person{"", 0},
		person{"Paul", 23},
		person{"Jim", 24},
		person{"Sam", 84},
		person{"Rob", 54},
		person{"", 0},
		person{"", 0},
		person{"", 0},
		person{"Karl", 10},
		person{"", 0},
	}

	// Declare and initialize an array of 10 persons, but let the compiler guess the size.
	array2 := [...]person{ // Substitute '...' instead of an integer size.
		person{"", 0},
		person{"Paul", 23},
		person{"Jim", 24},
		person{"Sam", 84},
		person{"Rob", 54},
		person{"", 0},
		person{"", 0},
		person{"", 0},
		person{"Karl", 10},
		person{"", 0}}

	if len(array1) == len(array2) {
		print("array1 and array2 have the same length: ")
	}
	if len(array1) == 10 {
		println("10")
	}
	if array1 == array2 {
		println("array1 and array2 are equals")
	}
}

func multiArray() {
	// declare and initialize an array of 2 arrays of 4 ints
	doubleArray_1 := [2][4]int{[4]int{1, 2, 3, 4}, [4]int{5, 6, 7, 8}}

	// simplify the previous declaration, with the '...' syntax
	doubleArray_2 := [2][4]int{
		[...]int{1, 2, 3, 4}, [...]int{5, 6, 7, 8}}

	// super simpification!
	doubleArray_3 := [2][4]int{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
	}

	if doubleArray_1 == doubleArray_2 && doubleArray_2 == doubleArray_3 {
		println("The three multi-dimensional arrays are equal")
	}
}

func main() {
	println("\n== testStruct()\n")
	testStruct()

	println("\n== testArray()\n")
	testArray()

	println("\n== initializeArray()\n")
	initializeArray()

	println("\n== multiArray()\n")
	multiArray()
}

/*
== testStruct()

Of Tom and Bob, Bob is older by 7 years
Of Tom and Paul, Paul is older by 25 years
Of Bob and Paul, Paul is older by 18 years

== testArray()

The older of the group is: Sam

== initializeArray()

array1 and array2 have the same length: 10
array1 and array2 are equals

== multiArray()

The three multi-dimensional arrays are equal
*/
