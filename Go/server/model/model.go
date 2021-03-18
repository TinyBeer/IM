package model

import (
	"ChartRoom/server/redisdb"
	"time"
)

func InitDao(address string, maxIdle, maxActive int, idleTimeout time.Duration) {
	redisdb.InitPool(address, maxIdle, maxActive, idleTimeout)
	InitUserDao()
}
