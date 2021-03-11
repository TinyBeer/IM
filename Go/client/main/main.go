package main

import (
	"ChartRoom/client/view"
	"fmt"
)

var (
	key     int
	userID  int
	userPwd string
	loop    bool
)

func main() {
	view := view.View{}
	loop = true

	for loop {
		view.MainMenu()
		fmt.Scanln(&key)

		switch key {
		case 1:
			fmt.Println("------登录聊天室")
		case 2:
			fmt.Println("------注册用户")
		case 3:
			fmt.Println("------退出聊天室")
			loop = false
		default:
			fmt.Print("您的输入有误，请重新输入：")
		}

		if key == 1 {
			fmt.Println("请输入用户ID：")
			fmt.Scanln(&userID)
			fmt.Println("请输入密码：")
			fmt.Scanln(&userPwd)
			err := login(userID, userPwd)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("登录成功")
			}
		} else if key == 2 {
			fmt.Println("进行用户注册处理")
		}
	}
}
