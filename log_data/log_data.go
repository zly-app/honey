package log_data

import (
	"bytes"
	"strconv"

	"github.com/spf13/cast"
	"go.uber.org/zap/zapcore"
)

type LogData struct {
	T       int64  `json:"t"`                // 纳秒级时间戳
	Level   string `json:"level"`            // 日志等级
	Msg     string `json:"msg,omitempty"`    // 日志内容
	Fields  string `json:"fields,omitempty"` // 日志自定义数据
	Line    string `json:"line,omitempty"`   // 日志所在行
	Func    string `json:"func,omitempty"`   // 函数名
	TraceID string `json:"traceID,omitempty"`
	SpanID  string `json:"spanID,omitempty"`

	Duration      int64  `json:"duration,omitempty"` // 持续时间, 纳秒
	Code          int64  `json:"code,omitempty"`     // Code
	CodeType      string `json:"codeType,omitempty"` // CodeType
	CallerService string `json:"callerService,omitempty"`
	CallerMethod  string `json:"callerMethod,omitempty"`
	CalleeService string `json:"calleeService,omitempty"`
	CalleeMethod  string `json:"calleeMethod,omitempty"`

	Req string `json:"Req,omitempty"`
	Rsp string `json:"Rsp,omitempty"`
}

func MakeLogData(ent *zapcore.Entry, fields []zapcore.Field) *LogData {
	// 解析fields
	enc := zapcore.NewMapObjectEncoder()
	for i := range fields {
		fields[i].AddTo(enc)
	}

	data := &LogData{
		T:      ent.Time.UnixNano(),
		Level:  ent.Level.String(),
		Msg:    ent.Message,
		Fields: "",
		Line:   ent.Caller.File + ":" + strconv.Itoa(ent.Caller.Line),
		Func:   ent.Caller.Function,
	}

	// 提取traceID
	if v, ok := enc.Fields["traceID"]; ok {
		data.TraceID = cast.ToString(v)
		delete(enc.Fields, "traceID")
	}

	// 提取spanID
	if v, ok := enc.Fields["spanID"]; ok {
		data.SpanID = cast.ToString(v)
		delete(enc.Fields, "spanID")
	}

	if v, ok := enc.Fields["duration"]; ok {
		data.Duration = cast.ToInt64(v)
		delete(enc.Fields, "duration")
	}

	if v, ok := enc.Fields["code"]; ok {
		data.Code = cast.ToInt64(v)
		delete(enc.Fields, "code")
	}
	if v, ok := enc.Fields["codeType"]; ok {
		data.CodeType = cast.ToString(v)
		delete(enc.Fields, "codeType")
	}
	if v, ok := enc.Fields["callerService"]; ok {
		data.CallerService = cast.ToString(v)
		delete(enc.Fields, "callerService")
	}
	if v, ok := enc.Fields["callerMethod"]; ok {
		data.CallerMethod = cast.ToString(v)
		delete(enc.Fields, "callerMethod")
	}
	if v, ok := enc.Fields["calleeService"]; ok {
		data.CalleeService = cast.ToString(v)
		delete(enc.Fields, "calleeService")
	}
	if v, ok := enc.Fields["calleeMethod"]; ok {
		data.CalleeMethod = cast.ToString(v)
		delete(enc.Fields, "calleeMethod")
	}

	if v, ok := enc.Fields["req"]; ok {
		data.Req = cast.ToString(v)
		delete(enc.Fields, "req")
	}
	if v, ok := enc.Fields["rsp"]; ok {
		data.Rsp = cast.ToString(v)
		delete(enc.Fields, "rsp")
	}

	// 序列化fields
	if len(enc.Fields) > 0 {
		var fieldsBuff bytes.Buffer
		for k, v := range enc.Fields {
			fieldsBuff.WriteByte(',')
			fieldsBuff.WriteString(k)
			fieldsBuff.WriteByte('=')
			fieldsBuff.WriteString(cast.ToString(v))
		}
		data.Fields = string(fieldsBuff.Bytes()[1:])
	}

	return data
}
