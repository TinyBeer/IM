package main

import (
	"ChatRoom/Go/client/processes"
	"ChatRoom/Go/client/view"
	"fmt"
	"net"
)

var (
	key      int
	userID   int
	userPwd  string
	userName string
	cChan    chan net.Conn
	pMgr     *view.PageMgr
)

func UserLogin() (conn net.Conn, err error) {
	fmt.Println("请输入用户ID：")
	fmt.Scanln(&userID)
	fmt.Println("请输入密码：")
	fmt.Scanln(&userPwd)
	up := &processes.UserProcess{}
	conn, err = up.Login(userID, userPwd)
	return
}

func UserRegister() (err error) {
	fmt.Println("请输入用户ID：")
	fmt.Scanln(&userID)
	fmt.Println("请输入密码：")
	fmt.Scanln(&userPwd)
	fmt.Println("请输入用户昵称：")
	fmt.Scanln(&userName)
	// 调用UserDao实例  实现注册
	up := &processes.UserProcess{}
	err = up.Register(userID, userPwd, userName)
	return
}
