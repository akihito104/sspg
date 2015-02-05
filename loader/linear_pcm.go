package loader

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strings"
)

type LnrPcmWav struct {
	ChCount  int16
	SampFreq int
	File     *os.File
}

func OpenWav(fname string) (wav LnrPcmWav, err error) {
	f, err := os.Open(fname)
	wav = LnrPcmWav{ChCount: int16(0), SampFreq: int(0), File: nil}
	if err != nil {
		return wav, err
	}
	if !strings.HasSuffix(f.Name(), ".wav") {
		return wav, errors.New(fmt.Sprintf("the file is not .wav: %s", fname))
	}

	rmask := []byte("\xFF\xFF\xFF\xFF\x00\x00\x00\x00\xFF\xFF\xFF\xFF")
	rpat := []byte("RIFF\x00\x00\x00\x00WAVE")
	riff := make([]byte, len(rmask))
	if _, e := f.Read(riff); e != nil {
		f.Close()
		return wav, errors.New("read error")
	}
	for i, r := range riff {
		b := r & rmask[i]
		if rpat[i] != b {
			f.Close()
			return wav, errors.New("illegal format: the file is not RIFF WAVE.")
		}
	}
	tag := make([]byte, 4)

	if _, e := f.Read(tag); e != nil {
		f.Close()
		return wav, e
	}
	if !compareTag("fmt ", tag) {
		f.Close()
		return wav, errors.New("illegal format: fmt chunk is not found.")
	}

	f.Read(make([]byte, 2))
	var chsize int32
	err = binary.Read(f, binary.LittleEndian, &chsize)
	if err != nil {
		f.Close()
		return wav, err
	}
	chc := int16(0)
	binary.Read(f, binary.LittleEndian, &chc)
	fs := int32(0)
	binary.Read(f, binary.LittleEndian, &fs)
	f.Read(make([]byte, 4+2+2))

	if _, e := f.Read(tag); e != nil {
		f.Close()
		return wav, e
	}
	if !compareTag("data", tag) {
		f.Close()
		return wav, errors.New("illegal format: data chunk is not found.")
	}
	f.Read(make([]byte, 4))

	wav.File = f
	return LnrPcmWav{ChCount: chc, SampFreq: int(fs), File: f}, nil
}

func compareTag(tag string, b []byte) bool {
	for i, t := range []byte(tag) {
		if b[i] != t {
			return false
		}
	}
	return true
}

func (w *LnrPcmWav) Close() error {
	if w.File != nil {
		return w.File.Close()
	}
	return nil
}
