package main

import (
	"sync/atomic"

	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/service"
	"go.uber.org/zap"

	"github.com/zly-app/honey/config"
	"github.com/zly-app/honey/input"
	"github.com/zly-app/honey/output"
)

// 默认服务类型
const DefaultServiceType core.ServiceType = "honey"

type Honey struct {
	app   core.IApp
	conf  *config.Config // 配置
	state int32          // 启动状态 0未启动, 1已启动

	rotateGroup *rotateEnvGroup // 旋转组
	rotateGPool core.IGPool     // 用于处理同时旋转的协程池

	inputs  map[string]input.IInput   // 输入设备
	outputs map[string]output.IOutput // 输出设备
}

func (h *Honey) Init() {
	h.app = zapp.App()
	// 解析配置
	conf := config.NewConfig()
	err := h.app.GetConfig().ParseServiceConfig(DefaultServiceType, conf, true)
	if err == nil {
		err = conf.Check()
	}
	if err != nil {
		h.app.Fatal("honey配置错误", zap.String("ServiceType", string(DefaultServiceType)), zap.Error(err))
	}
	h.conf = conf

	h.MakeRotateGroup()
}

func (h *Honey) Inject(a ...interface{}) {}

func (h *Honey) isStart() bool {
	return atomic.LoadInt32(&h.state) == 1
}

func (h *Honey) Start() error {
	atomic.StoreInt32(&h.state, 1)
	// 启动输入设备
	h.MakeInput()
	h.StartInput()
	// 启动输出设备
	h.MakeOutput()
	h.StartOutput()
	return nil
}

func (h *Honey) Close() error {
	atomic.StoreInt32(&h.state, 0)
	// 关闭输入设备
	h.CloseInput()

	// 立即旋转
	rotates := h.rotateGroup.GetAllRotate()
	for _, r := range rotates {
		r.Rotate()
	}
	return nil
}

func (h *Honey) AfterExit() {
	// 立即旋转
	rotates := h.rotateGroup.GetAllRotate()
	for _, r := range rotates {
		r.Rotate()
	}

	// 等待处理
	h.rotateGPool.Wait()

	// 关闭输出设备
	h.CloseOutput()
}

// 启用honey服务
func (h *Honey) WithHoneyService() zapp.Option {
	service.RegisterCreatorFunc(DefaultServiceType, func(app core.IApp) core.IService {
		return h
	})
	return zapp.WithService(DefaultServiceType)
}

func NewHoney() *Honey {
	return &Honey{}
}
