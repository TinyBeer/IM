package main

import (
	"ChartRoom/common/message"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

func readPkg(conn net.Conn) (mes message.Message, err error) {
	buf := make([]byte, 4094)

	fmt.Println("等待读取客户端发送的数据...")
	_, err = conn.Read(buf[:4])
	if err != nil {
		// fmt.Println("conn.Read failed, err=", err.Error())
		return
	}

	fmt.Printf("数据长度：%v\n", buf[:4])
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(buf[:4])

	// 读取消息体
	n, err := conn.Read(buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		// fmt.Println("conn.Read failed, err=", err.Error())
		return
	}

	// 反序列化pkg
	err = json.Unmarshal(buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal failed, err=", err.Error())
		return
	}

	return
}

func writePkg(conn net.Conn, data []byte) (err error) {

	// 发送data的长度给对方
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

	return
}

func serverProcessLogin(conn net.Conn, mes *message.Message) (err error) {
	// 1.先取出mes.Data，并反序列化
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		return
	}

	var resMes message.Message
	resMes.Type = message.LoginResMesType

	var loginResMes message.LoginResMes

	// 2.比对
	if loginMes.UserID == 100 && loginMes.UserPwd == "123456" {
		// 合法
		fmt.Println("登录成功")
		loginResMes.Code = 200 // 登录成功
	} else {
		// 不合法
		loginResMes.Code = 500 // 500 用户不存在
		loginResMes.Error = "该用户不存在"
	}

	// 3.序列化
	data, err := json.Marshal(&loginResMes)
	if err != nil {
		fmt.Println("json.Marshal failed, err=", err)
		return
	}

	// 4.将data赋值给mes.Data
	resMes.Data = string(data)

	// 5.序列化
	data, err = json.Marshal(&resMes)
	if err != nil {
		fmt.Println("json.Marshal failed, err=", err)
		return
	}

	// 6.发送mes
	writePkg(conn, data)

	return
}

func serverProcess(conn net.Conn, mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:
		serverProcessLogin(conn, mes)
	case message.RegisterMesType:
	default:
		err = errors.New("未知消息类型")
	}

	return
}

func process(conn net.Conn) {
	fmt.Printf("接收到%v的连接\n", conn.RemoteAddr().String())
	defer conn.Close()

	// 读取客户发送的消息
	for {
		mes, err := readPkg(conn)
		if err != nil {
			fmt.Println("客户端断开连接")
			return
		}

		fmt.Println("mes=", mes)
		serverProcess(conn, &mes)
	}

}

func main() {
	fmt.Println("服务器启动...")
	fmt.Println("服务器在8889端口监听...")

	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.Listen failed, err=", err)
		return
	}

	// 延时关闭监听
	defer listen.Close()

	// 循环等待用户连接
	for {
		fmt.Println("等待客户端连接服务器...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept() failed, err=", err.Error())
		}

		// 启动一个协程 与客户端交互
		go process(conn)
	}
}
