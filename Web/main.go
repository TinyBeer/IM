package main

import (
	"ChatRoom/Web/common/message"
	"ChatRoom/Web/common/utils"
	"ChatRoom/Web/processes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	UserID string `json:"userID"`
	PWD    string `json:"password"`
}

var dialogList []string

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("./templates/*")

	r.Static("/xxx", "statics")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	r.GET("/hall", func(c *gin.Context) {
		c.HTML(http.StatusOK, "hall.html", nil)
	})

	r.GET("/content", func(c *gin.Context) {
		if len(dialogList) != 0 {
			c.JSON(http.StatusOK, gin.H{
				"content": dialogList[0],
			})
			dialogList = dialogList[1:]
		} else {
			c.JSON(http.StatusOK, gin.H{
				"content": "",
			})
		}

	})

	r.POST("/login", func(c *gin.Context) {
		userInfo := UserInfo{}

		err := c.Bind(&userInfo)
		if err != nil {
			log.Println(err)
			return
		}

		up := processes.UserProcess{}
		userID, err := strconv.Atoi(userInfo.UserID)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(userID, userInfo)
		conn, err := up.Login(userID, userInfo.PWD)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"res": "fail",
				"err": err.Error(),
			})
		} else {
			go server(conn)
			c.JSON(http.StatusOK, gin.H{
				"res": "ok",
			})
		}
	})

	r.Run(":9090")
}

func server(conn net.Conn) {
	defer conn.Close()
	// 创建一个Transfer 不停的读取消息
	tf := utils.NewTransfer(conn)
	for {
		// fmt.Println("客户端正在读取服务器发送的消息")
		data, err := tf.ReadDate()
		if err != nil {
			log.Println("tf.ReadDate failed, err=", err.Error())
			return
		}
		var mes message.Message
		err = json.Unmarshal(data, &mes)
		if err != nil {
			log.Println("json.Unmarshal failed, err=", err.Error())
			return
		}

		switch mes.Type {
		case message.SmsMesType:
			var smsMes message.SmsMes
			json.Unmarshal([]byte(mes.Data), &smsMes)
			dialogList = append(dialogList, strconv.Itoa(smsMes.UserID)+":"+smsMes.Content)
		default:
			fmt.Println("获取到未知消息类型")
		}
	}
}
