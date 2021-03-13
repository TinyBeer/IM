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
	// 连接
	Conn net.Conn
	// 用户ID
	UserID int
}

// 编写通知用户在线的方法
func (up *UserProcess) NotifyOthersOnline(userId int) (err error) {
	// 遍历onlineUsers  一个一个发送消息
	for id, up := range userMgr.onlineUsers {
		if id == userId {
			continue
		}

		// 开始通知上线
		up.NotifyMeOnline(userId)
	}
	return
}

func (up *UserProcess) NotifyMeOnline(userId int) {

	// 开始组装NotifyUserStatusMes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserID = userId
	notifyUserStatusMes.UserStatus = message.USER_ONLINE

	// 序列换
	data, err := json.Marshal(&notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal failed, err=", err.Error())
		return
	}

	// 装填到mes包
	mes.Data = string(data)

	// 传输mes
	tf := utils.NewTransfer(up.Conn)

	err = tf.WritePkg(&mes)
	if err != nil {
		fmt.Println("tf.WritePkg failed, err=", err.Error())
		return
	}
}

// 处理注册mes
func (up *UserProcess) ServerProccessRegister(mes *message.Message) (err error) {
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		return
	}

	// 声明返回消息
	var resMes message.Message
	resMes.Type = message.RegisterMesResType
	// 声明注册返回消息体
	var registerResMes message.RegisterResMes

	// 进行注册
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		switch err {
		case model.ERROR_USER_EXIST:
			registerResMes.Code = 505
			registerResMes.Error = err.Error()
		default:
			registerResMes.Code = 506
			registerResMes.Error = "注册时发生未知错误"
		}
	} else {
		registerResMes.Code = 200
	}

	// 3.序列化
	data, err := json.Marshal(&registerResMes)
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
		// 用户登录成功  更行onlineUsers
		// 为up加入UserID
		up.UserID = loginMes.UserID
		userMgr.AddOnlineUser(up)
		fmt.Printf("用户%d登录\n", loginMes.UserID)
		// 通知其他用户上线
		up.NotifyOthersOnline(loginMes.UserID)

		loginResMes.UserName = user.UserName

		// fmt.Println(user)
		for id, _ := range userMgr.onlineUsers {
			if id == user.UserID {
				continue
			}
			loginResMes.OnlineUsersID = append(loginResMes.OnlineUsersID, id)
		}

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
