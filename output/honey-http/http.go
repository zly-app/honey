package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/zly-app/service/api"
	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zly-app/honey/component"
	"github.com/zly-app/honey/log_data"
	"github.com/zly-app/honey/output"
	"github.com/zly-app/honey/pkg/compress"
	"github.com/zly-app/honey/pkg/serializer"
)

const (
	// 如何从header中设置env
	HeaderNameEnv = "env"
	// 如何从header中设置service
	HeaderNameService = "service"
	// 如何从header中设置instance
	HeaderNameInstance = "instance"
)

type HttpOutput struct {
	conf       *Config
	compress   compress.ICompress
	serializer serializer.ISerializer
}

func (h *HttpOutput) Start() error { return nil }

func (h *HttpOutput) Close() error { return nil }

func (h *HttpOutput) Out(env, service, instance string, data []*log_data.LogData) {
	if h.conf.Disable {
		return
	}

	// 序列化
	buff := bytes.NewBuffer(nil)
	err := h.serializer.Marshal(data, buff)
	if err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("序列化日志数据失败: %v\n", err))
		return
	}

	// 编码
	body := bytes.NewBuffer(nil)
	err = h.compress.Compress(buff, body)
	if err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("编码失败: %v\n", err))
		return
	}

	// 构建请求
	ctx, cancel := context.WithTimeout(zapp.App().BaseContext(), time.Duration(h.conf.ReqTimeout)*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", h.conf.PushAddress, body)
	if err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("构建请求体失败: %v\n", err))
		return
	}
	req.Header.Add("Content-Encoding", h.conf.Compress)
	req.Header.Add("Content-Type", h.conf.Serializer)
	if h.conf.AuthToken != "" {
		req.Header.Add("token", h.conf.AuthToken)
	}
	req.Header.Add(HeaderNameEnv, env)
	req.Header.Add(HeaderNameService, service)
	req.Header.Add(HeaderNameInstance, instance)

	// 请求
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("上报失败, 请求失败: err:%v\n", err))
		return
	}
	defer rsp.Body.Close()
	rspBody, _ := io.ReadAll(rsp.Body)

	// 检查状态码
	if rsp.StatusCode != http.StatusOK {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("上报失败, 错误的状态码: code:%v, body:%v\n",
			rsp.StatusCode, string(rspBody)))
		return
	}

	// 解析body
	result := api.Response{}
	err = json.Unmarshal(rspBody, &result)
	if err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("解析rsp失败: body:%v, err:%v\n",
			string(rspBody), err))
		return
	}

	if result.ErrCode != 0 {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("上报失败, 错误的响应: errCode:%v, errMsg\n",
			result.ErrCode, result.ErrMsg))
		return
	}
}

func NewHttpOutput(iConfig component.IOutputConfig) *HttpOutput {
	conf := newConfig()
	err := iConfig.ParseOutputConf(HoneyHttpOutputName, conf, true)
	if err == nil {
		err = conf.Check()
	}
	if err != nil {
		logger.Log.Fatal("获取honey-http输出设备配置失败", zap.Error(err))
	}

	h := &HttpOutput{
		conf:       conf,
		compress:   compress.GetCompress(conf.Compress),
		serializer: serializer.GetSerializer(conf.Serializer),
	}
	return h
}

const HoneyHttpOutputName = "honey-http"

func init() {
	output.RegistryOutputCreator(HoneyHttpOutputName, func(iConfig component.IOutputConfig) output.IOutput {
		return NewHttpOutput(iConfig)
	})
}
