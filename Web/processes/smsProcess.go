package processes

import (
	"ChatRoom/Web/common/message"
	"ChatRoom/Web/common/utils"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

// import (
// 	"ChatRoom/Web/common/message"
// 	"ChatRoom/Web/common/utils"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// )

type SmsProcess struct {
}

// // SendMesToUser toUserID:接收方ID content 内容
// func (sp *SmsProcess) SendMessageToUser(toUserID int, content string) (err error) {
// 	// 1.创建mes MessageMes
// 	var mes message.Message
// 	mes.Type = message.MessageMesType
// 	var messageMes message.MessageMes

// 	// 2.装在messageMes数据
// 	messageMes.UserID = CurUser.User.UserID
// 	messageMes.ToUserID = toUserID
// 	messageMes.Content = content

// 	// 3.序列化messageMes
// 	data, err := json.Marshal(&messageMes)
// 	if err != nil {
// 		log.Println("SendMessage json.Marshal failed err=", err)
// 		return
// 	}
// 	// 4.装在mes.Data
// 	mes.Data = string(data)

// 	// 5.序列化mes
// 	data, err = json.Marshal(&mes)
// 	if err != nil {
// 		log.Println("SendMessage json.Marshal failed err=", err.Error())
// 		return
// 	}
// 	// 6.使用Transfer发送数据
// 	tf := utils.NewTransfer(CurUser.Conn)
// 	err = tf.WriteData(data)
// 	if err != nil {
// 		log.Println("SendMessage tf.Write failed, err=", err.Error())
// 		return
// 	}
// 	return
// }

// 发送群发消息
func (sp *SmsProcess) SendGroupMes(content string, userID int, conn net.Conn) (err error) {
	// 创建一个mes
	var mes message.Message
	mes.Type = message.SmsMesType

	var smsMes message.SmsMes
	smsMes.UserID = userID
	smsMes.Content = content

	// 序列化smsMes
	err = message.Pack(&mes, &smsMes)
	if err != nil {
		fmt.Println("Pack failed, err=", err.Error())
		return
	}

	// 序列化
	data, err := json.Marshal(&mes)
	if err != nil {
		return
	}

	// 发送
	tf := utils.NewTransfer(conn)
	err = tf.WriteData(data)
	if err != nil {
		log.Println("tf.WriteData failed, err=", err.Error())
		return
	}
	return
}
