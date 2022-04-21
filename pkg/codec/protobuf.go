package codec

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

const ProtoBufCodecName = "protobuf"

// ProtoBuf编解码器
type protoBufCodec struct{}

func (*protoBufCodec) Encode(a interface{}) ([]byte, error) {
	if m, ok := a.(proto.Message); ok {
		return proto.Marshal(m)
	}

	return nil, fmt.Errorf("<%T> can't convert to proto.Message", a)
}

func (*protoBufCodec) Decode(data []byte, a interface{}) error {
	if m, ok := a.(proto.Message); ok {
		return proto.Unmarshal(data, m)
	}

	return fmt.Errorf("<%T> can't convert to proto.Message", a)
}
