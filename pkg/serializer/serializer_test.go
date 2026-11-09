package serializer

import (
	"bytes"
	"testing"
)

type kv struct {
	A string `json:"a" msgpack:"a"`
	B int    `json:"b" msgpack:"b"`
}

func TestBuiltinSerializersRegistered(t *testing.T) {
	for _, name := range []string{JsonSerializerName, MsgPackSerializerName} {
		s, ok := TryGetSerializer(name)
		if !ok || s == nil {
			t.Fatalf("内置序列化器 %q 未注册", name)
		}
	}
}

func TestTryGetSerializerUnknown(t *testing.T) {
	if s, ok := TryGetSerializer("not-exist"); ok || s != nil {
		t.Fatalf("未知序列化器应返回 (nil,false), 实际: %v %v", s, ok)
	}
}

func TestSerializerRoundTrip(t *testing.T) {
	for _, name := range []string{JsonSerializerName, MsgPackSerializerName} {
		t.Run(name, func(t *testing.T) {
			s := GetSerializer(name)
			buf := &bytes.Buffer{}
			in := kv{A: "hello", B: 7}
			if err := s.Marshal(in, buf); err != nil {
				t.Fatalf("Marshal 失败: %v", err)
			}
			var out kv
			if err := s.Unmarshal(bytes.NewReader(buf.Bytes()), &out); err != nil {
				t.Fatalf("Unmarshal 失败: %v", err)
			}
			if out != in {
				t.Fatalf("往返不一致: %+v -> %+v", in, out)
			}
		})
	}
}

func TestRegistrySerializer(t *testing.T) {
	const name = "test-custom"
	if _, ok := TryGetSerializer(name); ok {
		t.Fatalf("测试用名字不应已存在: %s", name)
	}
	RegistrySerializer(name, GetSerializer(JsonSerializerName))
	s, ok := TryGetSerializer(name)
	if !ok || s == nil {
		t.Fatalf("注册后应能取到: %s", name)
	}
}
