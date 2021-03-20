package userinfo

// 为了序列化和反序列化  确保tag正确

// 最基础的用户信息
type BasicUserInfo struct {
	UserID int `json:"userID"`
}

// 鉴权信息
type AuthenticationUserInfo struct {
	BasicUserInfo
	UserPwd string `json:"userPwd"`
}

// 个性化信息
type PersonalUserInfo struct {
	UserName string `json:"userName"`
}

type RegisterUserInfo struct {
	AuthenticationUserInfo
	PersonalUserInfo
}

type User struct {
	UserID     int    `json:"userID"`
	UserPwd    string `json:"userPwd"`
	UserName   string `json:"userName"`
	UserStatus int    `josn:"userStatus"`
}
