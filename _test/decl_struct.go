// http://golang.org/doc/go_spec.html#Type_declarations

package test

// An empty struct.
type s0 struct{}

// A struct with 6 fields.
type s1 struct {
	a, b int
	f    float32
	_    float32 // padding
	A    *[]int
	//F    func()
}

// The tag strings define the protocol buffer field numbers.
type s2 struct {
	microsec  uint32 "field 1"
	serverIP6 uint32 "field 2"
	process   string "field 3"
}

//type IntArray [16]int

type (
	Point struct{ x, y float64 }
	//Polar Point
)

func main() {
	type Fa struct {
		a int
	}
}
