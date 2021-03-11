package main

import (
	"ChartRoom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

func login(userID int, userPwd string) (err error) {
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

	// 6.序列化mes
	data, err = json.Marshal(&mes)
	if err != nil {
		return err
	}

	// 发送data的长度给服务器
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)

	// 发送长度
	n, err := conn.Write(buf[0:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write failed, err=", err.Error())
		return
	}

	fmt.Println("客户端发送消息长度成功")
	fmt.Println("内容:", string(data))

	// 发送消息体
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write failed, err=", err.Error())
		return
	}

	mes, err = readPkg(conn)

	if err != nil {
		fmt.Println("err=", err.Error())
		return
	}
	// 反序列化 mes.Data
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		// fmt.Println("登陆成功")
		return
	} else {
		fmt.Println(loginResMes.Error)
	}
	return
}
