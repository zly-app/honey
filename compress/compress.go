package compress

import (
	"io"

	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"
)

type Compress interface {
	// 压缩
	Compress(in io.Reader, out io.Writer) error
	// 解压缩
	UnCompress(in io.Reader, out io.Writer) error
}

var compresses = map[string]Compress{
	ZStdCompressName: NewZStdCompress(),
}

// 注册压缩程序
func RegistryCompress(name string, c Compress) {
	compresses[name] = c
}

// 获取压缩程序, 压缩程序不存在会panic
func GetCompress(name string) Compress {
	c, ok := compresses[name]
	if !ok {
		logger.Log.Panic("未定义的CompressName", zap.String("name", name))
	}
	return c
}
