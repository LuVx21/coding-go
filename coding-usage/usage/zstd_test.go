package main

import (
	"fmt"
	"testing"

	"github.com/klauspost/compress/zstd"
)

var encoder, _ = zstd.NewWriter(nil)
var decoder, _ = zstd.NewReader(nil)

func compressZstd(bytes []byte) ([]byte, error) {
	return encoder.EncodeAll(bytes, nil), nil
}

func decompressZstd(bytes []byte) ([]byte, error) {
	return decoder.DecodeAll(bytes, nil)
}

func Test_zstd_00(t *testing.T) {
	bytes := []byte("aabbccddaabbccddaabbccddaabbccddaabbccddaabbccddaabbccddaabbccdd")
	compressed, _ := compressZstd(bytes)
	fmt.Printf("压缩: %d bytes -> %d bytes (%.1f%%)\n", len(bytes), len(compressed), float64(len(compressed))/float64(len(bytes))*100)
}
