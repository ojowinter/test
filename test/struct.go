package test

// An empty struct.
type s0 struct {}

// A struct with 6 fields.
type s1 struct {
	a, b int
	f    float32
	_    float32 // padding
	A *[]int
	//F func()
}

func main() {}
