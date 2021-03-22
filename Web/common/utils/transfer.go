package utils

import (
	"ChatRoom/Web/common/datasafe"
	"encoding/binary"
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

// 传输[]byte 数据
func (tf *Transfer) WriteData(data []byte) (err error) {
	// 传输前进行数据加密
	data, err = datasafe.EncryptoAES(data)
	if err != nil {
		return
	}

	// 发送data的长度给对方
	var pkgLen uint32 = uint32(len(data))
	binary.BigEndian.PutUint32(tf.Buf[0:4], pkgLen)

	// 发送长度
	n, err := tf.Conn.Write(tf.Buf[0:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write failed, err=", err.Error())
		return
	}

	// 发送消息体
	_, err = tf.Conn.Write(data)
	if err != nil {
		return
	}

	return
}

// ReadDate 读取数据 []byte
func (tf *Transfer) ReadDate() (data []byte, err error) {
	// 读取包长度
	_, err = tf.Conn.Read(tf.Buf[:4])
	if err != nil {
		// 读取包长度失败
		return
	}

	// 数据类型转换  获取包大小
	var pkgLen uint32 = binary.BigEndian.Uint32(tf.Buf[:4])

	// 读取消息体
	n, err := tf.Conn.Read(tf.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		return
	}

	// // 开辟新的存储空间
	// data = make([]byte, pkgLen)
	// copy(data, tf.Buf[:pkgLen])

	// 数据接收后，解密
	data, err = datasafe.DecryptoAES(tf.Buf[:pkgLen])
	return
}

// ReadPkg  读取数据包
// func (tf *Transfer) ReadPkg() (mes message.Message, err error) {

// 	// 读取包长度
// 	_, err = tf.Conn.Read(tf.Buf[:4])
// 	if err != nil {
// 		// 读取包长度失败
// 		return
// 	}

// 	// 数据类型转换  获取包大小
// 	var pkgLen uint32 = binary.BigEndian.Uint32(tf.Buf[:4])

// 	// 读取消息体
// 	n, err := tf.Conn.Read(tf.Buf[:pkgLen])
// 	if n != int(pkgLen) || err != nil {
// 		return
// 	}

// 	// 反序列化pkg
// 	err = json.Unmarshal(tf.Buf[:pkgLen], &mes)
// 	if err != nil {
// 		return
// 	}
// 	return
// }

// // 发送数据包
// func (tf *Transfer) WritePkg(mes *message.Message) (err error) {

// 	// 序列化 mes
// 	data, err := json.Marshal(mes)
// 	if err != nil {
// 		// serialization failed
// 		return
// 	}

// 	// 发送data的长度给对方
// 	var pkgLen uint32 = uint32(len(data))
// 	binary.BigEndian.PutUint32(tf.Buf[0:4], pkgLen)

// 	// 发送长度
// 	n, err := tf.Conn.Write(tf.Buf[0:4])
// 	if n != 4 || err != nil {
// 		fmt.Println("conn.Write failed, err=", err.Error())
// 		return
// 	}

// 	// 发送消息体
// 	_, err = tf.Conn.Write(data)
// 	if err != nil {
// 		return
// 	}

// 	return
// }
