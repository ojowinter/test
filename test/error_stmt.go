// Statements not supported

package test

func ret() (byte, byte) {
	ch := make(chan int)

	go print("hello!")
	defer println("bye!")

	panic("problem")
	recover()

	return 0, 0
}

func testGoto() {
	isFirst := true

_skipPoint:
	println("Using label")

	if isFirst {
		isFirst = false
		goto _skipPoint
		print("This part is skipped")
	}
}
