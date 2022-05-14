package component

import (
	"github.com/zly-app/honey/log_data"
)

// 日志收集器
type ILogCollector interface {
	// 收集
	Collect(env, app, instance string, log []*log_data.LogData)
}

// 输入器配置
type IInputConfig interface {
	/*解析输入设备配置数据到结构中
	  配置项在配置文件中为 [input.{输入设备名}]
	  name 输入设备名
	  outPtr 接收配置的变量
	  ignoreNotSet 如果无配置数据, 则忽略, 默认为false
	*/
	ParseInputConf(name string, conf interface{}, ignoreNotSet ...bool) error
}

// 输出器配置
type IOutputConfig interface {
	/*解析输出设备配置数据到结构中
	  配置项在配置文件中为 [output.{输出设备名}]
	  name 输出设备名
	  outPtr 接收配置的变量
	  ignoreNotSet 如果无配置数据, 则忽略, 默认为false
	*/
	ParseOutputConf(name string, conf interface{}, ignoreNotSet ...bool) error
}
