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
