package main

import (
	"github.com/zly-app/honey/log_data"
)

// 收集日志
func (h *Honey) Collect(env, service, instance string, log []*log_data.LogData) {
	rotate := h.rotateGroup.GetRotate(env, service, instance)
	// 写入旋转器
	for _, v := range log {
		rotate.Add(v)
	}
}
