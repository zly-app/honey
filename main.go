package main

import (
	"github.com/zly-app/zapp"
)

func main() {
	app := zapp.NewApp("honey",
		zapp.WithEnableDaemon(),
	)
	app.Run()
}
