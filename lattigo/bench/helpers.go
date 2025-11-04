package main

import (
	"bytes"
	"compress/zlib"

	"github.com/klauspost/compress/zstd"
)

// compressZLIB сжимает данные используя ZLIB
func compressZLIB(data []byte) []byte {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	w.Write(data)
	w.Close()
	return buf.Bytes()
}

// compressZstd сжимает данные используя Zstandard
func compressZstd(data []byte) []byte {
	encoder, _ := zstd.NewWriter(nil)
	return encoder.EncodeAll(data, make([]byte, 0, len(data)))
}
