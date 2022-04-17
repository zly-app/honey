package component

import (
	"github.com/zly-app/zapp/core"
)

type Component struct {
	core.IComponent
}

func (c *Component) Close() {
	c.IComponent.Close()
}

func GetComponent(c core.IComponent) *Component {
	return c.(*Component)
}
