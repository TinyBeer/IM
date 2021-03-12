package model

import (
	"encoding/json"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

// 服务器启动后初始化一个全局的UserDao
var (
	MyUserDao *UserDao
)

// 定义一个UserDao结构体
type UserDao struct {
	Pool *redis.Pool
}

// 使用工厂模式创建UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		Pool: pool,
	}
	return
}

// 以下是增删改查
// 根据一个用户id返回一个User实例
func (udao *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {
	// 通过给定的id 去redis查询用户
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		// 发生错误
		if err == redis.ErrNil {
			// 没有找到对应id
			err = ERROR_USER_NOTEXIST
		}
		return nil, err
	}
	fmt.Println(res)
	user = &User{}
	// 无错误  将res反序列化为User实例
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("josn.Unmarshal failed, err=", err.Error())
		return
	}
	return
}

// 完成登录校验 Login
// 完成对用户信息的校验
// 如果用户的id或pwd有错误 返回错误信息
func (udao *UserDao) Login(userID int, userPwd string) (user *User, err error) {
	// 从连接池取出连接
	conn := udao.Pool.Get()
	// 延时关闭
	defer conn.Close()

	user, err = udao.getUserById(conn, userID)
	if err != nil {
		return
	}

	// 用获取到了
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}
