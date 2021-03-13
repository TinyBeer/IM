package model

import (
	"ChartRoom/common/message"
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
func (udao *UserDao) getUserById(conn redis.Conn, id int) (user *message.User, err error) {
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
	user = &message.User{}
	// 无错误  将res反序列化为User实例
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("josn.Unmarshal failed, err=", err.Error())
		return
	}
	return
}

// Register向redis注册用户
func (udao *UserDao) Register(user *message.User) (err error) {
	// 从连接池取出连接
	conn := udao.Pool.Get()
	// 延时关闭
	defer conn.Close()

	_, err = udao.getUserById(conn, user.UserID)
	if err == nil {
		// 用户ID已经存在
		return ERROR_USER_EXIST
	}

	// 该用户ID可用
	data, err := json.Marshal(user)
	if err != nil {
		return
	}

	// 入库
	_, err = conn.Do("HSet", "users", user.UserID, string(data))
	if err != nil {
		fmt.Println("用户信息入库失败")
		return
	}
	return
}

// 完成登录校验 Login
// 完成对用户信息的校验
// 如果用户的id或pwd有错误 返回错误信息
func (udao *UserDao) Login(userID int, userPwd string) (user *message.User, err error) {
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
