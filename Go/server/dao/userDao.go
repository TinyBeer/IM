package dao

import (
	"ChatRoom/Go/common/userinfo"
	"ChatRoom/Go/server/cache"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// 服务器启动后初始化一个全局的UserDao
var (
	MyUserDao = &RedisUserDao{}
)

type RedisUserDao struct {
}

func (udao *RedisUserDao) DepositUserOfflineMesById(id int, data []byte) (err error) {
	// 将数据存入mesList[userId]中
	err = cache.RedisLpush("mesList"+strconv.Itoa(id), string(data))
	if err != nil {
		return err
	}
	// 退出
	return
}

func (udao *RedisUserDao) WithdrawOfflineMesById(id int) (dataSlice []string, err error) {
	// 将数据存入mesList[userId]中
	dataSlice, err = cache.RedisGetList("mesList" + strconv.Itoa(id))
	log.Println(dataSlice)
	if err != nil {
		return
	}

	// 如果留言数量不为零
	if len(dataSlice) != 0 {
		err = cache.RedisDel("mesList" + strconv.Itoa(id))
		if err != nil {
			log.Println(err.Error())
		}
	}

	// 退出
	return
}

func (rud *RedisUserDao) Signup(id int, pwd string, name string) error {
	_, err := rud.getUserById(id)
	if err != ERROR_USER_NOTEXIST {
		return err
	}

	// 存储前 用户id加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("用户密码加密失败，err:", err)
		return err
	}

	user := userinfo.User{
		UserID:   id,
		UserPwd:  string(hashedPassword),
		UserName: name,
	}

	// 该用户ID可用
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// 入库
	err = cache.RedisHSet("users", user.UserID, string(data))
	if err != nil {
		fmt.Println("用户信息入库失败")
		return err
	}
	return nil

}

func (rud *RedisUserDao) Signin(id int, pwd string) (*userinfo.User, error) {
	user, err := rud.getUserById(id)
	if err != nil {
		return nil, err
	}

	// 判断密码是否正确
	if err = bcrypt.CompareHashAndPassword([]byte(user.UserPwd), []byte(pwd)); err != nil {
		err = ERROR_USER_PWD
		return nil, err
	}
	return user, nil
}

func (rud *RedisUserDao) Delete(id int, pwd string) error {
	user, err := rud.getUserById(id)
	if err != nil {
		return err
	}

	// 判断密码是否正确
	if err = bcrypt.CompareHashAndPassword([]byte(user.UserPwd), []byte(pwd)); err != nil {
		err = ERROR_USER_PWD
		return err
	}

	return cache.RedisDel("users", id)
}

func (rud *RedisUserDao) IsExist(id int) bool {
	if _, err := rud.getUserById(id); err != nil {
		return false
	}

	return true
}

func (rud *RedisUserDao) getUserById(id int) (user *userinfo.User, err error) {
	// 通过给定的id 去redis查询用户
	res, err := cache.RedisHGetStr("users", id)
	fmt.Println(res, err)
	if err != nil {
		// 发生错误
		if err == cache.ErrNil {
			// 没有找到对应id
			err = ERROR_USER_NOTEXIST
		}
		return nil, err
	}
	// fmt.Println(res)
	user = &userinfo.User{}
	// 无错误  将res反序列化为User实例
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("josn.Unmarshal failed, err=", err.Error())
		return
	}
	return
}
