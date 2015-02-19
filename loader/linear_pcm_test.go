package loader

import (
	"encoding/binary"
	"fmt"
	"testing"
)

func TestOpenWav(t *testing.T) {
	fname := "../resources/test/yamanosususme_ss/natsuiro/01-AudioTrack.wav"
	wav, err := OpenWav(fname)
	if wav.File != nil {
		defer wav.Close()
	}

	if err != nil {
		t.Error("got error:", err.Error())
		fmt.Println(wav)
	}
	if wav.File == nil {
		t.Error("wav.File should have ", fname)
	} else if wav.File.Name() != fname {
		t.Error("wav.File should be", fname, ", but", wav.File.Name())
	}
	if wav.SampFreq != 44100 {
		t.Error("wav.SampFreq should be 44100, but", wav.SampFreq)
	}
	if wav.ChCount != 2 {
		t.Error("wav.ChCount should be 2, but", wav.ChCount)
	}
}

func TestCreateAndClose(t *testing.T) {
	wav, err := Create(2, 44100, "test.wav")
	if err != nil {
		t.Error(err.Error())
		return
	}
	for i := 0; i < 10; i++ {
		binary.Write(wav.File, binary.LittleEndian, make([]int16, 10))
	}
	if e := wav.Close(); e != nil {
		t.Error(e.Error())
	}

	act, err := OpenWav("test.wav")
	if err != nil {
		t.Error(err.Error())
	}
	defer act.File.Close()
	fi, _ := act.File.Stat()
	if fi.Size() != (44 + 10*10*2) {
		t.Error("size: ", fi.Size())
	}
	if act.ChCount != 2 {
		t.Error("chcount: ", act.ChCount)
	}
	if act.SampFreq != 44100 {
		t.Error("sample freq.: ", act.SampFreq)
	}
}
