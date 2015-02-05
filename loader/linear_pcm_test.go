package loader

import (
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
