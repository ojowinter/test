package main

import "fmt"

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
	rating := map[string]float32 {"C":5, "Go":4.5, "Python":4.5, "C++":2 }
	csharp_rating := rating["C#"]

	// Checking
	if csharp_rating == 0.00 {
		println("[OK]")
	} else {
		fmt.Println("[Error] value in key:", csharp_rating)
	}
	//==
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
}
