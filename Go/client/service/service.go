package service

import (
	"ChatRoom/Go/client/processes"
	"fmt"
)

// 通过ServiceMgr统一管理客户端操作
type ServiceMgr struct {
	key      int
	userID   int
	userPwd  string
	userName string
}

// 登录服务
func (sMgr *ServiceMgr) LoginService() {
	fmt.Println("\t------登录聊天室------")
	fmt.Println("\t请输入用户ID：")
	fmt.Scanln(&sMgr.userID)
	fmt.Println("\t请输入密码：")
	fmt.Scanln(&sMgr.userPwd)
	up := &processes.UserProcess{}
	up.Login(sMgr.userID, sMgr.userPwd)
}

// 注册服务
func (sMgr *ServiceMgr) RegisterService() {
	fmt.Println("\t------注册用户")
	fmt.Println("\t请输入用户ID：")
	fmt.Scanln(&sMgr.userID)
	fmt.Println("\t请输入密码：")
	fmt.Scanln(&sMgr.userPwd)
	fmt.Println("\t请输入用户昵称：")
	fmt.Scanln(&sMgr.userName)
	// 调用UserDao实例  实现注册
	up := &processes.UserProcess{}
	err := up.Register(sMgr.userID, sMgr.userPwd, sMgr.userName)
	if err != nil {
		fmt.Println("注册失败")
	} else {
		fmt.Println("注册成功，快来登录吧!")
	}

}
