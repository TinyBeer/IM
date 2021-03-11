package main

import (
	"ChartRoom/common/message"
	"encoding/binary"
	"encoding/json"
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
