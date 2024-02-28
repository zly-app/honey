package log_data

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
	"strconv"

	"github.com/bytedance/sonic"
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
		data.TraceID = ToString(v)
		delete(enc.Fields, "traceID")
	}

	// 提取spanID
	if v, ok := enc.Fields["spanID"]; ok {
		data.SpanID = ToString(v)
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
		data.CodeType = ToString(v)
		delete(enc.Fields, "codeType")
	}
	if v, ok := enc.Fields["callerService"]; ok {
		data.CallerService = ToString(v)
		delete(enc.Fields, "callerService")
	}
	if v, ok := enc.Fields["callerMethod"]; ok {
		data.CallerMethod = ToString(v)
		delete(enc.Fields, "callerMethod")
	}
	if v, ok := enc.Fields["calleeService"]; ok {
		data.CalleeService = ToString(v)
		delete(enc.Fields, "calleeService")
	}
	if v, ok := enc.Fields["calleeMethod"]; ok {
		data.CalleeMethod = ToString(v)
		delete(enc.Fields, "calleeMethod")
	}

	if v, ok := enc.Fields["req"]; ok {
		data.Req = ToString(v)
		delete(enc.Fields, "req")
	}
	if v, ok := enc.Fields["rsp"]; ok {
		data.Rsp = ToString(v)
		delete(enc.Fields, "rsp")
	}

	// 序列化fields
	if len(enc.Fields) > 0 {
		var fieldsBuff bytes.Buffer
		for k, v := range enc.Fields {
			fieldsBuff.WriteByte(',')
			fieldsBuff.WriteString(k)
			fieldsBuff.WriteByte('=')
			fieldsBuff.WriteString(ToString(v))
		}
		data.Fields = string(fieldsBuff.Bytes()[1:])
	}

	return data
}

func indirectToStringerOrError(a interface{}) interface{} {
	if a == nil {
		return nil
	}

	var errorType = reflect.TypeOf((*error)(nil)).Elem()
	var fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

	v := reflect.ValueOf(a)
	for !v.Type().Implements(fmtStringerType) && !v.Type().Implements(errorType) && v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

func ToString(i interface{}) string {
	i = indirectToStringerOrError(i)

	switch s := i.(type) {
	case string:
		return s
	case bool:
		return strconv.FormatBool(s)
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32)
	case int:
		return strconv.Itoa(s)
	case int64:
		return strconv.FormatInt(s, 10)
	case int32:
		return strconv.Itoa(int(s))
	case int16:
		return strconv.FormatInt(int64(s), 10)
	case int8:
		return strconv.FormatInt(int64(s), 10)
	case uint:
		return strconv.FormatInt(int64(s), 10)
	case uint64:
		return strconv.FormatInt(int64(s), 10)
	case uint32:
		return strconv.FormatInt(int64(s), 10)
	case uint16:
		return strconv.FormatInt(int64(s), 10)
	case uint8:
		return strconv.FormatInt(int64(s), 10)
	case []byte:
		return string(s)
	case template.HTML:
		return string(s)
	case template.URL:
		return string(s)
	case template.JS:
		return string(s)
	case template.CSS:
		return string(s)
	case template.HTMLAttr:
		return string(s)
	case nil:
		return "<nil>"
	case fmt.Stringer:
		return s.String()
	case error:
		return s.Error()
	default:
		v, _ := sonic.MarshalString(i)
		return v
	}
}
