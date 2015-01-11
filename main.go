// gos project main.go
package main

import (
	_ "bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

// D:\ecWork\ohlsample
// imp{L|R}{30|45}{L|R}_44100.DDB

func main() {
	if len(os.Args[1:]) != 1 {
		fmt.Println("Usage: ", os.Args[0], " filepath")
		return
	}
	fname := os.Args[1]
	impR30R := loadDdbIRes("D:\\ecWork\\ohlsample\\impR30R_44100.DDB")
	fmt.Println("impR30R_44100.DDB length: ", len(impR30R))

	f, err := os.Open(fname)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	size := 0
	d := make([]int16, 4096)
	for binary.Read(f, binary.LittleEndian, d) == nil {
		size = size + len(d)
	}
	fmt.Println("size: ", size)

	if e := f.Close(); e != nil {
		fmt.Println(e.Error())
	}
}

func loadDdbIRes(name string) []float64 {
	f, err := os.Open(name)
	if err != nil {
		fmt.Print(err.Error())
		return nil
	}
	b := make([]byte, 4096)
	res := make([]float64, 4096)
	for n, e := f.Read(b); e == nil; n, e = f.Read(b) {
		br := bytes.NewReader(b)
		d := make([]float64, n/8)
		binary.Read(br, binary.LittleEndian, d)
		res = append(res, d...)
	}
	if e := f.Close(); e != nil {
		fmt.Println(e.Error())
	}
	return res
}
