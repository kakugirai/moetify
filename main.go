package main

import (
	"github.com/kakugirai/moetify/app"
	"github.com/kakugirai/moetify/config"
)

func main() {
	a := app.App{}
	a.Initialize(config.GetRedisEnv())
	a.Run(config.GetAppEnv().Addr)
}
