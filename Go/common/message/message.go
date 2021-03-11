package message

var (
	// LoginMesType 登录消息类型
	LoginMesType = "LoginMes"
	// LoginResMesType 登录返回消息类型
	LoginResMesType = "LoginMesRes"
	// RegisterMesType 注册
	RegisterMesType = "RegisterMes"
)

// Message ：消息类型
type Message struct {
	Type string `json:"type"` // 消息类型
	Data string `json:"data"` // 消息体
}

// LoginMes :登录消息
type LoginMes struct {
	UserID   int    `json:"userID"`   // 用户ID
	UserPwd  string `json:"userPwd"`  // 用户密码
	UserName string `json:"userName"` // 用户名
}

// LoginResMes : 登录结果返回消息
type LoginResMes struct {
	Code  int    `json:"code"`  // 状态码
	Error string `json:"error"` // 放回错误信息
}

// RegisterMes :
type RegisterMes struct {
}
