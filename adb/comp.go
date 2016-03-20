package adb

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
)

func Enc(rec *Record) ([]byte, error) {
	var buf *bytes.Buffer
	w, err := gzip.NewWriterLevel(buf, 9)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(rec)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(b); err != nil {
		return nil, err
	}
	if err := w.Flush(); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Dec(dat []byte) (*Record, error) {
	return nil, nil
}
