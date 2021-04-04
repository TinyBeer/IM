package controller

import (
	"ChatRoom/Gin/dao"
	"ChatRoom/Gin/datasafe"
	"ChatRoom/Gin/model"
	"ChatRoom/Gin/processes"
	"ChatRoom/Gin/response"
	"fmt"

	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {

	// 获取参数
	var user model.User
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.Abort()
	}
	id, _ := strconv.Atoi(user.ID)
	if id <= 0 || id > 99999 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "ID应少于5位")
		return
	}

	if len(user.Name) > 10 || len(user.Name) == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity,
			422, nil, "昵称不能为空，且不能超过10位")
		return
	}

	if len(user.Password) < 6 || len(user.Password) > 10 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码因为6-10位")
		return
	}

	// 注册

	up := &processes.UserProcess{}
	err = up.Register(id, user.Password, user.Name)
	if err != nil {
		response.Response(ctx, http.StatusBadRequest, 500, nil, err.Error())
		return
	}

	// 发放token
	token, err := datasafe.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		return
	}

	// 缓存
	err = dao.MyUserDao.Insert(user.ID, user.Name)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 返回结果
	response.Success(ctx, gin.H{"token": token}, "注册成功")
}

func Login(ctx *gin.Context) {

	var user model.User
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.Abort()
	}

	id, _ := strconv.Atoi(user.ID)
	if id <= 0 || id > 99999 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "ID应少于5位")
		return
	}

	if len(user.Password) < 6 || len(user.Password) > 10 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码因为6-10位")
		return
	}

	up := &processes.UserProcess{}
	user.Name, err = up.Check(id, user.Password)
	if err != nil {
		response.Response(ctx, http.StatusBadRequest, 500, nil, err.Error())
		return
	}
	// 发放token
	token, err := datasafe.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token generate error:%v", err)
		return
	}

	// 缓存
	err = dao.MyUserDao.Insert(user.ID, user.Name)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 返回结果
	response.Success(ctx, gin.H{"token": token}, "登陆成功")
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	response.Success(ctx, gin.H{"code": 200, "user": user}, "")
}

var contentList []model.UserMes

func Speak(ctx *gin.Context) {
	var userMes model.UserMes
	userInter, _ := ctx.Get("user")
	userMes.Speaker, _ = userInter.(model.UserInfo)
	data, _ := ctx.GetRawData()
	userMes.Content = string(data)

	contentList = append(contentList, userMes)

	for _, content := range contentList {
		fmt.Println(content)
	}

	response.Success(ctx, nil, "发言成功")
}

func GetContent(ctx *gin.Context) {
	str := ctx.Param("id")
	id, _ := strconv.Atoi(str)
	fmt.Println(id)

	if id >= len(contentList) {
		response.Success(ctx, gin.H{"content": nil, "content_id": len(contentList)}, "接收成功")
	} else {
		response.Success(ctx, gin.H{"content": contentList[id:], "content_id": len(contentList)}, "接收成功")
	}

}
