package std_output

import (
	"os"
	"strings"
	"time"

	"github.com/zly-app/honey/component"
	"github.com/zly-app/honey/log_data"
	"github.com/zly-app/honey/output"
)

var StdFormat = "[{dev}.{service}][{instance}][{time}] {level} {msg} {fields} {line} {trace_id}"
var TimeFormat = "2006-01-02 15:04:05.999999"

type StdOutput struct{}

func (s *StdOutput) Start() error { return nil }
func (s *StdOutput) Close() error { return nil }

func (s *StdOutput) Out(env, service, instance string, data []*log_data.LogData) {
	for _, v := range data {
		text := StdFormat
		text = strings.ReplaceAll(text, "{dev}", env)
		text = strings.ReplaceAll(text, "{service}", service)
		text = strings.ReplaceAll(text, "{instance}", instance)
		text = strings.ReplaceAll(text, "{time}", time.Unix(0, v.T*1000).Format(TimeFormat))
		text = strings.ReplaceAll(text, "{level}", v.Level)
		text = strings.ReplaceAll(text, "{msg}", v.Msg)
		text = strings.ReplaceAll(text, "{fields}", v.Fields)
		text = strings.ReplaceAll(text, "{line}", v.Line)
		text = strings.ReplaceAll(text, "{trace_id}", v.TraceID)

		_, _ = os.Stdout.WriteString(text)
		_, _ = os.Stdout.Write([]byte{'\n'})
	}
}

// std输出设备名
const StdOutputName = "std"

func init() {
	output.RegistryOutputCreator(StdOutputName, func(c component.IComponent) output.IOutput {
		return &StdOutput{}
	})
}
