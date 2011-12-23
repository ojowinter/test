package main

var x, y = 14, 9

var (
	and = x & y
	or = x | y
	xor = x ^ y
	not = !y

	lShift = 9 << 2
	rShift = 9 >> 2
	rShiftNeg = -9 >> 2

	zRShift = 9 >>> 2
	zRShiftNeg = -9 >>> 2
)
