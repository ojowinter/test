// Not supported in JavaScript

package test

var (
	i  = MAX_INT_JS + 1
	ui = MAX_UINT_JS + 1
)

// === Chan
var c1 = make(chan int, 10)
var c2 = make(chan bool)
var c3 = <- 0

// === Struct
type i int

type s1 struct {
	a, b int
	c    float64
	_    float32 // padding
	F func()
}

type s2 struct {
	a int64
	i
	f float32
}

func main() {}
