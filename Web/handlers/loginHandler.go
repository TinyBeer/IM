package handlers

import (
	"ChatRoom/Web/processes"
	"fmt"
	"net"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserConn struct {
	UserID int
	Conn   net.Conn
}

var UserChin chan UserConn
var Conns map[int]net.Conn

func LoginHandler(c *gin.Context) (userID int, err error) {
	// data, _ := c.GetRawData()
	// fmt.Println(string(data))
	// 取出json格式数据
	var temp struct {
		UserID  string `json:"userID"`
		UserPwd string `josn:"userPwd"`
	}
	err = c.Bind(&temp)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 转换
	userID, _ = strconv.Atoi(temp.UserID)
	up := processes.UserProcess{}
	conn, err := up.Login(userID, temp.UserPwd)
	if err != nil {
		return
	}

	UserChin <- UserConn{userID, conn}

	if Conns == nil {
		Conns = make(map[int]net.Conn, 5)
	}
	Conns[userID] = conn
	return
}
