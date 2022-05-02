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
	"go.uber.org/zap"

	"github.com/zly-app/honey/component"
	"github.com/zly-app/honey/log_data"
	"github.com/zly-app/honey/output"
	"github.com/zly-app/honey/pkg/compress"
)

type HttpOutput struct {
	conf *Config
}

func (h *HttpOutput) Start() error { return nil }

func (h *HttpOutput) Close() error { return nil }

func (h *HttpOutput) Out(env, service, instance string, data []*log_data.LogData) {
	if h.conf.Disable {
		return
	}

	lokiData := MakeLokiBody(env, service, instance, data)

	// 序列化
	body := bytes.NewBuffer(nil)
	err := jsoniter.ConfigCompatibleWithStandardLibrary.NewEncoder(body).Encode(lokiData)
	if err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("序列化日志数据失败: %v\n", err))
		return
	}

	// 编码
	if h.conf.EnableCompress {
		tmp := bytes.NewBuffer(nil)
		err = compress.GetCompress(compress.GzipCompressName).Compress(body, tmp)
		if err != nil {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("编码失败: %v\n", err))
			return
		}
		body = tmp
	}

	// 构建请求
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.conf.ReqTimeout)*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", h.conf.PushAddress, body)
	if err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("构建请求体失败: %v\n", err))
		return
	}
	if h.conf.EnableCompress {
		req.Header.Add("Content-Encoding", compress.GzipCompressName)
	}
	req.Header.Add("Content-Type", "application/json")

	// 请求
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("上报失败, 请求失败: err:%v\n", err))
		return
	}
	defer rsp.Body.Close()
	rspBody, _ := io.ReadAll(rsp.Body)

	// 检查状态码
	if rsp.StatusCode != http.StatusOK && rsp.StatusCode != http.StatusNoContent {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("上报失败, 错误的状态码: code:%v, body:%v\n",
			rsp.StatusCode, string(rspBody)))
		return
	}
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
		conf: conf,
	}
	return h
}

const LokiHttpOutputName = "loki-http"

func init() {
	output.RegistryOutputCreator(LokiHttpOutputName, func(iConfig component.IOutputConfig) output.IOutput {
		return NewHttpOutput(iConfig)
	})
}
