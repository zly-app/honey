package compress

import (
	"io"
)

const RawCompressName = "raw"

// 原始数据, 不进行任何压缩
type RawCompress struct{}

func (r *RawCompress) Compress(in io.Reader, out io.Writer) error {
	_, err := io.Copy(out, in)
	return err
}

func (r *RawCompress) UnCompress(in io.Reader, out io.Writer) error {
	_, err := io.Copy(out, in)
	return err
}

func NewRawCompress() ICompress {
	return &RawCompress{}
}
