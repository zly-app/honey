package http

import (
	"github.com/zly-app/honey/component"
	"github.com/zly-app/honey/input"
)

// http输入设备名
const HttpInputName = "http"

func init() {
	input.RegistryInputCreator(HttpInputName, func(c component.ILogCollector, ic component.IInputConfig) input.IInput {
		return NewHttpInput(c, ic)
	})
}
