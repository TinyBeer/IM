package handlers

import (
	"ChatRoom/Web/common/message"
	"ChatRoom/Web/common/utils"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
)

var DialogList map[int][]string

func Server(userID int, conn net.Conn) {
	defer conn.Close()
	// 创建一个Transfer 不停的读取消息
	tf := utils.NewTransfer(conn)
	for {
		data, err := tf.ReadDate()

		if err != nil {
			log.Println("tf.ReadDate failed, err=", err.Error())
			return
		}
		var mes message.Message
		err = json.Unmarshal(data, &mes)
		fmt.Println(mes)
		if err != nil {
			log.Println("json.Unmarshal failed, err=", err.Error())
			return
		}

		switch mes.Type {
		case message.SmsMesType:
			var smsMes message.SmsMes
			json.Unmarshal([]byte(mes.Data), &smsMes)
			content := strconv.Itoa(smsMes.UserID) + ":" + smsMes.Content
			fmt.Println(userID, content)
			DialogList[userID] = append(DialogList[userID], content)
		default:
			fmt.Println("获取到未知消息类型")
		}
	}
}
