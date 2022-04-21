package main

import (
	"strings"

	"github.com/zly-app/honey/input"
)

// 生成输入设备
func (h *Honey) MakeInput() {
	if h.conf.Inputs == "" {
		return
	}

	names := strings.Split(h.conf.Inputs, ",")
	h.inputs = make([]input.IInput, len(names))
	for i := range names {
		in := input.MakeInput(h.c, names[i])
		h.inputs[i] = in
	}
}

// 启动输入设备
func (h *Honey) StartInput() {
	for _, in := range h.inputs {
		in.Start()
	}
}

// 关闭输入设备
func (h *Honey) CloseInput() {
	for _, in := range h.inputs {
		in.Close()
	}
}
