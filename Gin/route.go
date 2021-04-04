package main

import (
	"ChatRoom/Gin/controller"
	"ChatRoom/Gin/middleware"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) *gin.Engine {

	r.Use(middleware.CORSMiddleware())
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)
	r.POST("/api/auth/speak", middleware.AuthMiddleware(), controller.Speak)
	r.GET("/api/auth/content/:id", middleware.AuthMiddleware(), controller.GetContent)

	return r
}
