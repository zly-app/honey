package main

import (
	"github.com/zly-app/service/api"
	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/core"

	"github.com/zly-app/honey/component"
)

func main() {
	app := zapp.NewApp("honey",
		zapp.WithEnableDaemon(), // 守护进程
		zapp.WithCustomComponent(func(app core.IApp) core.IComponent {
			return &component.Component{
				IComponent: app.GetComponent(),
			}
		}),
		api.WithService(),
		zapp.WithCustomEnableService(func(app core.IApp, services map[core.ServiceType]bool) {

		}),
	)

	app.Run()
}
