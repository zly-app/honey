package std

import (
	"os"
	"strings"
	"time"

	"github.com/zly-app/honey/component"
	"github.com/zly-app/honey/log_data"
	"github.com/zly-app/honey/output"
)

var StdFormat = "[{env}.{app}][{instance}][{time}] {level} {msg} line={line} fields={fields} traceID={traceID} spanID={spanID}"
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
		text = strings.ReplaceAll(text, "{traceID}", v.TraceID)
		text = strings.ReplaceAll(text, "{spanID}", v.SpanID)

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
