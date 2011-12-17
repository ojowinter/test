// Statements not supported

package test

func ret() (byte, byte) {
	go print("hello!")

	return 0, 0
}
