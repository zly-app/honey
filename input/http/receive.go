package http

import (
	"bytes"
	"fmt"

	"github.com/zly-app/service/api"

	"github.com/zly-app/honey/log_data"
	"github.com/zly-app/honey/pkg/compress"
	"github.com/zly-app/honey/pkg/serializer"
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

	// 获取压缩程序
	compressName := ctx.GetHeader("Content-Encoding")
	if compressName == "" {
		compressName = compress.RawCompressName
	}
	cp, ok := compress.TryGetCompress(compressName)
	if !ok {
		return fmt.Errorf("不支持的压缩类型: %v", compressName)
	}

	// 获取序列化器
	ct := ctx.GetHeader("Content-Type")
	se, ok := serializer.TryGetSerializer(ct)
	if !ok {
		return fmt.Errorf("不支持的数据类型: %v", ct)
	}

	body := bytes.NewBuffer(nil)
	err := cp.UnCompress(ctx.Request().Body, body)
	if err != nil {
		return fmt.Errorf("解压缩失败: %v", err)
	}

	logs := []*log_data.LogData{}
	err = se.Unmarshal(body, &logs)
	if err != nil {
		return fmt.Errorf("解码失败: %v", err)
	}

	if len(logs) > 0 {
		h.c.Collect(env, service, instance, logs)
	}

	return nil
}
