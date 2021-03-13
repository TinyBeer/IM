package message

type User struct {
	// 为了序列化和反序列化  确保tag正确
	UserID     int    `json:"userID"`
	UserName   string `json:"userName"`
	UserPwd    string `json:"userPwd"`
	UserStatus int    `josn:"userStatus"`
}
