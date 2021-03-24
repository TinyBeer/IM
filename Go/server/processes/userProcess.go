package processes

import (
	"ChatRoom/Go/common/message"
	"ChatRoom/Go/common/utils"
	"ChatRoom/Go/server/model"

	"encoding/json"
	"fmt"
	"log"
	"net"
)

type UserProcess struct {
	// 连接
	Conn net.Conn
	// 用户ID
	UserID int
}

// 处理登录mes
func (up *UserProcess) ServerProcessLogout(mes *message.Message) (err error) {
	// 1.先取出mes.Data，并反序列化
	var logoutMes message.LogoutMes
	err = message.Unpack(mes, &logoutMes)
	if err != nil {
		return
	}
	// 删除map中断 对应数据
	delete(userMgr.onlineUsers, logoutMes.UserID)
	// 通知用户下线
	up.NotifyOthersOffline(logoutMes.UserID)
	fmt.Printf("用户%d下线\n", logoutMes.UserID)
	return
}

// 编写通知用户下线的方法
func (up *UserProcess) NotifyOthersOffline(userId int) (err error) {
	// 遍历onlineUsers  一个一个发送消息
	for id, up := range userMgr.onlineUsers {
		if id == userId {
			continue
		}
		// 开始通知下线
		up.NotifyMeOffline(userId)
	}
	return
}

func (up *UserProcess) NotifyMeOffline(userId int) {
	// 开始组装NotifyUserStatusMes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserID = userId
	notifyUserStatusMes.UserStatus = message.USER_OFFLINE

	// 封包
	err := message.Pack(&mes, &notifyUserStatusMes)
	if err != nil {
		fmt.Println("Pack failed, err=", err.Error())
		return
	}

	// 传输mes
	tf := utils.NewTransfer(up.Conn)

	// 序列化
	data, err := json.Marshal(&mes)
	if err != nil {
		return
	}

	err = tf.WriteData(data)
	if err != nil {
		fmt.Println("tf.WriteData failed, err=", err.Error())
		return
	}
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

	// 封包
	err := message.Pack(&mes, &notifyUserStatusMes)
	if err != nil {
		fmt.Println("Pack failed, err=", err.Error())
		return
	}

	data, err := json.Marshal(&mes)
	if err != nil {
		return
	}

	// 传输mes
	tf := utils.NewTransfer(up.Conn)
	err = tf.WriteData(data)
	if err != nil {
		fmt.Println("tf.WriteData failed, err=", err.Error())
		return
	}
}

// 处理注册mes
func (up *UserProcess) ServerProccessRegister(mes *message.Message) (err error) {
	var registerMes message.RegisterMes
	err = message.Unpack(mes, &registerMes)
	if err != nil {
		return
	}

	// 声明返回消息
	var resMes message.Message
	resMes.Type = message.RegisterMesResType
	// 声明注册返回消息体
	var registerResMes message.RegisterResMes

	// 进行注册
	err = model.MyUserDao.Register(&registerMes.RegisterUserInfo)
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

	// 3.封包
	err = message.Pack(&resMes, &registerResMes)
	if err != nil {
		fmt.Println("Pack failed, err=", err)
		return
	}

	// 序列化
	data, err := json.Marshal(&resMes)
	if err != nil {
		return
	}

	// 使用Transfer返回resMes
	tf := utils.NewTransfer(up.Conn)
	tf.WriteData(data)
	return
}

// 处理登录mes
func (up *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	// 1.先取出mes.Data，并反序列化
	var loginMes message.LoginMes
	err = message.Unpack(mes, &loginMes)
	if err != nil {
		return
	}

	var resMes message.Message
	resMes.Type = message.LoginResMesType
	var loginResMes message.LoginResMes

	// 2.比对数据库
	// 到redis数据库进行验证
	switch loginMes.AutenticationType {
	case message.PasswordType:
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

			_, ok := userMgr.onlineUsers[loginMes.UserID]

			if ok {
				loginResMes.Code = 501 // 用户已登录
				loginResMes.Error = "用户已登录"
			} else {
				loginResMes.Code = 200 // 登录成功
				// 用户登录成功  更行onlineUsers
				// 为up加入UserID
				up.UserID = loginMes.UserID
				loginResMes.UserName = user.UserName
				userMgr.AddOnlineUser(up)
				fmt.Printf("用户%d登录\n", loginMes.UserID)
				// 通知其他用户上线
				up.NotifyOthersOnline(loginMes.UserID)

				// fmt.Println(user)
				for id, _ := range userMgr.onlineUsers {
					if id == user.UserID {
						continue
					}
					loginResMes.OnlineUsersID = append(loginResMes.OnlineUsersID, id)
				}
			}

		}
	default:
		return
	}

	// 3.封包
	err = message.Pack(&resMes, &loginResMes)
	if err != nil {
		log.Println("Pack failed, err=", err)
		return
	}

	// 序列化
	data, err := json.Marshal(&resMes)
	if err != nil {
		return
	}

	// 使用Transfer返回resMes
	tf := utils.NewTransfer(up.Conn)
	tf.WriteData(data)

	// 登录确认消息发送后执行
	// 登录成功则发送离线留言
	if loginResMes.Code == 200 {
		smsProcess := &SmsProcess{}
		mesErr := smsProcess.SendOfflineMessage(up.UserID, up.Conn)
		if mesErr != nil {
			log.Println("SendOfflineMessage failed, err=", mesErr.Error())
		}
	}

	return
}
