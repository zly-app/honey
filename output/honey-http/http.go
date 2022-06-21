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
	"github.com/zly-app/zapp/logger"
	"github.com/zly-app/zapp/pkg/utils"
	"github.com/zlyuancn/zretry"
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
	// 如何从header中设置app
	HeaderNameApp = "app"
	// 如何从header中设置instance
	HeaderNameInstance = "instance"
)

type HttpOutput struct {
	conf       *Config
	compress   compress.ICompress
	serializer serializer.ISerializer
	client     *http.Client
}

func (h *HttpOutput) Start() error { return nil }

func (h *HttpOutput) Close() error { return nil }

func (h *HttpOutput) Out(env, app, instance string, data []*log_data.LogData) {
	if h.conf.Disable {
		return
	}

	_ = zretry.DoRetry(h.conf.RetryCount+1, time.Duration(h.conf.RetryIntervalMs)*time.Millisecond,
		func() error {
			return h.out(env, app, instance, data)
		},
		func(nowAttemptCount, remainCount int, err error) {
			_, _ = os.Stdout.WriteString(fmt.Sprintf("输出失败, 剩余重试 %d 次, err: %v\n", remainCount, err.Error()))
		})
}

func (h *HttpOutput) out(env, app, instance string, data []*log_data.LogData) error {

	// 序列化
	buff := bytes.NewBuffer(nil)
	err := h.serializer.Marshal(data, buff)
	if err != nil {
		return fmt.Errorf("序列化日志数据失败: %v", err)
	}

	// 编码
	body := bytes.NewBuffer(nil)
	err = h.compress.Compress(buff, body)
	if err != nil {
		return fmt.Errorf("编码失败: %v", err)
	}

	// 构建请求
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.conf.ReqTimeout)*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", h.conf.PushAddress, body)
	if err != nil {
		return fmt.Errorf("构建请求体失败: %v", err)
	}
	req.Header.Add("Content-Encoding", h.conf.Compress)
	req.Header.Add("Content-Type", h.conf.Serializer)
	if h.conf.AuthToken != "" {
		req.Header.Add("token", h.conf.AuthToken)
	}
	req.Header.Add(HeaderNameEnv, env)
	req.Header.Add(HeaderNameApp, app)
	req.Header.Add(HeaderNameInstance, instance)

	// 请求
	rsp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("上报失败, 请求失败: err:%v", err)
	}
	defer rsp.Body.Close()
	rspBody, _ := io.ReadAll(rsp.Body)

	// 检查状态码
	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("上报失败, 错误的状态码: code:%v, body:%v",
			rsp.StatusCode, string(rspBody))
	}

	// 解析body
	result := api.Response{}
	err = json.Unmarshal(rspBody, &result)
	if err != nil {
		return fmt.Errorf("解析rsp失败: body:%v, err:%v",
			string(rspBody), err)
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("上报失败, 错误的响应: errCode:%v, errMsg",
			result.ErrCode, result.ErrMsg)
	}
	return nil
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
		client:     &http.Client{},
	}

	if conf.ProxyAddress != "" {
		p, err := utils.NewHttpProxy(conf.ProxyAddress)
		if err != nil {
			logger.Log.Fatal("创建loki-http代理失败", zap.Error(err))
		}
		transport := &http.Transport{}
		p.SetProxy(transport)
		h.client.Transport = transport
	}

	return h
}

const HoneyHttpOutputName = "honey-http"

func init() {
	output.RegistryOutputCreator(HoneyHttpOutputName, func(iConfig component.IOutputConfig) output.IOutput {
		return NewHttpOutput(iConfig)
	})
}
