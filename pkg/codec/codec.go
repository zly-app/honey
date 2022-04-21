package codec

import (
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"
)

// 编解码器
type ICodec interface {
	// 编码
	Encode(a interface{}) ([]byte, error)
	// 解码
	Decode(data []byte, a interface{}) error
}

var codecs = map[string]ICodec{
	JsonCodecName:     &jsonCodec{},
	MsgPackCodecName:  &msgPackCodec{},
	ProtoBufCodecName: &protoBufCodec{},
}

// 注册编解码器
func RegistryCodec(name string, c ICodec) {
	if _, ok := codecs[name]; ok {
		logger.Log.Panic("codec重复注册", zap.String("name", name))
	}
	codecs[name] = c
}

// 获取编解码器
func GetCodec(name string) ICodec {
	c, ok := codecs[name]
	if !ok {
		logger.Log.Panic("试图获取未注册的编解码器", zap.String("name", name))
	}
	return c
}
