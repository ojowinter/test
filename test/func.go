package main

func min(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

func testIf() int {
	x, y := 3, 5
	max := 7

	if x > max {
		x = max
	}

	if z := 2; x < y {
		return x
	} else if x > z {
		return z
	} else {
		return y
	}
}

func testBreak() {
	for i := 10; i > 0; i--{
		if i < 5{
			break
		}
//		fmt.Println("i")
	}
}

func testContinue() {
	for i := 10; i > 0; i--{
		if i == 5{
			continue
		}
//		fmt.Println("i")
	}
}

func testSwitch(tag int) string {
	//var str string
	str := ""

	switch tag {
	default:
	case 1, 3, 5, 7, 9: return "odd"
	case 2, 4, 6, 8: return "even"
	}

	switch x := tag; {
	case x < 0:
		str = "negative"
	default:
		str = "positive"
	}

	y := 100
	switch {
	case x < y: str = "lesser than 100"
	case x > y: str = "greater than 100"
	case x == 0: str = "zero"
	}

	return str

}

func testSwitch2() {
	switch i := 10; {
	case 1:
		fmt.Println("i is equal to 1")
	case 2, 3, 4:
			fmt.Println("i is equal to 2, 3 or 4")
	case 10:
		fmt.Println("i is equal to 10")
	default:
		fmt.Println("All I know is that i is an integer")
	}

	i := 10
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
}

func testFor() int {
	a, b := 1, 2

	for a < b {
		a *= 2
	}

	for i := 0; i < 10; i++ {
		a = i
	}

	s := []int{2, 3, 5}
	/*for i, v := range s {
		a = s
	}*/
}
