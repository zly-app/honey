package serializer

import (
	"io"

	jsoniter "github.com/json-iterator/go"
)

const JsonSerializerName = "json"

// 使用第三方包json-iterator进行序列化
type jsonSerializer struct{}

func (*jsonSerializer) Marshal(a interface{}, w io.Writer) error {
	return jsoniter.NewEncoder(w).Encode(a)
}

func (*jsonSerializer) Unmarshal(in io.Reader, a interface{}) error {
	return jsoniter.NewDecoder(in).Decode(a)
}
