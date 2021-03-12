package main

import (
	"ChartRoom/client/processes"
	"fmt"
)

var (
	key     int
	userID  int
	userPwd string
	loop    bool
)

func main() {
	loop = true

	for loop {
		fmt.Println("------------欢迎登录海量用户聊天系统------------")
		fmt.Println("\t\t1：登录聊天室")
		fmt.Println("\t\t2：注册用户")
		fmt.Println("\t\t3：退出系统")
		fmt.Print("请选择1-3：")

		fmt.Scanln(&key)
		switch key {
		case 1:
			fmt.Println("------登录聊天室")
			fmt.Println("请输入用户ID：")
			fmt.Scanln(&userID)
			fmt.Println("请输入密码：")
			fmt.Scanln(&userPwd)
			up := &processes.UserProcess{}
			up.Login(userID, userPwd)
		case 2:
			fmt.Println("------注册用户")
		case 3:
			fmt.Println("------退出聊天室")
			loop = false
		default:
			fmt.Print("您的输入有误，请重新输入：")
		}

		// if key == 1 {
		// 	fmt.Println("请输入用户ID：")
		// 	fmt.Scanln(&userID)
		// 	fmt.Println("请输入密码：")
		// 	fmt.Scanln(&userPwd)
		// 	err := login(userID, userPwd)
		// 	if err != nil {
		// 		fmt.Println(err.Error())
		// 	} else {
		// 		fmt.Println("登录成功1")
		// 	}
		// } else if key == 2 {
		// 	fmt.Println("进行用户注册处理")
		// }
	}
}
