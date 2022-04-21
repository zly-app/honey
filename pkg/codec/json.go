package codec

import (
	jsoniter "github.com/json-iterator/go"
)

const JsonCodecName = "json"

// 使用第三方包json-iterator进行编解码
type jsonCodec struct{}

func (*jsonCodec) Encode(a interface{}) ([]byte, error) {
	return jsoniter.Marshal(a)
}

func (*jsonCodec) Decode(data []byte, a interface{}) error {
	return jsoniter.Unmarshal(data, a)
}
