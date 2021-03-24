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

type PreProcessor struct {
	Conn net.Conn
}

// 预处理阶段可以处理的服务类型
const (
	REGISTER_SERVICE = iota
	LOGIN_SERVICE
)

// PreviousProcess完成登录和注册相关任务的处理
// servicType: 与服务类型  注册服务  登陆服务
// err 结果
func (ppro *PreProcessor) PreviousProcess() (servicType int, err error) {
	// 使用Transfer读写数据
	tf := utils.NewTransfer(ppro.Conn)
	// 读取客户发送的消息
	data, err := tf.ReadDate()
	if err != nil {
		switch err {
		case io.EOF:
			log.Println("客户端断开连接")
		default:
			log.Println("客户端连接中断")
		}
		// 断开连接
		ppro.Conn.Close()
		return
	}
	var mes message.Message
	err = json.Unmarshal(data, &mes)
	if err != nil {
		log.Println("json.Unmarshal failed, err=", err.Error())
		return
	}

	servicType, err = ppro.serverProcess(&mes)

	return
}

// serverProcess仅用于处理登录和注册消息
// servicType: 与服务类型  注册服务  登陆服务
// err 结果
func (ppro *PreProcessor) serverProcess(mes *message.Message) (servicType int, err error) {
	// 创建UserProcess实例
	up := &UserProcess{Conn: ppro.Conn}

	// 根据消息类型进行处理
	switch mes.Type {
	case message.LoginMesType:
		// 用户登录
		servicType = LOGIN_SERVICE
		err = up.ServerProcessLogin(mes)
		return
	case message.RegisterMesType:
		// 用户注册
		servicType = REGISTER_SERVICE
		// 处理注册消息
		err = up.ServerProccessRegister(mes)
		return
	default:
		// 消息类型不是注册或者登录
		err = errors.New("不可处理的消息类型")
		log.Println("mes=", mes)
	}
	return
}
