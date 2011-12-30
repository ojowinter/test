package test

// Anonymous function
var Mul = func(x, y int) int {
	return x * y
}

// Function in the same line.
func hello() { print("Hello world!") }

// Returns the maximum between two int a, and b.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	Add := func(x, y int) int { return x + y }

	x := 3
	y := 4
	z := 5

	max_xy := max(x, y) //calling max(x, y)
	max_xz := max(x, z) //calling max(x, z)
}
