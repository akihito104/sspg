// gos project main.go
package main

import (
	_ "bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"runtime/pprof"
)

// D:\ecWork\ohlsample
// imp{L|R}{30|45}{L|R}_44100.DDB

func main() {
	if len(os.Args[1:]) != 1 {
		fmt.Println("Usage: ", os.Args[0], " filepath")
		return
	}
	fname := os.Args[1]
	pproff, _ := os.Create("convo.cpuprofile")
	pprof.StartCPUProfile(pproff)
	defer pprof.StopCPUProfile()

	f, err := os.Open(fname)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	defer f.Close()

	out, err := os.OpenFile("out.pcm", os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer out.Close()

	// todo: convert to DSB format
	impR30R := loadDdbIRes("D:\\ecWork\\ohlsample\\impR30R_44100.DDB")
	impR30L := loadDdbIRes("D:\\ecWork\\ohlsample\\impR30L_44100.DDB")
	//impL30R := loadDdbIRes("D:\\ecWork\\ohlsample\\impL30R_44100.DDB")
	//impL30L := loadDdbIRes("D:\\ecWork\\ohlsample\\impL30L_44100.DDB")

	size := 0
	b := make([]byte, len(impR30R)*4)
	outArr := make([]int16, len(b)/2)
	for n, e := f.Read(b); e == nil; n, e = f.Read(b) {
		size = size + n
		br := bytes.NewReader(b)
		d := make([]int16, n/2)
		if e := binary.Read(br, binary.LittleEndian, d); e != nil {
			fmt.Println("binary.Read: ", e.Error())
		}
		dR := make([]int16, len(d)/2)
		dL := make([]int16, len(dR))
		for i := 0; i < len(dR); i++ {
			dR[i] = d[2*i]
			dL[i] = d[2*i+1]
		}
		tmpL := convolve(dR, impR30L)
		tmpR := convolve(dR, impR30R)
		for i := 0; i < len(dR); i++ {
			outArr[2*i] = tmpL[i]
			outArr[2*i+1] = tmpR[i]
		}

		if e := binary.Write(out, binary.LittleEndian, outArr); e != nil {
			fmt.Println("binaly.Write: ", e.Error())
		}
		fmt.Println("size: ", size)
	}
}
func convolve(sound []int16, imp []float64) []int16 {
	res := make([]int16, len(sound)+len(imp))
	var tmp float64
	for i, s := range sound {
		for j, p := range imp {
			tmp = float64(s) * p
			res[i+j] = res[i+j] + int16(tmp*p)
		}
	}
	return res
}

func loadDdbIRes(name string) []float64 {
	f, err := os.Open(name)
	if err != nil {
		fmt.Print(err.Error())
		return nil
	}
	defer f.Close()

	b := make([]byte, 4096)
	res := make([]float64, 4096)
	for n, e := f.Read(b); e == nil; n, e = f.Read(b) {
		br := bytes.NewReader(b)
		d := make([]float64, n/8)
		binary.Read(br, binary.LittleEndian, d)
		res = append(res, d...)
	}
	return res
}
