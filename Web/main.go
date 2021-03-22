package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	UserID string `json:"userID"`
	PWD    string `json:"password"`
}

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
	var flag = true
	r.GET("/content", func(c *gin.Context) {
		var content string
		if flag {
			content = "你好呀" + time.Now().Format("2006-01-02 15:04:05")
		}

		flag = !flag

		c.JSON(http.StatusOK, gin.H{
			"content": content,
		})

	})

	r.POST("/login", func(c *gin.Context) {
		userInfo := UserInfo{}

		err := c.Bind(&userInfo)
		if err != nil {
			log.Println(err)
			return
		}

		if userInfo.UserID == "100" && userInfo.PWD == "123456" {
			c.JSON(http.StatusOK, gin.H{
				"res": "ok",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"res": "fail",
			})
		}
	})

	r.Run(":9090")
}
