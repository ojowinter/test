package main

func testIf() {
	x := 5

	// Simple
	if x > 10 {
		fmt.Println("x is greater than 10")
	} else {
		fmt.Println("x is less than 10")
	}

	// Leading initial short
	if x := 12; x > 10 {
		fmt.Println("x is greater than 10")
	} else {
		fmt.Println("x is less than 10")
	}

	// Multiple if/else
	i := 7

	if i == 3 {
		fmt.Println("i is equal to 3")
	} else if i < 3 {
		fmt.Println("i is less than 3")
	} else {
		fmt.Println("i is greater than 3")
	}
}

func testSwitch() string {
	i := 10

	// Simple
	switch i {
	case 1:
		fmt.Println("i is equal to 1")
	case 2, 3, 4:
		fmt.Println("i is equal to 2, 3 or 4")
	case 10:
		fmt.Println("i is equal to 10")
	default:
		fmt.Println("All I know is that i is an integer")
	}

	// Without expression
	switch {
	case i < 10:
		fmt.Println("i is less than 10")
	case i > 10, i < 0:
		fmt.Println("i is either bigger than 10 or less than 0")
	case i == 10:
		fmt.Println("i is equal to 10")
	default:
		fmt.Println("This won't be printed anyway")
	}

	// With fallthrough
	switch i := 6; {
	case 4:
		fmt.Println("was <= 4")
		fallthrough
	case 5:
		fmt.Println("was <= 5")
		fallthrough
	case 6:
		fmt.Println("was <= 6")
		fallthrough
	case 7:
		fmt.Println("was <= 7")
		fallthrough
	case 8:
		fmt.Println("was <= 8")
		//fallthrough
	default:
		fmt.Println("default case")
	}

	// With return
	switch i {
	default:
	case 1, 3, 5, 7, 9:
		return "odd"
	case 2, 4, 6, 8:
		return "even"
	}

	return ""
}

func testFor() {
	sum := 0

	// Simple
	for i := 0; i < 10; i++ {
		sum += i
	}
	fmt.Println("sum is equal to ", sum)

	// Expression1 and expression3 are omitted here
	sum := 1
	for sum < 1000 {
		sum += sum
	}
	fmt.Println("sum is equal to ", sum)

	// Expression1 and expression3 are omitted here, and semicolons gone
	sum := 1
	for sum < 1000 {
		sum += sum
	}
	fmt.Println("sum is equal to ", sum)

	// Infinite loop, no semicolons at all
	for {
		fmt.Println("I loop for ever!")
	}

	// break
	for i := 10; i > 0; i-- {
		if i < 5 {
			break
		}
		fmt.Println(i)
	}

	// continue
	for i := 10; i > 0; i-- {
		if i == 5 {
			continue
		}
		fmt.Println(i)
	}
}

func testRange() {
	s := []int{2, 3, 5}

	/*for i, v := range s {
		a = s
	}*/
}
