package main

import (
	"ChartRoom/server/model"
	"fmt"
	"net"
	"time"
)

func process(conn net.Conn) {
	fmt.Printf("接收到%v的连接\n", conn.RemoteAddr().String())
	defer conn.Close()

	// 创建总控
	processor := &Processor{
		Conn: conn,
	}
	processor.Process2()

}

func initUserDao() {
	model.MyUserDao = model.NewUserDao(pool)
}

func main() {
	fmt.Println("服务器启动...")

	// 初始化连接池
	initPool("localhost:6379", 16, 0, 300*time.Second)
	initUserDao()
	fmt.Println("服务器[新结构]在8889端口监听...")

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
