package main

import (
	"log"
	"net/http"

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
