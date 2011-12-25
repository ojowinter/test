// Statements not supported

package test

func ret() (byte, byte) {
	go print("hello!")
	defer println("bye!")

	panic("problem")
	recover()

	return 0, 0
}
