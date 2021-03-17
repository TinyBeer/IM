package redisdb

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

var Pool *redis.Pool

func InitPool(address string, maxIdle, maxActive int, idleTimeout time.Duration) {
	Pool = &redis.Pool{
		MaxIdle:     maxIdle,     // 最大空闲数
		MaxActive:   maxActive,   // 最大连接数
		IdleTimeout: idleTimeout, // 最大空闲事件
		Dial: func() (redis.Conn, error) { // 创建连接的函数
			return redis.Dial("tcp", address)
		},
	}
}
