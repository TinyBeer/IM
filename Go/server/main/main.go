package main

import (
	"ChatRoom/Go/server/model"
	"ChatRoom/Go/server/processes"
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("服务器启动...")
	// 初始化Dao功能
	model.InitDao("localhost:6379", 16, 0, 300*time.Second)

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
		// 未登录状态下的服务
		ppro := &processes.PreProcessor{
			Conn: conn,
		}
		sType, err := ppro.PreviousProcess()
		if err != nil {
			// 服务失败
			continue
		}

		// 将用户登录注册  和 上线后业务分离
		// 可以更加方便的调节负载
		switch sType {
		case processes.LOGIN_SERVICE:
			// 如果是登录
			// 启动一个协程 与客户端交互
			go func(conn net.Conn) {
				// fmt.Printf("%v\n", conn.RemoteAddr().String())
				defer conn.Close()

				// 创建总控
				processor := &processes.Processor{
					Conn: conn,
				}
				processor.Process()
				// 传入ppro中的连接 避免连接信息变化
			}(ppro.Conn)
		// case processes.REGISTER_SERVICE:
		default:
			// 否则断开连接
			conn.Close()
			continue
		}

	}
}
