package honey

import (
	"sync"
	"sync/atomic"

	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/core"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/zly-app/honey/log_data"
	"github.com/zly-app/honey/pkg/rotate"
	"github.com/zly-app/honey/zapp_plugin/config"
	"github.com/zly-app/honey/zapp_plugin/reporter"
)

type HoneyPlugin struct {
	app    core.IApp
	conf   *config.Config
	isInit bool  // 是否完成初始化
	state  int32 // 启动状态 0未启动, 1已启动

	rotate      rotate.IRotator // 旋转器
	rotateGPool core.IGPool     // 用于处理同时旋转的协程池

	reporters map[string]reporter.IReporter // 上报器
}

func (h *HoneyPlugin) Init() {
	h.app = zapp.App()
	// 解析配置
	conf := config.NewConfig()
	err := h.app.GetConfig().ParsePluginConfig(nowPluginType, conf, true)
	if err == nil {
		err = conf.Check()
	}
	if err != nil {
		h.app.Fatal("honey插件配置错误", zap.String("PluginType", string(nowPluginType)), zap.Error(err))
	}
	h.conf = conf

	h.MakeRotateGroup()
	h.isInit = true
}

func (h *HoneyPlugin) isStart() bool {
	return atomic.LoadInt32(&h.state) == 1
}

func (h *HoneyPlugin) Start() {
	if !h.isInit {
		return
	}
	atomic.StoreInt32(&h.state, 1)

	// 启动上报者
	h.MakeReporter()
	h.StartReporter()
}

func (h *HoneyPlugin) Close() {
	atomic.StoreInt32(&h.state, 0)
	if !h.isInit {
		return
	}

	// 立即旋转
	h.rotate.Rotate()
	// 关闭上报者
	h.CloseReporter()
}

// 日志拦截函数
func (h *HoneyPlugin) LogInterceptorFunc(ent *zapcore.Entry, fields []zapcore.Field) (cancel bool) {
	log := log_data.MakeLogData(ent, fields)
	h.rotate.Add(log)
	return h.conf.StopLogOutput && h.isStart() // 设置了拦截并且在服务启动后才允许拦截
}

var honey *HoneyPlugin
var once sync.Once

func NewHoneyPlugin() *HoneyPlugin {
	once.Do(func() {
		honey = &HoneyPlugin{}
	})
	return honey
}