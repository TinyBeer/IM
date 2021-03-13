package processes

import (
	"ChartRoom/common/message"
	"ChartRoom/common/utils"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

// 显示登录成功后的界面
func ShowMenu(userName string) {

	fmt.Printf("----恭喜%3s登录成功-----\n", userName)
	fmt.Println("----1：在线用户列表-----")
	fmt.Println("----2：发送消息---------")
	fmt.Println("----3：信息列表---------")
	fmt.Println("----4：退出系统---------")
	fmt.Print("请选择1-4：")

	var key int
	fmt.Scanln(&key)
	switch key {
	case 1:
		// fmt.Println("显示在线用户列表")
		outputOnlineUsers()
	case 2:
		fmt.Println("发送消息")
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Print("退出系统")
		os.Exit(0)
	default:
		fmt.Println("输入有误，从新输入：")

	}

}

// 和服务器保持通信
func serverMesProcess(conn net.Conn) {
	// 创建一个Transfer 不停的读取消息
	tf := utils.NewTransfer(conn)
	for {
		// fmt.Println("客户端正在读取服务器发送的消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg failed, err=", err.Error())
			return
		}

		// fmt.Printf("mes = %v\n", mes)

		switch mes.Type {
		case message.NotifyUserStatusMesType:
			// 处理用户上上线消息
			// 取出NotifyUserStatusMes
			var notifyUserStatusMes message.NotifyUserStatusMes
			// 反序列化
			err = json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			if err != nil {
				fmt.Println("json.Unmarshall failed, err=", err.Error())
				continue
			}
			updateUserStatus(&notifyUserStatusMes)
			outputOnlineUsers()
		default:
			fmt.Println("获取到未知消息类型")
		}

	}
}
