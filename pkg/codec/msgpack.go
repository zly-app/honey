package codec

import (
	"bytes"

	"github.com/vmihailenco/msgpack/v5"
)

const MsgPackCodecName = "msgpack"

// MsgPack编解码器
type msgPackCodec struct{}

func (*msgPackCodec) Encode(a interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	enc.SetCustomStructTag("json") // 如果没有 msgpack 标记, 使用 json 标记
	err := enc.Encode(a)
	return buf.Bytes(), err
}

func (*msgPackCodec) Decode(data []byte, a interface{}) error {
	dec := msgpack.NewDecoder(bytes.NewReader(data))
	dec.SetCustomStructTag("json") // 如果没有 msgpack 标记, 使用 json 标记
	err := dec.Decode(a)
	return err
}
