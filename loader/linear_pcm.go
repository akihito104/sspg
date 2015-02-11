package loader

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	riffHeader = "RIFF\x00\x00\x00\x00WAVE"
)

type LnrPcmWav struct {
	ChCount  int16
	SampFreq int32
	File     *os.File
}

func OpenWav(fname string) (wav LnrPcmWav, err error) {
	f, err := os.Open(fname)
	wav = LnrPcmWav{ChCount: int16(0), SampFreq: int32(0), File: nil}
	if err != nil {
		return wav, err
	}
	if !strings.HasSuffix(f.Name(), ".wav") {
		return wav, errors.New(fmt.Sprintf("the file is not .wav: %s", fname))
	}

	rmask := []byte("\xFF\xFF\xFF\xFF\x00\x00\x00\x00\xFF\xFF\xFF\xFF")
	riff := make([]byte, len(rmask))
	if _, e := f.Read(riff); e != nil {
		f.Close()
		return wav, errors.New("read error")
	}
	for i, r := range riff {
		b := r & rmask[i]
		if byte(riffHeader[i]) != b {
			f.Close()
			return wav, errors.New("illegal format: the file is not RIFF WAVE.")
		}
	}

	if e := checkTag(f, "fmt "); e != nil {
		f.Close()
		return wav, e
	}
	var chsize int32
	if e := binary.Read(f, binary.LittleEndian, &chsize); e != nil {
		f.Close()
		return wav, e
	}
	f.Read(make([]byte, 2))
	chc := int16(0)
	binary.Read(f, binary.LittleEndian, &chc)
	fs := int32(0)
	binary.Read(f, binary.LittleEndian, &fs)
	f.Read(make([]byte, chsize-8))

	if e := checkTag(f, "data"); e != nil {
		f.Close()
		return wav, e
	}
	f.Read(make([]byte, 4))

	return LnrPcmWav{ChCount: chc, SampFreq: fs, File: f}, nil
}

func Create(chCount int16, fs int32, fname string) (wav LnrPcmWav, err error) {
	f, err := os.Create(fname)
	f.Write([]byte(riffHeader))

	f.Write([]byte("fmt "))
	f.Write([]byte("\x10")) // length of fmt chunk (bytes)
	f.Write([]byte("\x01")) // format id (linear pcm)
	binary.Write(f, binary.LittleEndian, chCount)
	binary.Write(f, binary.LittleEndian, fs)
	binary.Write(f, binary.LittleEndian, int32(2*int32(chCount)*fs))
	f.Write([]byte("\x10")) // bit/sample

	f.Write([]byte("data"))
	f.Write([]byte("\x00")) // all of sound data length filled at called when Close
	return LnrPcmWav{ChCount: chCount, SampFreq: fs, File: f}, err
}

func (w *LnrPcmWav) Close() error {
	if w.File != nil {
		return w.File.Close()
	}
	return nil
}

func checkTag(f *os.File, tag string) error {
	b := make([]byte, len(tag))
	if _, e := f.Read(b); e != nil {
		return e
	}
	if !equals(tag, b) {
		return errors.New("illegal format: " + tag + "chunk is not found.")
	}
	return nil
}

func equals(tag string, b []byte) bool {
	for i, t := range []byte(tag) {
		if b[i] != t {
			return false
		}
	}
	return true
}
