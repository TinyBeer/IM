package message

import (
	"encoding/json"
)

func Pack(mes *Message, mess interface{}) (err error) {
	// 1.序列化loginMes
	data, err := json.Marshal(mess)
	if err != nil {
		return
	}

	// 2.填充 mes 的Data
	mes.Data = string(data)
	return
}

func Unpack(mes *Message, mess interface{}) (err error) {
	err = json.Unmarshal([]byte(mes.Data), mess)
	if err != nil {
		return
	}
	return
}
