package main

import (
	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/core"

	_ "github.com/zly-app/honey/input/http"

	_ "github.com/zly-app/honey/output/honey-http"
	_ "github.com/zly-app/honey/output/loki-http"
	_ "github.com/zly-app/honey/output/std"
)

func main() {
	honey := NewHoney()

	app := zapp.NewApp("honey",
		zapp.WithEnableDaemon(),  // 守护进程
		honey.LogHook(),          // 日志hook
		honey.WithHoneyService(), // honey服务
		zapp.WithHandler(zapp.AfterExitHandler, func(app core.IApp, handlerType zapp.HandlerType) {
			honey.AfterExit()
		}),
	)

	app.Run()
}
