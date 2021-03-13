package main

import (
	"ChartRoom/common/message"
	"ChartRoom/common/utils"
	"ChartRoom/server/processes"
	"errors"
	"fmt"
	"io"
	"net"
)

type Processor struct {
	Conn net.Conn
}

func (pro *Processor) serverProcess(mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:
		// 创建UserProcess实例
		up := &processes.UserProcess{Conn: pro.Conn}
		up.ServerProcessLogin(mes)
	case message.RegisterMesType:
		up := &processes.UserProcess{Conn: pro.Conn}
		// 处理注册消息
		up.ServerProccessRegister(mes)
	default:
		err = errors.New("未知消息类型")
	}

	return
}

func (pro *Processor) Process2() (err error) {
	// 使用Transfer读写数据
	tf := utils.NewTransfer(pro.Conn)
	// 读取客户发送的消息
	for {
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端断开连接")
				return err
			} else {
				fmt.Print("readPkg err=", err.Error())
				return err
			}

		}

		// fmt.Println("mes=", mes)

		err = pro.serverProcess(&mes)
		if err != nil {
			fmt.Println("通讯协程错误， err=", err.Error())
		}
	}
}
