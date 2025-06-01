package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

func Decompress(data []byte) ([]byte, error) {
	r, _ := gzip.NewReader(bytes.NewReader(data))
	defer r.Close()

	var b bytes.Buffer
	_, err := b.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}

	return b.Bytes(), nil
}
