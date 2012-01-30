package main

func declaration() {
	// A map that associates strings to int
	// eg. "one" --> 1, "two" --> 2...
	var numbers map[string]int //declare a map of strings to ints

	numbers = make(map[string]int)
	numbers["one"] = 1
	numbers["ten"] = 10
	numbers["trois"] = 3 //trois is "three" in french. I know that you know.
	//...
	println("Trois is the french word for the number:", numbers[3])
	// Trois is the french word for the number: 3. Also a good time.
}

func main() {
	println("\n== declaration()\n")
	declaration()
}

/*

*/
