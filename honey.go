package main

import (
	"sync"
	"sync/atomic"

	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/service"
	"go.uber.org/zap"

	"github.com/zly-app/honey/component"
	"github.com/zly-app/honey/config"
)

// 默认服务类型
const DefaultServiceType core.ServiceType = "honey"

type Honey struct {
	app         core.IApp
	conf        *config.Config
	rotateGroup *rotateEnvGroup
	mx          sync.Mutex
	state       int32 // 启动状态 0未启动, 1已启动
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
	// 创建旋转器组
	h.rotateGroup = newRotateGroup(h.rotateCreator)
}

func (h *Honey) Inject(a ...interface{}) {}

func (h *Honey) isStart() bool {
	return atomic.LoadInt32(&h.state) == 1
}

func (h *Honey) Start() error {
	atomic.StoreInt32(&h.state, 1)
	return nil
}

func (h *Honey) Close() error {
	atomic.StoreInt32(&h.state, 0)
	// 立即旋转
	rotates := h.rotateGroup.GetAllRotate()
	for _, r := range rotates {
		r.Rotate()
	}
	return nil
}

// 自定义component
func (h *Honey) WithCustomComponent() zapp.Option {
	return zapp.WithCustomComponent(func(app core.IApp) core.IComponent {
		return &component.Component{
			IComponent:   app.GetComponent(),
			LogCollector: h,
		}
	})
}

// 自定义服务
func (h *Honey) WithCustomEnableService() zapp.Option {
	return zapp.WithCustomEnableService(func(app core.IApp, services map[core.ServiceType]bool) {
		// todo 待实现
	})
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
