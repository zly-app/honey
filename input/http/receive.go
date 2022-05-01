package http

import (
	"bytes"
	"fmt"

	"github.com/zly-app/service/api"

	"github.com/zly-app/honey/log_data"
)

const (
	// 如何从header中获取env
	HeaderNameEnv = "env"
	// 如何从header中获取service
	HeaderNameService = "service"
	// 如何从header中获取instance
	HeaderNameInstance = "instance"
)

// 接收
func (h *HttpInput) Receive(ctx *api.Context) error {
	env := ctx.GetHeader(HeaderNameEnv)
	service := ctx.GetHeader(HeaderNameService)
	instance := ctx.GetHeader(HeaderNameInstance)

	if env == "" {
		return fmt.Errorf("env为空")
	}
	if service == "" {
		return fmt.Errorf("service为空")
	}
	if instance == "" {
		instance = ctx.RemoteAddr()
	}

	body := bytes.NewBuffer(nil)
	err := h.compress.UnCompress(ctx.Request().Body, body)
	if err != nil {
		return fmt.Errorf("解压缩失败: %v", err)
	}

	logs := []*log_data.LogData{}
	err = h.serializer.Unmarshal(body, &logs)
	if err != nil {
		return fmt.Errorf("解码失败: %v", err)
	}

	if len(logs) > 0 {
		h.c.Collect(env, service, instance, logs)
	}

	return nil
}
