package main

import (
	"fmt"
	"testing"

	"github.com/golang/snappy"
	"github.com/google/brotli/go/cbrotli"
	"github.com/klauspost/compress/zstd"
)

var encoder, _ = zstd.NewWriter(nil)
var decoder, _ = zstd.NewReader(nil)

var originalData = []byte("aabbccddaabbccddaabbccddaabbccddaabbccddaabbccddaabbccddaabbccdd")

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
func Test_brotli_00(t *testing.T) {
	bytes, _ := cbrotli.Encode(originalData, cbrotli.WriterOptions{Quality: 11})
	compressed, _ := cbrotli.Decode(bytes)
	fmt.Printf("压缩: %d bytes -> %d bytes (%.1f%%)\n", len(bytes), len(compressed), float64(len(compressed))/float64(len(bytes))*100)
}
func Test_snappy_00(t *testing.T) {
	bytes := snappy.Encode(nil, originalData)
	compressed, _ := snappy.Decode(nil, bytes)
	fmt.Printf("压缩: %d bytes -> %d bytes (%.1f%%)\n", len(bytes), len(compressed), float64(len(compressed))/float64(len(bytes))*100)
}
