package cmd

import (
	"fmt"

	"github.com/dotdancer/gogofly/core"
	"github.com/dotdancer/gogofly/global"
	"github.com/dotdancer/gogofly/model" // Added for model access
	"github.com/dotdancer/gogofly/router"
)

func Start() {
	core.InitConfig()
	global.Logger = core.InitLogger()
	global.DB = core.InitMysql() // Uncommented to initialize DB

	if global.DB != nil {
		global.Logger.Info("Attempting to migrate database tables...")
		// Register table migrations
		err := global.DB.AutoMigrate(
			&model.UserInfo{},
			&model.Tenable{},
			&model.Account{},
			&model.Transaction{},
			&model.Category{},
			&model.Budget{}, // Add Budget model for migration
		)
		if err != nil {
			global.Logger.Error("Failed to migrate database tables: " + err.Error())
			// Consider how to handle this error in a production environment.
			// For now, logging the error and continuing.
		} else {
			global.Logger.Info("Database tables migrated successfully or no changes needed.")
		}
	} else {
		global.Logger.Warn("Database not initialized (global.DB is nil), skipping migrations.")
	}

	// 只在配置启用Redis时初始化Redis连接
	if global.Config.System.UseRedis {
		global.Logger.Info("Initializing Redis connection...")
		core.InitRedis()
	} else {
		global.Logger.Info("Redis is disabled in config, skipping Redis initialization.")
	}

	router.InitRouter()
}

func Clear() {
	fmt.Println("Program execution stops")
}
