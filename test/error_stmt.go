// Directives not supported

package test

func unvalidDirectives() {
	ch := make(chan int)

	go print("hello!")
	defer println("bye!")

	panic("problem")
	recover()
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
