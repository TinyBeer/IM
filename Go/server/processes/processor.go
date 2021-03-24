package processes

import (
	"ChatRoom/Go/common/message"
	"ChatRoom/Go/common/utils"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
)

type Processor struct {
	Conn net.Conn
}

// 仅处理用户上限后的业务
func (pro *Processor) serverProcess(mes *message.Message) (err error) {
	switch mes.Type {
	// case message.LoginMesType:
	// 	// 创建UserProcess实例
	// 	up := &UserProcess{Conn: pro.Conn}
	// 	up.ServerProcessLogin(mes)
	// case message.RegisterMesType:
	// 	up := &UserProcess{Conn: pro.Conn}
	// 	// 处理注册消息
	// 	up.ServerProccessRegister(mes)
	case message.SmsMesType:
		smsProcess := &SmsProcess{}
		smsProcess.SendGroupMes(mes)
	case message.LogoutMesType:
		up := &UserProcess{Conn: pro.Conn}
		up.ServerProcessLogout(mes)
		return errors.New("用户登出")
	case message.MessageMesType:
		smsProcess := &SmsProcess{}
		smsProcess.SendMessage(mes)
	default:
		err = errors.New("未知消息类型")
		log.Println("mes=", mes)
	}
	return
}

func (pro *Processor) Process() (err error) {
	// 使用Transfer读写数据
	tf := utils.NewTransfer(pro.Conn)
	// 读取客户发送的消息
	for {
		data, err := tf.ReadDate()
		if err != nil {
			switch err {
			case io.EOF:
				log.Println("客户端断开连接")
			default:
				log.Println("客户端连接中断")
			}
			// 断开连接
			pro.Conn.Close()

			var temp int
			for userID, v := range userMgr.onlineUsers {
				if v.Conn == pro.Conn {
					temp = userID
				}
			}

			if temp != 0 {
				delete(userMgr.onlineUsers, temp)
			}
			return err
		}
		var mes message.Message
		err = json.Unmarshal(data, &mes)
		if err != nil {
			log.Println("json.Unmarshal failed, err=", err.Error())
			continue
		}

		err = pro.serverProcess(&mes)
		if err != nil {
			// 连接断开
			// 重连或者判断为用户离线
			continue
		}
	}
}
