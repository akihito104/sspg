// gos project main.go
package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

func main() {
	if len(os.Args[1:]) != 1 {
		fmt.Println("Usage: ", os.Args[0], " filepath")
		return
	}
	fname := os.Args[1]

	f, err := os.Open(fname)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	size := 0
	d := make([]int16, 4096)
	for binary.Read(f, binary.LittleEndian, d) == nil {
		size = size + len(d)
		fmt.Println("size: ", size)
	}

	if e := f.Close(); e != nil {
		fmt.Println(e.Error())
	}
}
