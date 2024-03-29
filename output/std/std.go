package std

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/cast"

	"github.com/zly-app/honey/component"
	"github.com/zly-app/honey/log_data"
	"github.com/zly-app/honey/output"
)

var StdFormat = "[{env}.{app}][{instance}][{time}] {level} {msg} line={line},{func} traceID={traceID} spanID={spanID} duration={duration} code={code} codeType={codeType} callerService={callerService} callerMethod={callerMethod} calleeService={calleeService} calleeMethod={calleeMethod} req={req} rsp={rsp} fields={fields}"
var TimeFormat = "2006-01-02 15:04:05.999999"

type StdOutput struct{}

func (s *StdOutput) Start() error { return nil }
func (s *StdOutput) Close() error { return nil }

func (s *StdOutput) Out(env, app, instance string, data []*log_data.LogData) {
	for _, v := range data {
		text := StdFormat
		text = strings.ReplaceAll(text, "{env}", env)
		text = strings.ReplaceAll(text, "{app}", app)
		text = strings.ReplaceAll(text, "{instance}", instance)
		text = strings.ReplaceAll(text, "{time}", time.Unix(0, v.T).Format(TimeFormat))
		text = strings.ReplaceAll(text, "{level}", v.Level)
		text = strings.ReplaceAll(text, "{msg}", v.Msg)
		text = strings.ReplaceAll(text, "{fields}", v.Fields)
		text = strings.ReplaceAll(text, "{line}", v.Line)
		text = strings.ReplaceAll(text, "{func}", v.Func)
		text = strings.ReplaceAll(text, "{traceID}", v.TraceID)
		text = strings.ReplaceAll(text, "{spanID}", v.SpanID)
		text = strings.ReplaceAll(text, "{req}", v.Req)
		text = strings.ReplaceAll(text, "{rsp}", v.Rsp)

		text = strings.ReplaceAll(text, "{duration}", cast.ToString(v.Duration))
		text = strings.ReplaceAll(text, "{code}", cast.ToString(v.Code))
		text = strings.ReplaceAll(text, "{codeType}", v.CodeType)
		text = strings.ReplaceAll(text, "{callerService}", v.CallerService)
		text = strings.ReplaceAll(text, "{callerMethod}", v.CallerMethod)
		text = strings.ReplaceAll(text, "{calleeService}", v.CalleeService)
		text = strings.ReplaceAll(text, "{calleeMethod}", v.CalleeMethod)

		_, _ = os.Stdout.WriteString(text)
		_, _ = os.Stdout.Write([]byte{'\n'})
	}
}

// std输出设备名
const StdOutputName = "std"

func init() {
	output.RegistryOutputCreator(StdOutputName, func(iConfig component.IOutputConfig) output.IOutput {
		return &StdOutput{}
	})
}
