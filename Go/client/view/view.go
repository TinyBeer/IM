package view

import (
	"fmt"
)

// View :sturct using to display client window
type View struct {
}

// MainMenu :display main menu
func (*View) MainMenu() {
	fmt.Println("------------欢迎登录海量用户聊天系统------------")
	fmt.Println("\t\t1：登录聊天室")
	fmt.Println("\t\t2：注册用户")
	fmt.Println("\t\t3：退出系统")
	fmt.Print("请选择1-3：")
}
