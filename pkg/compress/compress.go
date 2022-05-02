package compress

import (
	"io"

	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"
)

type ICompress interface {
	// 压缩
	Compress(in io.Reader, out io.Writer) error
	// 解压缩
	UnCompress(in io.Reader, out io.Writer) error
}

var compresses = map[string]ICompress{
	RawCompressName:  NewRawCompress(),
	ZStdCompressName: NewZStdCompress(),
	GzipCompressName: NewGzipCompress(),
}

// 注册压缩程序
func RegistryCompress(name string, c ICompress) {
	if _, ok := compresses[name]; ok {
		logger.Log.Panic("Compress重复注册", zap.String("name", name))
	}
	compresses[name] = c
}

// 获取压缩程序, 压缩程序不存在会panic
func GetCompress(name string) ICompress {
	c, ok := compresses[name]
	if !ok {
		logger.Log.Panic("未定义的CompressName", zap.String("name", name))
	}
	return c
}

// 尝试获取压缩程序
func TryGetCompress(name string) (ICompress, bool) {
	c, ok := compresses[name]
	return c, ok
}
