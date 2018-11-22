// Services > Compress
// This module handless the compression and uncompression of simple strings.

package compress

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func Zip(input string) []byte {
	b := make([]byte, 0, len(input))
	buf := bytes.NewBuffer(b)
	w := gzip.NewWriter(buf)
	w.Write([]byte(input))
	w.Close()
	return buf.Bytes()
}

func Unzip(input []byte) (string, error) {
	// This is where we deal with possible NULL values. If so, just gate it out.
	if len(input) == 0 {
		return "", nil
	}
	r, err := gzip.NewReader(bytes.NewReader(input))
	if err != nil {
		return "", err
	}
	defer r.Close()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
