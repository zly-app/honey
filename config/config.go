package config

import (
	"github.com/zly-app/zapp"

	"github.com/zly-app/honey/pkg/instance"
)

var Conf *Config

const (
	DefaultBatchSize = 10000 // 默认批次大小
	// 默认自动旋转时间(秒)
	DefaultAutoRotateTime = 5

	// 默认环境名
	DefaultThisLogEnv = "dev"
	// 停止原有的日志输出
	DefaultThisLogStopLogOutput = true
)

type Config struct {
	ThisLog struct {
		Disable       bool   // 关闭honey服务本身的日志收集
		Env           string // 环境名
		Service       string // 服务名
		Instance      string // 实例名
		StopLogOutput bool   // 停止原有的日志输出, honey启动后不会输出日志到屏幕或原有的
	} // honey服务本身的log处理文件

	BatchSize      int // 批次大小, 累计达到这个大小立即写入一次日志, 不用等待时间
	AutoRotateTime int // 自动旋转时间(秒), 如果没有达到累计写入批次大小, 在指定时间后也会立即写入

	HttpReceiver bool // 启用http接收器
}

func NewConfig() *Config {
	conf := &Config{}
	conf.ThisLog.StopLogOutput = DefaultThisLogStopLogOutput
	Conf = conf
	return conf
}

func (conf *Config) Check() error {
	if conf.BatchSize < 1 {
		conf.BatchSize = DefaultBatchSize
	}
	if conf.AutoRotateTime < 1 {
		conf.AutoRotateTime = DefaultAutoRotateTime
	}

	if conf.ThisLog.Env == "" {
		conf.ThisLog.Env = DefaultThisLogEnv
	}
	if conf.ThisLog.Service == "" {
		conf.ThisLog.Service = zapp.App().Name()
	}
	if conf.ThisLog.Instance == "" {
		conf.ThisLog.Instance = utils.GetInstance()
	}
	return nil
}
