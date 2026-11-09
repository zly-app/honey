package log_data

import (
	"errors"
	"testing"
	"time"

	"go.uber.org/zap/zapcore"
)

func TestToString(t *testing.T) {
	cases := []struct {
		name string
		in   interface{}
		want string
	}{
		{"string", "hello", "hello"},
		{"error", errors.New("boom"), "boom"},
		{"int", 42, "42"},
		{"nil", nil, "<nil>"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := ToString(c.in); got != c.want {
				t.Fatalf("ToString(%v) = %q, 期望 %q", c.in, got, c.want)
			}
		})
	}
}

func TestMakeLogData(t *testing.T) {
	ent := &zapcore.Entry{
		Level:   zapcore.InfoLevel,
		Time:    time.Unix(1700000000, 123),
		Message: "msg",
		Caller:  zapcore.EntryCaller{File: "a.go", Line: 12, Function: "pkg.F"},
	}
	data := MakeLogData(ent, nil)
	if data == nil {
		t.Fatal("MakeLogData 返回 nil")
	}
	if data.Level != "info" || data.Msg != "msg" {
		t.Fatalf("基础字段不符: %+v", data)
	}
	if data.T != ent.Time.UnixNano() {
		t.Fatalf("时间戳不符: %d", data.T)
	}
	if data.Line != "a.go:12" || data.Func != "pkg.F" {
		t.Fatalf("调用信息不符: %q/%q", data.Line, data.Func)
	}
}
