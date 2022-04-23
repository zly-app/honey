package config

import (
	"github.com/zly-app/zapp"

	"github.com/zly-app/honey/pkg/utils"
)

var Conf *Config

const (
	// 默认环境名
	DefaultThisLogEnv = "dev"
	// 停止原有的日志输出
	DefaultThisLogStopLogOutput = true

	// 默认批次大小
	DefaultLogBatchSize = 10000
	// 默认自动旋转时间(秒)
	DefaultAutoRotateTime = 5
	// 默认最大旋转线程数
	DefaultMaxRotateThreadNum = 10

	// 默认输出设备列表
	DefaultOutputs = "std"
)

type Config struct {
	ThisLog struct {
		Disable       bool   // 关闭honey服务本身的日志收集
		Env           string // 环境名
		Service       string // 服务名
		Instance      string // 实例名
		StopLogOutput bool   // 停止原有的日志输出, honey启动后不会输出日志到屏幕或原有的
	} // honey服务本身的log处理文件

	LogBatchSize       int // 日志批次大小, 累计达到这个大小立即写入一次日志, 不用等待时间
	AutoRotateTime     int // 自动旋转时间(秒), 如果没有达到累计写入批次大小, 在指定时间后也会立即写入
	MaxRotateThreadNum int // 最大旋转线程数, 表示同时允许多少批次发送到输出设备

	Inputs  string // 输入设备列表, 多个输入设备用半角逗号`,`分隔, 目前支持的输入设备: http
	Outputs string // 输出设备列表, 多个输出设备用半角逗号`,`分隔, 目前支持的输出设备: std
}

func NewConfig() *Config {
	conf := &Config{}
	conf.ThisLog.StopLogOutput = DefaultThisLogStopLogOutput
	Conf = conf
	return conf
}

func (conf *Config) Check() error {
	if conf.ThisLog.Env == "" {
		conf.ThisLog.Env = DefaultThisLogEnv
	}
	if conf.ThisLog.Service == "" {
		conf.ThisLog.Service = zapp.App().Name()
	}
	if conf.ThisLog.Instance == "" {
		conf.ThisLog.Instance = utils.GetInstance()
	}

	if conf.LogBatchSize < 1 {
		conf.LogBatchSize = DefaultLogBatchSize
	}
	if conf.AutoRotateTime < 1 {
		conf.AutoRotateTime = DefaultAutoRotateTime
	}
	if conf.MaxRotateThreadNum < 1 {
		conf.MaxRotateThreadNum = DefaultMaxRotateThreadNum
	}

	if conf.Outputs == "" {
		conf.Outputs = DefaultOutputs
	}

	return nil
}
