package test

var n int
var U, V, W float64
var k = 0
var x, y float32 = -1, -2
var (
	j        int
	u, v, st = 2.0, 3.0, "bar"
)
//var re, im = complexSqrt(-1)
//var _, found = entries[name] // map lookup; only interested in "found"

var b = true   // t has type bool
var i = 0      // i has type int
var f = 3.0    // f has type float64
var s = "OMDB" // s has type string

// Array
var (
	a1 = new([32]byte)
	a2 = new([2][4]uint)
	//a3 = [2*N] struct { x, y int32 }
	a4 = [1000]*float64{}
	a5 = [4]byte{}
	a6 = [3][5]int{}
	a7 = [2][2][2]float64{} // same as [2]([2]([2]float64))
)

func main() {
	// === Short variable declarations
	i, j := 0, 10
	f := func() int { return 7 }
	ch := make(chan int)
	r, w := os.Pipe(fd) // os.Pipe() returns two values
	_, y, _ := coord(p) // coord() returns three values; only interested in y coordinate

	field1, offset := nextField(str, 0)
	field2, offset := nextField(str, offset) // redeclares offset
}
