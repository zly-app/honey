package component

import (
	"github.com/zly-app/zapp/core"

	"github.com/zly-app/honey/log_data"
)

// 日志收集器
type LogCollector interface {
	// 收集
	Collect(env, service, instance string, log []*log_data.LogData)
}

type IComponent interface {
	core.IComponent
	LogCollector
}

type Component struct {
	core.IComponent
	LogCollector
}

func (c *Component) Close() {
	c.IComponent.Close()
}

// 获取Component
func GetComponent(c core.IComponent) IComponent {
	return c.(IComponent)
}
