package http

import (
	"github.com/zly-app/service/api"
	"github.com/zly-app/service/api/config"
	"github.com/zly-app/zapp/core"
	"go.uber.org/zap"

	"github.com/zly-app/honey/component"
	"github.com/zly-app/honey/pkg/compress"
	"github.com/zly-app/honey/pkg/serializer"
)

type HttpInput struct {
	c          component.IComponent
	api        *api.ApiService
	compress   compress.ICompress
	serializer serializer.ISerializer
}

func (h *HttpInput) Start() error {
	return h.api.Start()
}

func (h *HttpInput) Close() error {
	return h.api.Close()
}

func NewHttpInput(c component.IComponent) *HttpInput {
	conf := newConfig()
	err := c.ParseInputConf(HttpInputName, conf, true)
	if err == nil {
		err = conf.Check()
	}
	if err != nil {
		c.Fatal("获取http输入器配置失败", zap.Error(err))
	}
	h := &HttpInput{
		c:          c,
		compress:   compress.GetCompress(conf.Compress),
		serializer: serializer.GetSerializer(conf.Serializer),
	}

	apiConf := config.NewConfig()
	apiConf.Bind = conf.Bind
	apiConf.PostMaxMemory = conf.PostMaxMemory
	apiConf.SendDetailedErrorInProduction = true

	opts := []api.Option{
		api.WithMiddleware(AuthTokenMiddleware(conf.AuthToken)),
	}
	apiService := api.NewApiService(c.App(), apiConf, opts...)
	apiService.RegistryRouter(func(c core.IComponent, router api.Party) {
		router.Post(conf.ReportPath, api.Wrap(h.Receive))
	})
	h.api = apiService

	return h
}
