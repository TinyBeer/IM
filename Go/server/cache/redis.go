package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	redisPool *redis.Pool
	ErrNil    = errors.New("nil returned")
)

func InitPool(addr string, port int, maxIdle, maxActive int, idleTimeout time.Duration) {
	address := fmt.Sprintf("%s:%d", addr, port)
	redisPool = &redis.Pool{
		MaxIdle:     maxIdle,     // 最大空闲数
		MaxActive:   maxActive,   // 最大连接数
		IdleTimeout: idleTimeout, // 最大空闲事件
		Dial: func() (redis.Conn, error) { // 创建连接的函数
			return redis.Dial("tcp", address)
		},
	}
	fmt.Println(redisPool.Get().Do("ping"))
}

func RedisDel(args ...interface{}) error {
	// 从连接池取出连接
	conn := redisPool.Get()
	// 延时关闭连接
	defer conn.Close()

	_, err := conn.Do("del", args...)
	return err
}

func RedisGetList(key string) ([]string, error) {
	// 从连接池取出连接
	conn := redisPool.Get()
	// 延时关闭连接
	defer conn.Close()

	return redis.Strings(conn.Do("lrange", key, 0, -1))
}

func RedisLpush(args ...interface{}) error {
	// 从连接池取出连接
	conn := redisPool.Get()
	// 延时关闭连接
	defer conn.Close()

	_, err := conn.Do("lpsh", args...)
	return err

}

func RedisHGetStr(args ...interface{}) (string, error) {
	// 从连接池取出连接
	conn := redisPool.Get()
	// 延时关闭连接
	defer conn.Close()
	res, err := redis.String(conn.Do("HGet", args...))

	if err == redis.ErrNil {
		return "", ErrNil
	}
	return res, err
}

func RedisHSet(args ...interface{}) error {
	// 从连接池取出连接
	conn := redisPool.Get()
	// 延时关闭连接
	defer conn.Close()
	_, err := conn.Do("HSet", args...)
	return err
}
