package middleware

import (
	"ChatRoom/Gin/dao"
	"ChatRoom/Gin/datasafe"
	"ChatRoom/Gin/response"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取authorization header
		tokenString := ctx.GetHeader("Authorization")
		// validate token formate
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}

		tokenString = tokenString[7:]

		token, claims, err := datasafe.ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}

		// 验证通过后获取claim中的userId
		userId := fmt.Sprintf("%d", claims.UserId)
		user, err := dao.MyUserDao.GetUserByID(userId)
		if err != nil {
			response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
			ctx.Abort()
			return
		}
		// 用户存在 将user信息写入上下文
		ctx.Set("user", *user)
		ctx.Next()
	}
}
