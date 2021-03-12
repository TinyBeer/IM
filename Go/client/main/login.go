package main

import (
	"ChartRoom/common/message"
	"ChartRoom/common/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

func login(userID int, userPwd string) (err error) {
	// 1.连接到服务器
	conn, err := net.Dial("tcp", "192.168.68.166:8889")
	if err != nil {
		return err
	}
	// 延迟断开
	defer conn.Close()

	// 2.准备通过conn发送消息
	var mes message.Message
	mes.Type = message.LoginMesType

	// 3.创建loginMes结构体
	var loginMes message.LoginMes
	loginMes.UserID = userID
	loginMes.UserPwd = userPwd

	// 4.序列化loginMes
	data, err := json.Marshal(&loginMes)
	if err != nil {
		// fmt.Println("json.Marshal failed, err=", err.Error())
		return err
	}

	// 5.填充mes.Data
	mes.Data = string(data)

	// 使用Transfer发送数据

	tf := utils.NewTransfer(conn)
	tf.WritePkg(&mes)

	// 读取客服务端返回的mes
	resMes, err := tf.ReadPkg()

	if err != nil {
		// fmt.Println("err=", err.Error())
		return
	}
	// 反序列化 resMes.Data
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(resMes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		fmt.Println("登陆成功")
		return nil
	} else {
		err = errors.New(loginResMes.Error)
	}
	return
}
