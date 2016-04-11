package rand

import (
	"bufio"
	"encoding/binary"
	"os"
)

var r *bufio.Reader

func Inititalize() {
	f, err := os.Open("/home/_/go/src/github.com/karlek/vanilj/cmd/buddha/rand/source")
	if err != nil {
		panic(err)
	}
	// defer f.Close()
	r = bufio.NewReader(f)
}

func Float64() (ret float64) {
	err := binary.Read(r, binary.LittleEndian, &ret)
	if err != nil {
		panic(err)
	}
	return ret
}

func Sign() float64 {
	var sign uint8
	err := binary.Read(r, binary.LittleEndian, &sign)
	if err != nil {
		panic(err)
	}
	return float64(sign)
}
