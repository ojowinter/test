package main

import "fmt"

var rating = map[string]float32{"C": 5, "Go": 4.5, "Python": 4.5, "C++": 2}

func valueNil() {
	var n map[string]int

	// Checking
	msg := "declaration"
	if n == nil {
		println("[OK]", msg)
	} else {
		fmt.Println("[Error]", msg)
	}
	//==

	n = make(map[string]int)

	// Checking
	msg = "using make"
	if n != nil {
		println("[OK]", msg)
	} else {
		fmt.Println("[Error]", msg)
	}
	//==
}

func declare_1() {
	// A map that associates strings to int
	// eg. "one" --> 1, "two" --> 2...
	var numbers map[string]int //declare a map of strings to ints
	numbers = make(map[string]int)

	numbers["one"] = 1
	numbers["ten"] = 10
	numbers["trois"] = 3 //trois is "three" in french. I know that you know.

	// Checking
	if numbers["trois"] == 3 {
		println("[OK]")
	} else {
		fmt.Println("[Error] Trois is the french word for the number:", numbers["trois"])
	}
	//==
}

func declare_2() {
	// A map representing the rating given to some programming languages.
	rating2 := map[string]float32{"C": 5, "Go": 4.5, "Python": 4.5, "C++": 2}

	// This is equivalent to writing more verbosely
	rating := make(map[string]float32)
	rating["C"] = 5
	rating["Go"] = 4.5
	rating["Python"] = 4.5
	rating["C++"] = 2 //Linus would put 1 at most. Go ask him

	// Checking
	code := ""
	if rating["Go"] == rating2["Go"] {
		println("[OK] comparing same value")
	} else {
		fmt.Printf("[Error] rating[\"Go\"]: %f\trating2[\"Go\"]: %f\n",
			rating["Go"], rating2["Go"])
	}
	//==

	rating["Go"] = 4.7
	// Checking
	if rating["Go"] != rating2["Go"] {
		code = "OK"
	} else {
		code = "Error"
	}
	println("[" + code + "] comparing different value")
	//==
}

func reference() {
	//let's say a translation dictionary
	m := make(map[string]string)
	m["Hello"] = "Bonjour"

	m1 := m
	m1["Hello"] = "Salut" // Now: m["Hello"] == "Salut"

	// Checking
	if m["Hello"] == m1["Hello"] {
		println("[OK]")
	} else {
		fmt.Println("[Error] value in key:", m["Hello"])
	}
	//==
}

func checkKey() {
	csharp_rating := rating["C#"]
	// Checking
	if csharp_rating == 0.00 {
		println("[OK] single key")
	} else {
		fmt.Println("[Error] value in key:", csharp_rating)
	}
	//==

	multMap := map[int]map[int]string{1: {1: "one"}, 2: {2: "two"}}
	k_multMap := multMap[1][2]
	// Checking
	if k_multMap == "" {
		println("[OK] multi-dimensional key")
	} else {
		fmt.Println("[Error] value in multi-dimensional key:", k_multMap)
	}
	//==

	csharp_rating2, ok := rating["C#"]
	// Checking
	if ok {
		fmt.Println("[Error] using comma")
	} else {
		println("[OK] using comma")
	}
	if csharp_rating2 == 0.00 {
		println("[OK] value (using comma)")
	} else {
		fmt.Println("[Error] value in key (using comma):", csharp_rating2)
	}
	// ==
}

func deleteKey() {
	delete(rating, "C++") // We delete the entry with key "C++"

	_, ok := rating["C++"]
	// Checking
	if ok {
		fmt.Println("[Error]")
	} else {
		println("[OK]")
	}
	// ==
}

func testRange() {
	hasError := false

	// == Iterate over the ratings map
	for key, value := range rating {
		switch key {
		case "C":
			if value != 5 {
				fmt.Println("[Error] key 'C': expected '5', got", value)
				hasError = true
			}
		case "Go":
			if value != 4.5 {
				fmt.Println("[Error] key 'Go': expected '4.5', got", value)
				hasError = true
			}
		case "Python":
			if value != 4.5 {
				fmt.Println("[Error] key 'Python': expected '4.5', got", value)
				hasError = true
			}
		default:
			fmt.Println("[Error] key not expected:", key)
			hasError = true
		}
	}
	if !hasError {
		println("[OK]")
	}

	// == Omit the value.
	for key := range rating {
		if key != "C" && key != "Go" && key != "Python" {
			fmt.Println("[Error] key not expected:", key)
			hasError = true
		}
	}
	if !hasError {
		println("[OK] omitting value")
	}
}

func blankIdentifierInRange() {
	hasError := false

	// Return the biggest value in a slice of ints.
	Max := func(slice []int) int { // The input parameter is a slice of ints.
		max := slice[0]               //the first element is the max for now.
		for _, value := range slice { // Notice how we iterate!
			if value > max { // We found a bigger value in our slice.
				max = value
			}
		}
		return max
	}

	// Declare three arrays of different sizes, to test the function Max.
	A1 := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	A2 := [4]int{1, 2, 3, 4}
	A3 := [1]int{1}

	//declare a slice of ints
	var slice []int

	slice = A1[:] // Take all A1 elements.
	if Max(slice) != 9 {
		fmt.Println("[Error] 'A1': value expected '9', got", Max(slice))
		hasError = true
	}
	slice = A2[:] // Take all A2 elements.
	if Max(slice) != 4 {
		fmt.Println("[Error] 'A2': value expected '4', got", Max(slice))
		hasError = true
	}
	slice = A3[:] // Take all A3 elements.
	if Max(slice) != 1 {
		fmt.Println("[Error] 'A3': value expected '1', got", Max(slice))
		hasError = true
	}

	if !hasError {
		println("[OK]")
	}
}

func main() {
	println("\n== valueNil")
	valueNil()
	println("\n== declare_1")
	declare_1()
	println("\n== declare_2")
	declare_2()
	println("\n== reference")
	reference()
	println("\n== checkKey")
	checkKey()
	println("\n== deleteKey")
	deleteKey()
	println("\n== testRange")
	testRange()
	println("\n== blankIdentifierInRange")
	blankIdentifierInRange()
}
