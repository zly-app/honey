package main

import (
	"strings"

	"github.com/zly-app/honey/output"
)

// 生成输出设备
func (h *Honey) MakeOutput() {
	if h.conf.Outputs == "" {
		return
	}

	names := strings.Split(h.conf.Outputs, ",")
	h.outputs = make([]output.IOutput, len(names))
	for i := range names {
		out := output.MakeOutput(h.c, names[i])
		h.outputs[i] = out
	}
}

// 启动输出设备
func (h *Honey) StartOutput() {
	for _, out := range h.outputs {
		out.Start()
	}
}

// 关闭输出设备
func (h *Honey) CloseOutput() {
	for _, out := range h.outputs {
		out.Close()
	}
}
