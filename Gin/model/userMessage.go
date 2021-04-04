package model

type UserMes struct {
	Speaker UserInfo `json:"speaker"`
	Content string   `json:"content"`
}
