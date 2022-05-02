package main

import (
	"fmt"
	"strings"

	"github.com/zly-app/zapp/service"
	"go.uber.org/zap"

	"github.com/zly-app/honey/input"
)

// 生成输入设备
func (h *Honey) MakeInput() {
	if h.conf.Inputs == "" {
		return
	}

	names := strings.Split(h.conf.Inputs, ",")
	h.inputs = make(map[string]input.IInput, len(names))
	for i := range names {
		in := input.MakeInput(h, h, names[i])
		h.inputs[names[i]] = in
	}
}

// 解析输入设备配置数据到结构中
func (h *Honey) ParseInputConf(name string, conf interface{}, ignoreNotSet ...bool) error {
	key := fmt.Sprintf("input.%s", name)
	return h.app.GetConfig().Parse(key, conf, ignoreNotSet...)
}

// 启动输入设备
func (h *Honey) StartInput() {
	for name, in := range h.inputs {
		err := service.WaitRun(h.app, &service.WaitRunOption{
			ServiceType:        DefaultServiceType,
			ExitOnErrOfObserve: true,
			RunServiceFn: func() error {
				err := in.Start()
				if err == nil {
					return nil
				}
				return fmt.Errorf("启动Input失败, input: %s, err: %v", name, err)
			},
		})
		if err != nil {
			h.app.Fatal("启动Input失败", zap.String("input", name), zap.Error(err))
		}
	}
}

// 关闭输入设备
func (h *Honey) CloseInput() {
	for name, in := range h.inputs {
		err := in.Close()
		if err != nil {
			h.app.Error("关闭Input失败", zap.String("input", name), zap.Error(err))
		}
	}
}
