package utils

import (
	"ChartRoom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

// 数据传输去
type Transfer struct {
	Conn net.Conn
	Buf  [4096]byte
}

// 工厂模式生成传输器
func NewTransfer(conn net.Conn) *Transfer {
	return &Transfer{Conn: conn}
}

// ReadPkg  读取数据包
func (tf *Transfer) ReadPkg() (mes message.Message, err error) {

	// 创建数据缓存buf
	// buf := make([]byte, 4094)

	// fmt.Println("等待读取客户端发送的数据...")

	// 读取包长度
	_, err = tf.Conn.Read(tf.Buf[:4])
	if err != nil {
		// 读取包长度失败
		// fmt.Println("conn.Read failed, err=", err.Error())
		return
	}

	// fmt.Printf("数据长度：%v\n", buf[:4])
	// 数据类型转换  获取包大小
	var pkgLen uint32 = binary.BigEndian.Uint32(tf.Buf[:4])

	// 读取消息体
	n, err := tf.Conn.Read(tf.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		// fmt.Println("conn.Read failed, err=", err.Error())
		return
	}

	// 反序列化pkg
	err = json.Unmarshal(tf.Buf[:pkgLen], &mes)
	if err != nil {
		// fmt.Println("json.Unmarshal failed, err=", err.Error())
		return
	}

	return
}

// 发送数据包
func (tf *Transfer) WritePkg(mes *message.Message) (err error) {

	// 序列化 mes
	data, err := json.Marshal(mes)
	if err != nil {
		// serialization failed
		return
	}

	// 发送data的长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)

	// 发送长度
	n, err := tf.Conn.Write(buf[0:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write failed, err=", err.Error())
		return
	}

	// fmt.Println("客户端发送消息长度成功")
	// fmt.Println("内容:", string(data))

	// 发送消息体
	_, err = tf.Conn.Write(data)
	if err != nil {
		// fmt.Println("conn.Write failed, err=", err.Error())
		return
	}

	return
}
