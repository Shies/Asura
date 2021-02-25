package dao

import (
	"bytes"
	"compress/zlib"
	"io"
)

func CompressEncode(value string) string {
	var in bytes.Buffer
	b := []byte(value)
	w := zlib.NewWriter(&in)
	w.Write(b)
	w.Close()

	return in.String()
}

func CompressDecode(value string) string {
	var out bytes.Buffer
	b := []byte(value)
	in := bytes.NewReader(b)
	r, _ := zlib.NewReader(in)
	io.Copy(&out, r)

	return out.String()
}