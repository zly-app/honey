package log_data

import (
	"bytes"
	"fmt"

	"go.uber.org/zap/zapcore"
)

type LogData struct {
	T       int64  `json:"t"`                  // 纳秒级时间戳
	Level   string `json:"level"`              // 日志等级
	Msg     string `json:"msg,omitempty"`      // 日志内容
	Fields  string `json:"fields,omitempty"`   // 日志自定义数据
	Line    string `json:"line,omitempty"`     // 日志所在行
	TraceID string `json:"trace_id,omitempty"` // 链路id
}

func MakeLogData(ent *zapcore.Entry, fields []zapcore.Field) *LogData {
	// 解析fields
	enc := zapcore.NewMapObjectEncoder()
	for i := range fields {
		fields[i].AddTo(enc)
	}

	// 提取traceID
	traceID := ""
	if v, ok := enc.Fields["traceID"]; ok {
		traceID = fmt.Sprint(v)
		delete(enc.Fields, "traceID")
	}

	data := &LogData{
		T:       ent.Time.UnixNano(),
		Level:   ent.Level.String(),
		Msg:     ent.Message,
		Fields:  "",
		Line:    fmt.Sprintf("%s:%d,%s", ent.Caller.File, ent.Caller.Line, ent.Caller.Function),
		TraceID: traceID,
	}

	// 序列化fields
	if len(enc.Fields) > 0 {
		var fieldsBuff bytes.Buffer
		for k, v := range enc.Fields {
			fieldsBuff.WriteByte(',')
			fieldsBuff.WriteString(k)
			fieldsBuff.WriteByte('=')
			fieldsBuff.WriteString(fmt.Sprint(v))
		}
		data.Fields = string(fieldsBuff.Bytes()[1:])
	}

	return data
}
