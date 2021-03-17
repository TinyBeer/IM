package main

import (
	"ChartRoom/server/model"
	"ChartRoom/server/processes"
	"ChartRoom/server/redisdb"
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("服务器启动...")

	// 初始化连接池
	redisdb.InitPool("localhost:6379", 16, 0, 300*time.Second)
	model.InitUserDao()
	// fmt.Println("服务器[新结构]在8889端口监听...")

	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.Listen failed, err=", err)
		return
	}

	// 延时关闭监听
	defer listen.Close()
	fmt.Println("等待客户端连接服务器...")
	// 循环等待用户连接
	for {

		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept() failed, err=", err.Error())
		}

		// 启动一个协程 与客户端交互
		go func(conn net.Conn) {
			// fmt.Printf("%v\n", conn.RemoteAddr().String())
			defer conn.Close()

			// 创建总控
			processor := &processes.Processor{
				Conn: conn,
			}
			processor.Process()

		}(conn)
	}
}
