package output

import (
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zly-app/honey/component"
	"github.com/zly-app/honey/log_data"
)

// 输出设备建造者
type OutputCreator func(iConfig component.IOutputConfig) IOutput

// 输出设备
type IOutput interface {
	Start() error
	Close() error
	// 输出
	Out(env, app, instance string, data []*log_data.LogData)
}

var outputCreators = make(map[string]OutputCreator)

// 注册输出设备建造者
func RegistryOutputCreator(name string, oc OutputCreator) {
	if _, ok := outputCreators[name]; ok {
		logger.Log.Fatal("重复注册Output建造者", zap.String("name", name))
	}
	outputCreators[name] = oc
}

// 生成输出设备
func MakeOutput(iConfig component.IOutputConfig, name string) IOutput {
	oc, ok := outputCreators[name]
	if !ok {
		logger.Log.Fatal("试图构建未注册建造者的Output", zap.String("name", name))
	}
	return oc(iConfig)
}
