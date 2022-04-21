package main

import (
	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/pkg/zlog"
	"go.uber.org/zap/zapcore"

	"github.com/zly-app/honey/log_data"
)

// 日志拦截函数
func (h *Honey) logInterceptorFunc(ent *zapcore.Entry, fields []zapcore.Field) (cancel bool) {
	conf := h.conf.ThisLog
	if conf.Disable {
		return false
	}

	log := log_data.MakeLogData(ent, fields)
	h.Collect(conf.Env, conf.Service, conf.Instance, []*log_data.LogData{log})
	return conf.StopLogOutput && h.isStart() // 设置了拦截并且在服务启动后才允许拦截
}

// 提供日志hook的zapp选项
func (h *Honey) LogHook() zapp.Option {
	logConf := zlog.NewHookConfig().
		AddStartHookCallbacks(h.Init).           // 通过日志启动hook的能力提供初始化
		AddInterceptorFunc(h.logInterceptorFunc) // 添加拦截函数
	return zapp.WithLoggerOptions(zlog.WithHookByConfig(logConf))
}

// 收集日志
func (h *Honey) Collect(env, service, instance string, log []*log_data.LogData) {
	rotate := h.rotateGroup.GetRotate(env)
	// 写入旋转器
	for _, v := range log {
		data := log_data.MakeCollectData(service, instance, v)
		rotate.Add(data)
	}
}
