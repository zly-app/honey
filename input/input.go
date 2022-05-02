package input

import (
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zly-app/honey/component"
)

// 输入设备建造者
type InputCreator func(c component.ILogCollector, iConfig component.IInputConfig) IInput

// 输入设备
type IInput interface {
	Start() error
	Close() error
}

var inputCreators = make(map[string]InputCreator)

// 注册输入设备建造者
func RegistryInputCreator(name string, rc InputCreator) {
	if _, ok := inputCreators[name]; ok {
		logger.Log.Fatal("重复注册Input建造者", zap.String("name", name))
	}
	inputCreators[name] = rc
}

// 生成输入设备
func MakeInput(c component.ILogCollector, iConfig component.IInputConfig, name string) IInput {
	ic, ok := inputCreators[name]
	if !ok {
		logger.Log.Fatal("试图构建未注册建造者的Input", zap.String("name", name))
	}
	return ic(c, iConfig)
}
