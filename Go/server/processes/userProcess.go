package processes

import (
	"ChartRoom/common/message"
	"ChartRoom/common/utils"
	"ChartRoom/server/model"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	Conn net.Conn
}

// 处理登录mes
func (up *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	// 1.先取出mes.Data，并反序列化
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		return
	}

	var resMes message.Message
	resMes.Type = message.LoginResMesType

	var loginResMes message.LoginResMes

	// 2.比对
	// 到redis数据库进行验证
	user, err := model.MyUserDao.Login(loginMes.UserID, loginMes.UserPwd)
	if err != nil {
		switch err {
		case model.ERROR_USER_NOTEXIST:
			loginResMes.Code = 500 // 500 用户不存在
			loginResMes.Error = err.Error()
		case model.ERROR_USER_PWD:
			loginResMes.Code = 403 // 403 密码不正确
			loginResMes.Error = err.Error()
		default:
			loginResMes.Code = 505 // 505 服务器内部错误
			loginResMes.Error = "服务器内部错误"
		}
	} else {
		loginResMes.Code = 200 // 登录成功
		fmt.Println(user)
	}

	// 3.序列化
	data, err := json.Marshal(&loginResMes)
	if err != nil {
		fmt.Println("json.Marshal failed, err=", err)
		return
	}

	// 4.将data赋值给mes.Data
	resMes.Data = string(data)

	// 使用Transfer返回resMes
	tf := utils.NewTransfer(up.Conn)
	tf.WritePkg(&resMes)
	return
}
