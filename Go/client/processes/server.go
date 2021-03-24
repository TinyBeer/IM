package processes

import (
	"ChatRoom/Go/common/message"
	"ChatRoom/Go/common/utils"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

// 和服务器保持通信
func ServerMesProcess(conn net.Conn) {
	defer conn.Close()
	// 创建一个Transfer 不停的读取消息
	tf := utils.NewTransfer(conn)
	for {
		// fmt.Println("客户端正在读取服务器发送的消息")
		data, err := tf.ReadDate()
		if err != nil {
			log.Println("tf.ReadDate failed, err=", err.Error())
			return
		}
		var mes message.Message
		err = json.Unmarshal(data, &mes)
		if err != nil {
			log.Println("json.Unmarshal failed, err=", err.Error())
			return
		}

		switch mes.Type {
		case message.NotifyUserStatusMesType:
			// 处理用户状态更新消息
			var notifyUserStatusMes message.NotifyUserStatusMes
			err = message.Unpack(&mes, &notifyUserStatusMes)
			if err != nil {
				log.Println("Unpack failed, err=", err.Error())
				continue
			}
			updateUserStatus(&notifyUserStatusMes)
		case message.SmsMesType:
			outputMes(&mes)
		default:
			fmt.Println("获取到未知消息类型")
		}

	}
}
