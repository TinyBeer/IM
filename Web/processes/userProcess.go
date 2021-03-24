package processes

import (
	"ChatRoom/Web/common/message"
	"ChatRoom/Web/common/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
)

type UserProcess struct {
	// 暂时不需要字段
}

// 完成注册任务
func (up *UserProcess) Register(userID int, userPwd, userName string) (err error) {
	conn, err := net.Dial("tcp", "192.168.68.166:8889")
	if err != nil {
		return err
	}
	// 延迟断开
	defer conn.Close()

	// 2.准备通过conn发送消息
	var mes message.Message
	mes.Type = message.RegisterMesType

	// 3.创建registerMes结构体
	var registerMes message.RegisterMes
	registerMes.UserID = userID
	registerMes.UserPwd = userPwd
	registerMes.UserName = userName

	// 封包
	err = message.Pack(&mes, &registerMes)
	if err != nil {
		return err
	}

	// 序列化
	data, err := json.Marshal(&mes)
	if err != nil {
		return
	}

	// 使用Transfer发送数据
	tf := utils.NewTransfer(conn)
	err = tf.WriteData(data)
	if err != nil {
		fmt.Println("注册消息发送失败")
		return
	}

	resData, err := tf.ReadDate()
	if err != nil {
		log.Println("tf.ReadDate failed, err=", err.Error())
		return
	}
	var resMes message.Message
	err = json.Unmarshal(resData, &resMes)
	if err != nil {
		log.Println("json.Unmarshal failed, err=", err.Error())
		return
	}

	// 解包
	var registerResMes message.RegisterResMes
	err = message.Unpack(&resMes, &registerResMes)
	if err != nil {
		log.Println("Unpack failed, err=", err.Error())
		return
	}

	if registerResMes.Code != 200 {
		err = errors.New(registerResMes.Error)
	}
	return
}

// func (up *UserProcess) Logout() {
// 	// 1.创建mes
// 	var mes message.Message
// 	mes.Type = message.LogoutMesType
// 	// 2.创建logoutMes
// 	var logoutMes message.LogoutMes
// 	logoutMes.UserID = CurUser.User.UserID
// 	// 3.封包
// 	err := message.Pack(&mes, &logoutMes)
// 	if err != nil {
// 		log.Println("Pack failed, err=", err.Error())
// 		return
// 	}

// 	// 序列化
// 	data, err := json.Marshal(&mes)
// 	if err != nil {
// 		return
// 	}

// 	// 4.发送
// 	tf := utils.NewTransfer(CurUser.Conn)
// 	tf.WriteData(data)
// }

func (up *UserProcess) Login(userID int, userPwd string) (conn net.Conn, err error) {
	// 1.连接到服务器
	conn, err = net.Dial("tcp", "192.168.68.166:8889")
	if err != nil {
		return
	}

	// 2.准备通过conn发送消息
	var mes message.Message
	mes.Type = message.LoginMesType

	// 3.创建loginMes结构体
	var loginMes message.LoginMes
	loginMes.AutenticationType = message.PasswordType
	loginMes.UserID = userID
	loginMes.UserPwd = userPwd

	// 4.封包
	err = message.Pack(&mes, &loginMes)
	if err != nil {
		log.Println("pack failed, err=", err.Error())
		return
	}
	// 序列化
	data, err := json.Marshal(&mes)
	if err != nil {
		return
	}

	// 使用Transfer发送数据
	tf := utils.NewTransfer(conn)
	err = tf.WriteData(data)
	if err != nil {
		fmt.Println("登录消息发送失败")
		return
	}

	// 读取返回消息
	resData, err := tf.ReadDate()
	if err != nil {
		log.Println("tf.ReadDate failed, err=", err.Error())
		return
	}

	var resMes message.Message
	err = json.Unmarshal(resData, &resMes)
	if err != nil {
		log.Println("json.Unmarshal failed, err=", err.Error())
		return
	}

	// 解包
	var loginResMes message.LoginResMes
	err = message.Unpack(&resMes, &loginResMes)
	if err != nil {
		fmt.Println("Unpack failed, err=", err.Error())
		return
	}

	if loginResMes.Code != 200 {
		err = errors.New(loginResMes.Error)
	}

	return
}
