// gos project main.go
package main

import (
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

	frame := len(impR30R)
	b := make([]byte, frame*4)
	outArr := make([]int16, frame*2)
	nextArr := make([]int16, len(impR30R)*2)
	for n, e := f.Read(b); e == nil; n, e = f.Read(b) {
		br := bytes.NewReader(b)
		d := make([]int16, n/2)
		if e := binary.Read(br, binary.LittleEndian, d); e != nil {
			fmt.Println("binary.Read: ", e.Error())
		}
		curLen := n / 4
		dR := make([]int16, curLen)
		dL := make([]int16, curLen)
		for i := 0; i < curLen; i++ {
			dR[i] = d[2*i]
			dL[i] = d[2*i+1]
		}
		tmpL := convolve(dR, impR30L)
		tmpR := convolve(dR, impR30R)

		for i := 0; i < curLen; i++ {
			outArr[2*i] = int16(tmpL[i] / 100000)
			outArr[2*i+1] = int16(tmpR[i] / 100000)
		}
		for i, na := range nextArr {
			outArr[i] += na
		}

		if e := binary.Write(out, binary.LittleEndian, outArr); e != nil {
			fmt.Println("binaly.Write: ", e.Error())
		}

		for i := 0; i < len(impR30R)-1; i++ {
			nextArr[2*i] = int16(tmpL[curLen+i] / 100000)
			nextArr[2*i+1] = int16(tmpR[curLen+i] / 100000)
		}
	}
	if e := binary.Write(out, binary.LittleEndian, nextArr); e != nil {
		fmt.Println("binary.Write: ", e.Error())
	}
}

func convolve(sound []int16, imp []int) []int {
	res := make([]int, len(sound)+len(imp)-1)
	for i, s := range sound {
		ss := int(s)
		for j, p := range imp {
			res[i+j] += ss * p
		}
	}
	return res
}

func loadDdbIRes(name string) []int {
	f, err := os.Open(name)
	if err != nil {
		fmt.Print(err.Error())
		return nil
	}
	defer f.Close()

	b := make([]byte, 4096)
	res := make([]float64, 0)
	for n, e := f.Read(b); e == nil; n, e = f.Read(b) {
		br := bytes.NewReader(b)
		d := make([]float64, n/8)
		binary.Read(br, binary.LittleEndian, d)
		res = append(res, d...)
	}

	ires := make([]int, 1400)
	for i, r := range res[190:1590] {
		ires[i] = int(r * 32768 * 4)
	}

	return ires
}

func findPeak(sig []int) (int, float64) {
	max := float64(-1)
	ind := 0
	for i, s := range sig {
		fs := float64(s)
		if ss := fs * fs; max < ss {
			max = ss
			ind = i
		}
	}
	return ind, max
}
