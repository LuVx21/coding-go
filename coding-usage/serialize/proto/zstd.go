package serialize

import "github.com/klauspost/compress/zstd"

func compressZstd(data []byte) ([]byte, error) {
	encoder, err := zstd.NewWriter(nil)
	if err != nil {
		return nil, err
	}
	defer encoder.Close()

	return encoder.EncodeAll(data, nil), nil
}

func decompressZstd(data []byte) ([]byte, error) {
	decoder, err := zstd.NewReader(nil)
	if err != nil {
		return nil, err
	}
	defer decoder.Close()

	return decoder.DecodeAll(data, nil)
}
