// gos project main.go
package main

import (
	"fmt"
	"github.com/akihito104/sspg/dsp"
	"github.com/akihito104/sspg/loader"
	"os"
	"runtime/pprof"
)

func main() {
	if len(os.Args[1:]) != 1 {
		fmt.Println("Usage: ", os.Args[0], " filepath")
		return
	}
	fname := os.Args[1]
	pproff, _ := os.Create("convo.cpuprofile")
	pprof.StartCPUProfile(pproff)
	defer pprof.StopCPUProfile()

	f, err := loader.OpenWav(fname)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	defer f.Close()

	out, err := loader.Create(2, 44100, "out.wav")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer out.Close()

	impR30R := loadDdbIRes("../resources/main/impR30R_44100.DDB")
	if impR30R == nil {
		return
	}
	impR30L := loadDdbIRes("../resources/main/impR30L_44100.DDB")
	if impR30L == nil {
		return
	}
	impL30R := loadDdbIRes("../resources/main/impL30R_44100.DDB")
	if impL30R == nil {
		return
	}
	impL30L := loadDdbIRes("../resources/main/impL30L_44100.DDB")
	if impL30L == nil {
		return
	}

	frame := len(impR30R)
	d := make([]int16, frame*2)
	nextArr := make([]int16, len(impR30R)*2)
	for n, e := f.Read(d); e == nil; n, e = f.Read(d) {
		curLen := n / 2
		dR := make([]int16, curLen)
		dL := make([]int16, curLen)
		for i := 0; i < curLen; i++ {
			dR[i] = d[2*i]
			dL[i] = d[2*i+1]
		}

		chR := convoCh(dR, impR30R, impR30L)
		chL := convoCh(dL, impL30R, impL30L)
		outArr := make([]int16, len(chR))
		add(outArr, nextArr)
		add(outArr, chR)
		add(outArr, chL)

		if _, e := out.Write(outArr[:curLen*2]); e != nil {
			fmt.Println("binaly.Write: ", e.Error())
		}
		nextArr = outArr[curLen*2:]
	}
	if _, e := out.Write(nextArr); e != nil {
		fmt.Println("binary.Write: ", e.Error())
	}
}

func add(out, adder []int16) {
	for i, a := range adder {
		out[i] += a
	}
}

func convoCh(s []int16, iR []int32, iL []int32) []int16 {
	tmpR := dsp.Convolve(s, iR)
	tmpL := dsp.Convolve(s, iL)
	outLen := len(tmpR) * 2
	outArr := make([]int16, outLen)
	for i, a := range tmpR {
		outArr[2*i] = int16(a / 100000)
	}
	for i, a := range tmpL {
		outArr[2*i+1] = int16(a / 100000)
	}
	return outArr
}

func loadDdbIRes(name string) []int32 {
	res := loader.LoadDdb(name)
	ires := loader.ResliceToIntArr(190, 1590, res, 32768*4)
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
