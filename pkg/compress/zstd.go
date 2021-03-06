package compress

import (
	"io"

	"github.com/klauspost/compress/zstd"
)

const ZStdCompressName = "zstd"

type ZStdCompress struct{}

func (Z *ZStdCompress) Compress(in io.Reader, out io.Writer) error {
	opts := []zstd.EOption{
		zstd.WithEncoderLevel(zstd.SpeedFastest), // 最快压缩
	}
	enc, err := zstd.NewWriter(out, opts...)
	if err != nil {
		return err
	}
	_, err = io.Copy(enc, in)
	if err != nil {
		_ = enc.Close()
		return err
	}
	return enc.Close()
}

func (Z *ZStdCompress) UnCompress(in io.Reader, out io.Writer) error {
	d, err := zstd.NewReader(in)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(out, d)
	return err
}

func NewZStdCompress() ICompress {
	return &ZStdCompress{}
}
