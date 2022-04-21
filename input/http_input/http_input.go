package http_input

import (
	"github.com/zly-app/honey/component"
)

type HttpInput struct {
	c component.IComponent
}

func NewHttpInput(c component.IComponent) *HttpInput {
	return &HttpInput{
		c: c,
	}
}
