package handlers

import (
	"ChatRoom/Web/processes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 请求数据
func GetContentHandler(c *gin.Context, userID int) {
	if len(DialogList[userID]) != 0 {
		c.JSON(http.StatusOK, gin.H{
			"content": DialogList[userID][0],
		})
		DialogList[userID] = DialogList[userID][1:]
	} else {
		c.JSON(http.StatusOK, gin.H{
			"content": "",
		})
	}
	return
}

// 发送消息
func PostContentHandler(c *gin.Context, userID int) (err error) {
	var temp struct {
		Content string `json:"content"`
	}
	err = c.ShouldBind(&temp)
	if err != nil {
		return
	}
	fmt.Println(temp)
	sp := processes.SmsProcess{}
	sp.SendGroupMes(temp.Content, userID, Conns[userID])
	return
}
