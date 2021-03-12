package main

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

func initPool(address string, maxIdle, maxActive int, idleTimeout time.Duration) {
	pool = &redis.Pool{
		MaxIdle:     maxIdle,     // 最大空闲数
		MaxActive:   maxActive,   // 最大连接数
		IdleTimeout: idleTimeout, // 最大空闲事件
		Dial: func() (redis.Conn, error) { // 创建连接的函数
			return redis.Dial("tcp", address)
		},
	}
}
