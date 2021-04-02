package dao

import (
	"ChatRoom/Go/common/userinfo"
	"errors"
)

var (
	ERROR_USER_NOTEXIST = errors.New("用户不存在")
	ERROR_USER_EXIST    = errors.New("用户已经存在")
	ERROR_USER_PWD      = errors.New("密码不正确")
)

type IUserDao interface {
	Signup(int, string, string) error
	Signin(int, string) (*userinfo.User, error)
	// Signout(int, string) error
	Delete(int, string) error
	getUserById(int) (*userinfo.User, error)
	IsExist(int) error
}
