// Not supported in JavaScript

package test

// Chan
var c1 = make(chan int, 10)
var c2 = make(chan bool)
var c3 = <- 0

// Struct
type s1 struct {
	a, b int
	f    float64
	_    float32 // padding
	F func()
}

type i int
type s2 struct {
	a int64
	i
	f float32
}

func main() {}
