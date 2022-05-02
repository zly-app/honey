package compress

import (
	"compress/gzip"
	"io"
)

const GzipCompressName = "gzip"

type GzipCompress struct{}

func (r *GzipCompress) Compress(in io.Reader, out io.Writer) error {
	w := gzip.NewWriter(out)
	_, err := io.Copy(w, in)
	if err != nil {
		return err
	}
	return w.Close()
}

func (r *GzipCompress) UnCompress(in io.Reader, out io.Writer) error {
	read, err := gzip.NewReader(in)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, read)
	if err != nil {
		return err
	}
	return read.Close()
}

func NewGzipCompress() ICompress {
	return &GzipCompress{}
}
