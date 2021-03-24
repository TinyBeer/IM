package router

import (
	"ChatRoom/Web/handlers"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetupRouter() (r *gin.Engine) {
	r = gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.Static("/xxx", "statics")

	checkCookie := func(c *gin.Context) {
		cookie, err := c.Cookie("gin_cookie")
		if err != nil {
			fmt.Println("cookie err")
			c.Abort()
			c.Request.URL.Path = "/"
			r.HandleContext(c)
			return
		}
		c.SetCookie("gin_cookie", cookie, 30, "/", "localhost", false, true)
	}

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	r.GET("/hall", checkCookie, func(c *gin.Context) {
		c.HTML(http.StatusOK, "hall.html", nil)
	})

	r.GET("/cookie", checkCookie, func(c *gin.Context) {
		cookie, err := c.Cookie("gin_cookie")
		if err != nil {
			cookie = "NotSet"

		}

		fmt.Println("cookie value = ", cookie)
	})

	r.GET("/content", checkCookie, func(c *gin.Context) {
		cookie, err := c.Cookie("gin_cookie")
		if err != nil {
			c.Abort()
			return
		}
		userID, _ := strconv.Atoi(cookie)
		handlers.GetContentHandler(c, userID)
	})

	r.POST("/content", checkCookie, func(c *gin.Context) {
		cookie, err := c.Cookie("gin_cookie")
		fmt.Println(cookie)
		if err != nil {
			c.Abort()
			return
		}
		userID, _ := strconv.Atoi(cookie)
		fmt.Println(userID)
		handlers.PostContentHandler(c, userID)
	})

	r.POST("/login", func(c *gin.Context) {
		userID, err := handlers.LoginHandler(c)
		if err != nil {
			return
		}
		fmt.Println(userID)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"res": "fail",
				"err": "",
			})
		} else {
			c.SetCookie("gin_cookie", strconv.Itoa(userID), 100, "/", "localhost", false, true)
			c.JSON(http.StatusOK, gin.H{
				"res": "ok",
			})
		}
	})

	r.POST("/register", func(c *gin.Context) {
		err := handlers.RegisterHandler(c)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"res": "fail",
				"err": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"res": "ok",
			})
		}
	})

	return
}
