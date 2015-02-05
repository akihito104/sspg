package loader

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

func LoadDdb(path string) []float64 {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer f.Close()

	b := make([]byte, 4096)
	res := make([]float64, 0, 4096)
	for n, e := f.Read(b); e == nil; n, e = f.Read(b) {
		br := bytes.NewReader(b)
		d := make([]float64, n/8)
		binary.Read(br, binary.LittleEndian, d)
		res = append(res, d...)
	}
	return res
}

func ToIntArr(in []float64, scale float64) []int {
	out := make([]int, len(in))
	for i, a := range in {
		out[i] = int(a * scale)
	}
	return out
}

func ResliceToIntArr(from, to int, in []float64, scale float64) []int {
	return ToIntArr(in[from:to], scale)
}
