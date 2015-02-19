package loader

import (
	"bytes"
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
	Flag     int
}

func OpenWav(fname string) (wav LnrPcmWav, err error) {
	f, err := os.Open(fname)
	wav = LnrPcmWav{ChCount: int16(0), SampFreq: int32(0), File: nil, Flag: os.O_RDONLY}
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

	return LnrPcmWav{ChCount: chc, SampFreq: fs, File: f, Flag: os.O_RDONLY}, nil
}

func Create(chCount int16, fs int32, fname string) (wav LnrPcmWav, err error) {
	f, err := os.Create(fname)
	f.Write([]byte(riffHeader))

	f.Write([]byte("fmt "))
	binary.Write(f, binary.LittleEndian, int32(16))                  // length of fmt chunk (bytes)
	binary.Write(f, binary.LittleEndian, int16(1))                   // format id (linear pcm)
	binary.Write(f, binary.LittleEndian, chCount)                    // channle count
	binary.Write(f, binary.LittleEndian, fs)                         // sampling frequency (Hz)
	binary.Write(f, binary.LittleEndian, int32(2*int32(chCount)*fs)) // data speed (bytes/sec)
	binary.Write(f, binary.LittleEndian, int16(4))                   // bytes/sample
	binary.Write(f, binary.LittleEndian, int16(16))                  // quantity size

	f.Write([]byte("data"))
	binary.Write(f, binary.LittleEndian, int32(0)) // all of sound data length filled at called when Close
	return LnrPcmWav{ChCount: chCount, SampFreq: fs, File: f, Flag: os.O_RDWR}, err
}

func (w *LnrPcmWav) Close() error {
	if w.File == nil {
		return errors.New("File is nil.")
	}
	if w.Flag == os.O_RDWR {
		fi, _ := w.File.Stat()
		size := fi.Size()
		if _, e := w.File.WriteAt(toBytes(int32(size-8)), 4); e != nil {
			return e
		}
		w.File.WriteAt(toBytes(int32(size-44)), 40)
	}
	return w.File.Close()
}

func toBytes(num int32) []byte {
	res := make([]byte, 4)
	w := bytes.NewBuffer(res)
	binary.Write(w, binary.LittleEndian, int32(num))
	return w.Bytes()[4:]
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
