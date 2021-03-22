package message

import "ChatRoom/Web/common/userinfo"

// 定义用户在线状态
const (
	USER_ONLINE = iota
	USER_OFFLINE
	USER_BUSY
)

var (
	// LoginMesType 登录消息类型
	LoginMesType = "LoginMes"
	// LoginResMesType 登录返回消息类型
	LoginResMesType = "LoginMesRes"
	// RegisterMesType 注册
	RegisterMesType = "RegisterMes"
	// RegisterMesResType 注册结果返回消息类型
	RegisterMesResType = "RegisterMesRes"
	// NotifyUserStatusMes 推送用户登录消息
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	// SmsMesType 短消息类型
	SmsMesType = "SmsMes"
	// LogoutMesType 登出消息
	LogoutMesType = "LogoutMes"
	// MessageMesType 留言消息
	MessageMesType = "MessageMes"
	// MessageResMes 留言结果
	MessageResMesType = "MessageResMes"
)

// 登录方式
var (
	// 账号密码方式
	PasswordType = "Password"
)

// Message ：消息类型
type Message struct {
	Type string `json:"type"` // 消息类型
	Data string `json:"data"` // 消息体
}

// LoginMes :登录消息
type LoginMes struct {
	AutenticationType string `json:"autenticationType"`
	userinfo.AuthenticationUserInfo
}

// LoginResMes : 登录结果返回消息
type LoginResMes struct {
	Code int `json:"code"` // 状态码
	userinfo.PersonalUserInfo
	OnlineUsersID []int  `json:"onlineUsers"` // 在线用户
	Error         string `json:"error"`       // 返回错误信息
}

// RegisterMes :注册消息
type RegisterMes struct {
	userinfo.RegisterUserInfo
}

// RegisterResMes
type RegisterResMes struct {
	Code  int    `json:"code"`  // 400 用户已存在  200 注册成功
	Error string `json:"error"` // 错误提示
}

// 配合服务器推送上线通知
type NotifyUserStatusMes struct {
	userinfo.BasicUserInfo
	UserStatus int `json:"userStatus"`
}

// SmsMes 增加结构体
type SmsMes struct {
	userinfo.BasicUserInfo
	Content string `json:"content"`
}

// LogoutMes
type LogoutMes struct {
	userinfo.BasicUserInfo // 匿名结构体
}

// MessageMes 留言
type MessageMes struct {
	userinfo.BasicUserInfo        // 匿名结构体
	ToUserID               int    `json:"toUserID"`
	Content                string `json:"content"`
}

// MessageResMes 留言结果
type MessageResMes struct {
	Code  int    `json:"code"`  // 200 留言成功 300 留言失败
	Error string `json:"error"` // 错误提示
}
