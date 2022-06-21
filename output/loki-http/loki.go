package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/zly-app/zapp/logger"
	"github.com/zly-app/zapp/pkg/utils"
	"go.uber.org/zap"

	"github.com/zlyuancn/zretry"

	"github.com/zly-app/honey/component"
	"github.com/zly-app/honey/log_data"
	"github.com/zly-app/honey/output"
	"github.com/zly-app/honey/pkg/compress"
)

type HttpOutput struct {
	conf   *Config
	client *http.Client
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
	lokiData := MakeLokiBody(env, app, instance, data)

	// 序列化
	body := bytes.NewBuffer(nil)
	err := jsoniter.ConfigCompatibleWithStandardLibrary.NewEncoder(body).Encode(lokiData)
	if err != nil {
		return fmt.Errorf("序列化日志数据失败: %v", err)
	}

	// 编码
	if h.conf.EnableCompress {
		tmp := bytes.NewBuffer(nil)
		err = compress.GetCompress(compress.GzipCompressName).Compress(body, tmp)
		if err != nil {
			return fmt.Errorf("编码失败: %v", err)
		}
		body = tmp
	}

	// 构建请求
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.conf.ReqTimeout)*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", h.conf.PushAddress, body)
	if err != nil {
		return fmt.Errorf("构建请求体失败: %v", err)
	}
	if h.conf.EnableCompress {
		req.Header.Add("Content-Encoding", compress.GzipCompressName)
	}
	req.Header.Add("Content-Type", "application/json")

	// 请求
	rsp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("上报失败, 请求失败: err:%v", err)
	}
	defer rsp.Body.Close()
	rspBody, _ := io.ReadAll(rsp.Body)

	// 检查状态码
	if rsp.StatusCode != http.StatusOK && rsp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("上报失败, 错误的状态码: code:%v, body:%v",
			rsp.StatusCode, string(rspBody))
	}
	return nil
}

func NewHttpOutput(iConfig component.IOutputConfig) *HttpOutput {
	conf := newConfig()
	err := iConfig.ParseOutputConf(LokiHttpOutputName, conf, true)
	if err == nil {
		err = conf.Check()
	}
	if err != nil {
		logger.Log.Fatal("获取loki-http输出设备配置失败", zap.Error(err))
	}

	h := &HttpOutput{
		conf:   conf,
		client: &http.Client{},
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

const LokiHttpOutputName = "loki-http"

func init() {
	output.RegistryOutputCreator(LokiHttpOutputName, func(iConfig component.IOutputConfig) output.IOutput {
		return NewHttpOutput(iConfig)
	})
}
