package std_output

import (
	"fmt"
	"os"

	"github.com/zly-app/honey/component"
	"github.com/zly-app/honey/log_data"
	"github.com/zly-app/honey/output"
)

type StdOutput struct{}

func (s *StdOutput) Start() {}
func (s *StdOutput) Close() {}

func (s *StdOutput) Out(env string, data []*log_data.CollectData) {
	for _, v := range data {
		_, _ = os.Stdout.WriteString(fmt.Sprint(env, v.Service, v.Instance, *v.LogData, "\n"))
	}
}

// std输出设备名
const StdOutputName = "std"

func init() {
	output.RegistryOutputCreator(StdOutputName, func(c component.IComponent) output.IOutput {
		return &StdOutput{}
	})
}
