package cmd

import (
	"fmt"
	"github.com/dotdancer/gogofly/core"
	"github.com/dotdancer/gogofly/global"
	"github.com/dotdancer/gogofly/router"
)

func Start() {
	core.InitConfig()
	global.Logger = core.InitLogger()
	global.DB = core.InitMysql()
	if global.Config.System.UseRedis {
		core.InitRedis()
	}
	router.InitRouter()
}

func Clear() {
	fmt.Println("Program execution stops")
}
