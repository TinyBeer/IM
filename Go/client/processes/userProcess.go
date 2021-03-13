package processes

import (
	"ChartRoom/common/message"
	"ChartRoom/common/utils"
	"encoding/json"
	"fmt"
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
	registerMes.User.UserID = userID
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	// 4.序列化loginMes
	data, err := json.Marshal(&registerMes)
	if err != nil {
		// fmt.Println("json.Marshal failed, err=", err.Error())
		return err
	}

	// 5.填充mes.Data
	mes.Data = string(data)

	// 使用Transfer发送数据
	tf := utils.NewTransfer(conn)
	err = tf.WritePkg(&mes)
	if err != nil {
		fmt.Println("注册消息发送失败")
		return
	}

	// 读取客服务端返回的mes
	resMes, err := tf.ReadPkg()
	if err != nil {
		return
	}

	// 反序列化 resMes.Data
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(resMes.Data), &registerResMes)
	if err != nil {
		fmt.Println("json.Umarshal failed, err=", err.Error())
		return
	}

	if registerResMes.Code == 200 {
		fmt.Println("注册成功，来请登录吧")
	} else {
		fmt.Println(registerResMes.Error)
	}
	return
}

func (up *UserProcess) Login(userID int, userPwd string) (err error) {
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
	err = tf.WritePkg(&mes)
	if err != nil {
		fmt.Println("登录消息发送失败")
		return
	}

	// 读取客服务端返回的mes
	resMes, err := tf.ReadPkg()
	if err != nil {
		// fmt.Println("err=", err.Error())
		return
	}
	// 反序列化 resMes.Data
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(resMes.Data), &loginResMes)
	if err != nil {
		fmt.Println("json.Umarshal failed, err=", err.Error())
		return
	}
	if loginResMes.Code == 200 {
		// 可以显示当前在线用户id列表
		fmt.Println("当前在线用户列表如下:")
		for _, onlineUserID := range loginResMes.OnlineUsersID {
			fmt.Printf("用户id:%d\n", onlineUserID)
			// 初始化onlineUsers
			user := &message.User{
				UserID:     onlineUserID,
				UserStatus: message.USER_ONLINE,
			}
			onlineUsers[onlineUserID] = user
		}
		fmt.Println()
		// 启一个协程保持和服务器的练习
		go serverMesProcess(conn)

		// 1.显示二级菜单
		for {
			ShowMenu(loginResMes.UserName)
		}

	} else {
		fmt.Println(loginResMes.Error)
	}
	return
}
