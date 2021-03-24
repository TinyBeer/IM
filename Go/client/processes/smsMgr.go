package processes

import (
	"ChatRoom/Go/common/message"
	"fmt"
)

// 传入smsMes类型数据
func outputMes(mes *message.Message) {
	// 反序列化
	var smsMes message.SmsMes
	err := message.Unpack(mes, &smsMes)
	if err != nil {
		fmt.Println("Unpack failed, err=", err.Error())
		return
	}

	fmt.Printf("收到来自用户%d的消息:\n", smsMes.UserID)
	fmt.Println(smsMes.Content)
}
