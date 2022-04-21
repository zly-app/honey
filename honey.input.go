package main

import (
	"strings"

	"github.com/zly-app/honey/input"
)

// 初始化输入设备
func (h *Honey) InitInput() {
	if h.conf.Inputs == "" {
		return
	}

	names := strings.Split(h.conf.Inputs, ",")
	for _, name := range names {
		input.MakeInput(h.c, name)
	}
}
