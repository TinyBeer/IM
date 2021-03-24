package processes

import (
	"ChatRoom/Go/common/message"
	"ChatRoom/Go/common/utils"
	"ChatRoom/Go/server/model"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

// 声明结构体
type SmsProcess struct {
	//..
}

// 擦看并发送离线消息
func (sp *SmsProcess) SendOfflineMessage(userID int, conn net.Conn) (err error) {
	// 获取离线留言
	dataSlice, mesErr := model.MyUserDao.WithdrawOfflineMesById(userID)
	if mesErr != nil {
		log.Println("WithdrawOfflineMesById failed, err=", mesErr.Error())
		return
	}

	// 创建 messageMes
	var mes message.Message
	mes.Type = message.SmsMesType
	var messageMes message.MessageMes
	var smsMes message.SmsMes
	for _, messageString := range dataSlice {
		if err != nil {
			log.Println("Unmarshal failed, err=", mesErr.Error())
			continue
		}
		// 反序列化 messageMes
		mesErr = json.Unmarshal([]byte(messageString), &messageMes)
		if mesErr != nil {
			log.Println("Unmarshal failed, err=", mesErr.Error())
			return
		}

		// 装填smsMes信息
		smsMes.Content = messageMes.Content
		smsMes.UserID = messageMes.UserID

		// 封包
		mesErr = message.Pack(&mes, &smsMes)
		if err != mesErr {
			fmt.Println("Pack failed, err=", err)
			return
		}

		// 发送
		sp.SendMesToEachOnlineUser(&mes, conn)
	}
	return
}

func (sp *SmsProcess) SendMessage(mes *message.Message) (err error) {

	// 1.取出mes.Data,并反序列化
	var messageMes message.MessageMes
	err = message.Unpack(mes, &messageMes)
	if err != nil {
		log.Println("ServerProcessMessage message.Unpack failed, err=", err.Error())
		return
	}

	conn := model.MyUserDao.Pool.Get()
	defer conn.Close()

	// 判断用户是否存在
	_, err = model.MyUserDao.GetUserById(conn, messageMes.ToUserID)
	if err != nil {
		log.Println("GetUserById failed, err", err.Error())
		return
	}

	// 2.判断用户是否在线
	up, ok := userMgr.onlineUsers[messageMes.ToUserID]
	if ok {
		// 4.1 在线  转发消息
		// 构建mes
		var sendMes message.Message
		sendMes.Type = message.SmsMesType

		var smsMes message.SmsMes
		smsMes.Content = messageMes.Content
		smsMes.UserID = messageMes.UserID
		// 封包
		err = message.Pack(&sendMes, &smsMes)
		if err != nil {
			log.Println("ServerProcessMessage message.Pack failed, err=", err.Error())
			return
		}
		// 使用SendMesToEachOnlineUser函数发送sendMes
		sp.SendMesToEachOnlineUser(&sendMes, up.Conn)
	} else {
		// 4.2 不在线 转存消息
		err = model.MyUserDao.DepositUserOfflineMesById(messageMes.ToUserID, []byte(mes.Data))
		if err != nil {
			log.Println("DepositUserOfflineMesById failed, err=", err.Error())
			return
		}
	}
	return
}

// 转发消息
func (sp *SmsProcess) SendGroupMes(mes *message.Message) (err error) {

	// 取出smsMes
	var smsMes message.SmsMes
	err = message.Unpack(mes, &smsMes)
	if err != nil {
		log.Println("Unpack failed, err=", err.Error())
		return
	}

	// 遍历服务端的onlineUsers
	// 转发消息
	for id, up := range userMgr.onlineUsers {
		if id == smsMes.UserID {
			continue
		}
		sp.SendMesToEachOnlineUser(mes, up.Conn)
	}
	return
}

// 发送消息
func (sp *SmsProcess) SendMesToEachOnlineUser(mes *message.Message, conn net.Conn) (err error) {
	tf := utils.NewTransfer(conn)

	// 序列化
	data, err := json.Marshal(&mes)
	if err != nil {
		return
	}

	err = tf.WriteData(data)
	if err != nil {
		fmt.Println("tf.WriteData failed, err=", err.Error())
		return
	}
	return
}
