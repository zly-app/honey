package main

import (
	"github.com/zly-app/zapp"
)

func main() {
	honey := NewHoney()

	app := zapp.NewApp("honey",
		zapp.WithEnableDaemon(),         // 守护进程
		honey.LogHook(),                 // 日志hook
		honey.WithCustomComponent(),     // 自定义component
		honey.WithCustomEnableService(), // 自定义启动服务
		honey.WithHoneyService(),        // honey服务
	)

	app.Run()
}
