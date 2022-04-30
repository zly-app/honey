package serializer

import (
	"io"

	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"
)

// 序列化器
type ISerializer interface {
	// 序列化
	Marshal(a interface{}, w io.Writer) error
	// 反序列化
	Unmarshal(in io.Reader, a interface{}) error
}

var serializers = map[string]ISerializer{
	JsonSerializerName:    &jsonSerializer{},
	MsgPackSerializerName: &msgPackSerializer{},
}

// 注册序列化器
func RegistrySerializer(name string, c ISerializer) {
	if _, ok := serializers[name]; ok {
		logger.Log.Panic("Serializer重复注册", zap.String("name", name))
	}
	serializers[name] = c
}

// 获取序列化器
func GetSerializer(name string) ISerializer {
	c, ok := serializers[name]
	if !ok {
		logger.Log.Panic("试图获取未注册的序列化器", zap.String("name", name))
	}
	return c
}
