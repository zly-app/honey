package main

import (
	"fmt"
	"strings"

	"github.com/zly-app/zapp/service"
	"go.uber.org/zap"

	"github.com/zly-app/honey/output"
)

// 生成输出设备
func (h *Honey) MakeOutput() {
	if h.conf.Outputs == "" {
		return
	}

	names := strings.Split(h.conf.Outputs, ",")
	h.outputs = make(map[string]output.IOutput, len(names))
	for i := range names {
		out := output.MakeOutput(h.c, names[i])
		h.outputs[names[i]] = out
	}
}

// 启动输出设备
func (h *Honey) StartOutput() {
	for name, out := range h.outputs {
		err := service.WaitRun(h.app, &service.WaitRunOption{
			ServiceType:        DefaultServiceType,
			ExitOnErrOfObserve: true,
			RunServiceFn: func() error {
				err := out.Start()
				if err == nil {
					return nil
				}
				return fmt.Errorf("启动Output失败, output: %s, err: %v", name, err)
			},
		})
		if err != nil {
			h.app.Fatal("启动Output失败", zap.String("output", name), zap.Error(err))
		}
	}
}

// 关闭输出设备
func (h *Honey) CloseOutput() {
	for name, out := range h.outputs {
		err := out.Close()
		if err != nil {
			h.app.Error("关闭Output失败", zap.String("output", name), zap.Error(err))
		}
	}
}
