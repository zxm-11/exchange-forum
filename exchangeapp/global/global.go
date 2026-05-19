package global

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// 全局变量定义
var (
	Db          *gorm.DB
	RedisDB *redis.Client
)
