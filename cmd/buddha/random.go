package main

import (
	"io"
)

var p = make([]byte, 1)

func sign(random io.Reader) float64 {
	random.Read(p)
	r := int(p[0] % 2)
	if r == 1 {
		return -1.0
	}
	return 1.0
}

var p1 = make([]byte, 4)

func randfloat(random io.Reader) float64 {
	random.Read(p1)
	b0, b1, b2, b3 := float64(p1[0]), float64(p1[1]), float64(p1[2]), float64(p1[3])
	return (1 / 256.0) * (b0 + (1/256.0)*(b1+(1/256.0)*(b2+(1/256.0)*b3)))
}
