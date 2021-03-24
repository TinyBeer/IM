package model

import (
	"ChatRoom/Go/common/userinfo"
	"net"
)

// 创建全局
type CurUser struct {
	userinfo.User
	Conn net.Conn
}
