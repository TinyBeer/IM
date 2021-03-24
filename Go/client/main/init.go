package main

import (
	"ChatRoom/Go/client/processes"
	"ChatRoom/Go/client/view"
	"fmt"
	"log"
	"net"
	"os"
)

func InitSystem() {
	cChan = make(chan net.Conn, 1)
	go func() {
		for {
			conn := <-cChan
			go processes.ServerMesProcess(conn)
		}
	}()
}

func InitView() *view.PageMgr {
	var content string
	var toUserID int
	smsProcess := &processes.SmsProcess{}
	pMgr = view.NewPageMgr()
	p := pMgr.AddPage("MainPage", "", "------------欢迎登录海量用户聊天系统------------", "")
	p.AddOption("\t\t登录聊天室", func() {
		conn, err := UserLogin()
		p.SetDescription(fmt.Sprintf("------恭喜%s登陆成功------", processes.CurUser.UserName))
		if err != nil {
			log.Println(err.Error())
			return
		} else {
			// 启一个协程保持和服务器的连接
			processes.OutputOnlineUsers()
			cChan <- conn
			pMgr.TurnToPage("HallPage")
		}
	})

	p.AddOption("\t\t注册用户", func() {
		err := UserRegister()
		if err != nil {
			fmt.Println("用户注册失败,err=", err.Error())
			return
		} else {
			fmt.Println("注册成功，快来登录吧！")
		}
	})

	p.AddOption("\t\t退出系统", func() {
		os.Exit(0)
	})

	p = pMgr.AddPage("HallPage", "------聊天室大厅界面------", "恭喜xxx登录成功", "MainPage")

	p.AddOption("\t在线用户列表", func() {
		processes.OutputOnlineUsers()
	})
	p.AddOption("\t群发消息", func() {
		fmt.Println("请输入要发送的消息:")
		fmt.Scanf("%s\n", &content)
		smsProcess.SendGroupMes(content)
	})
	p.AddOption("\t信息列表", func() {

	})

	p.AddOption("\t发送消息", func() {
		fmt.Println("请输入要给用户的ID:")
		fmt.Scanln(&toUserID)
		fmt.Println("请输入要发送的消息:")
		fmt.Scanln(&content)
		smsProcess.SendMessageToUser(toUserID, content)
	})

	p.AddOption("\t退出聊天室", func() {
		up := &processes.UserProcess{}
		up.Logout()
		pMgr.GoBack()
	})
	return pMgr
}
