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
	if err != nil {
		fmt.Println("json.Umarshal failed, err=", err.Error())
		return
	}
	if loginResMes.Code == 200 {

		// 启一个协程保持和服务器的练习
		go serverMesProcess(conn)

		// 1.显示二级菜单
		for {
			ShowMenu()
		}

	} else {
		fmt.Println(loginResMes.Error)
	}
	return
}

// 和服务器保持通信
func serverMesProcess(conn net.Conn) {
	// 创建一个Transfer 不停的读取消息
	tf := utils.NewTransfer(conn)
	for {
		fmt.Println("客户端正在读取服务器发送的消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg failed, err=", err.Error())
			return
		}

		fmt.Printf("mes = %v\n", mes)
	}
}
