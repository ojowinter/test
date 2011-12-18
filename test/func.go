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

func testSwitch(tag int) string {
	//var str string
	str := ""

	switch tag {
	default: return "zero"
	case 1, 3, 5, 7, 9: return "odd"
	case 2, 4, 6, 8: return "even"
	}

	switch x := tag; {
	case x < 0:
		str = "negative"
	default:
		str = "positive"
	}

/*	switch {
	case x < y: f1()
	case x < z: f2()
	case x == 4: f3()
	}*/
}
