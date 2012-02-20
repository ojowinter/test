package main

import "fmt"

// Our struct representing a person
type person struct {
	name string
	age  int
}

// Return true, and the older person in a group of persons
// Or false, and nil if the group is empty.
func Older(people ...person) (bool, person) { // Variadic function.
	if len(people) == 0 {
		return false, person{}
	} // The group is empty.
	older := people[0] // The first one is the older FOR NOW.
	// Loop through the slice people.
	for _, value := range people { // We don't need the index.
		// Compare the current person's age with the oldest one so far
		if value.age > older.age {
			older = value //if value is older, replace older
		}
	}
	return true, older
}

func main() {

	// Two variables to be used by our program.
	var (
		ok    bool
		older person
	)

	// Declare some persons.
	paul := person{"Paul", 23}
	jim := person{"Jim", 24}
	sam := person{"Sam", 84}
	rob := person{"Rob", 54}
	karl := person{"Karl", 19}

	// Who is older? Paul or Jim?
	_, older = Older(paul, jim) //notice how we used the blank identifier
	fmt.Println("The older of Paul and Jim is: ", older.name)
	// Who is older? Paul, Jim or Sam?
	_, older = Older(paul, jim, sam)
	fmt.Println("The older of Paul, Jim and Sam is: ", older.name)
	// Who is older? Paul, Jim , Sam or Rob?
	_, older = Older(paul, jim, sam, rob)
	fmt.Println("The older of Paul, Jim, Sam and Rob is: ", older.name)
	// Who is older in a group containing only Karl?
	_, older = Older(karl)
	fmt.Println("When Karl is alone in a group, the older is: ", older.name)
	// Is there an older person in an empty group?
	ok, older = Older() //this time we use the boolean variable ok
	if !ok {
		fmt.Println("In an empty group there is no older person")
	}
}

// * * *
