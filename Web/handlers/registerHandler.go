package handlers

import (
	"ChatRoom/Web/processes"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) (err error) {
	var temp struct {
		UserID   string `json:"userID"`
		UserPwd  string `josn:"userPwd"`
		UserName string `json:"userName"`
	}
	err = c.Bind(&temp)
	if err != nil {
		fmt.Println(err)
		return
	}

	userID, err := strconv.Atoi(temp.UserID)
	if err != nil {
		return
	}

	up := processes.UserProcess{}
	err = up.Register(userID, temp.UserPwd, temp.UserName)

	return
}
