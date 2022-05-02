package http

import (
	"github.com/zly-app/service/api"
	"github.com/zly-app/service/api/config"
	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zly-app/honey/component"
)

type HttpInput struct {
	c   component.ILogCollector
	api *api.ApiService
}

func (h *HttpInput) Start() error {
	return h.api.Start()
}

func (h *HttpInput) Close() error {
	return h.api.Close()
}

func NewHttpInput(c component.ILogCollector, iConfig component.IInputConfig) *HttpInput {
	conf := newConfig()
	err := iConfig.ParseInputConf(HttpInputName, conf, true)
	if err == nil {
		err = conf.Check()
	}
	if err != nil {
		logger.Log.Fatal("获取http输入器配置失败", zap.Error(err))
	}
	h := &HttpInput{
		c: c,
	}

	apiConf := config.NewConfig()
	apiConf.Bind = conf.Bind
	apiConf.PostMaxMemory = conf.PostMaxMemory
	apiConf.SendDetailedErrorInProduction = true

	opts := []api.Option{
		api.WithMiddleware(AuthTokenMiddleware(conf.AuthToken)),
	}
	apiService := api.NewApiService(zapp.App(), apiConf, opts...)
	apiService.RegistryRouter(func(c core.IComponent, router api.Party) {
		router.Post(conf.PushPath, api.Wrap(h.Receive))
	})
	h.api = apiService

	return h
}
