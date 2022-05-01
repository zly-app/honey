package honey

import (
	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/pkg/zlog"

	_ "github.com/zly-app/honey/zapp_plugin/reporter/stdout"
)

// 默认插件类型
const DefaultPluginType core.PluginType = "honey"

// 当前插件类型
var nowPluginType = DefaultPluginType

// 设置插件类型, 这个函数应该在 zapp.NewApp 之前调用
func SetPluginType(t core.PluginType) {
	nowPluginType = t
}

// 启用插件
func WithPlugin() zapp.Option {
	h := NewHoneyPlugin()
	logConf := zlog.NewHookConfig().
		AddStartHookCallbacks(h.Init).           // 通过日志启动hook的能力提供初始化
		AddInterceptorFunc(h.LogInterceptorFunc) // 添加拦截函数
	return zapp.WithLoggerOptions(zlog.WithHookByConfig(logConf))
}

func init() {
	zapp.AddHandler(zapp.AfterInitializeHandler, func(app core.IApp, handlerType zapp.HandlerType) {
		h := NewHoneyPlugin()
		h.Start()
	})
	zapp.AddHandler(zapp.BeforeExitHandler, func(app core.IApp, handlerType zapp.HandlerType) {
		h := NewHoneyPlugin()
		h.Close()
	})
}