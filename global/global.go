package global

import (
	"github.com/dotdancer/gogofly/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	Redis  redis.UniversalClient
	Config config.Server
	Logger *zap.SugaredLogger
)
