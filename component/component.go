package component

import (
	"fmt"

	"github.com/zly-app/zapp/core"

	"github.com/zly-app/honey/log_data"
)

// 日志收集器
type ILogCollector interface {
	// 收集
	Collect(env, service, instance string, log []*log_data.LogData)
}

type IComponent interface {
	core.IComponent
	ILogCollector
	/*解析输入设备配置数据到结构中
	  配置项在配置文件中为 [input.{输入设备名}]
	  name 输入设备名
	  outPtr 接收配置的变量
	  ignoreNotSet 如果无配置数据, 则忽略, 默认为false
	*/
	ParseInputConf(name string, conf interface{}, ignoreNotSet ...bool) error
	/*解析输出设备配置数据到结构中
	  配置项在配置文件中为 [output.{输出设备名}]
	  name 输出设备名
	  outPtr 接收配置的变量
	  ignoreNotSet 如果无配置数据, 则忽略, 默认为false
	*/
	ParseOutputConf(name string, conf interface{}, ignoreNotSet ...bool) error
}

var _ IComponent = (*Component)(nil)

type Component struct {
	core.IComponent
	ILogCollector
}

// 解析输入设备配置数据到结构中
func (c *Component) ParseInputConf(name string, conf interface{}, ignoreNotSet ...bool) error {
	key := fmt.Sprintf("input.%s", name)
	return c.App().GetConfig().Parse(key, conf, ignoreNotSet...)
}

// 解析输出设备配置数据到结构中
func (c *Component) ParseOutputConf(name string, conf interface{}, ignoreNotSet ...bool) error {
	key := fmt.Sprintf("output.%s", name)
	return c.App().GetConfig().Parse(key, conf, ignoreNotSet...)
}

func (c *Component) Close() {
	c.IComponent.Close()
}

// 获取Component
func GetComponent(c core.IComponent) IComponent {
	return c.(IComponent)
}
