package processes

import (
	"ChartRoom/common/message"
	"fmt"
)

/// 客户端需要维护的map
var onlineUsers map[int]*message.User = make(map[int]*message.User, 10)

// 哭护短显示当前在线用户
func outputOnlineUsers() {
	// 遍历onlineUsers
	fmt.Println("当前在线用户列表:")
	for id, _ := range onlineUsers {
		fmt.Printf("用户id:%d\n", id)
	}
}

// 编写一个方法  处理返回的 NotifyUsersStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {
	user, ok := onlineUsers[notifyUserStatusMes.UserID]

	if !ok {
		onlineUsers[notifyUserStatusMes.UserID] = &message.User{
			UserID:     notifyUserStatusMes.UserID,
			UserStatus: notifyUserStatusMes.UserStatus,
		}
	} else {
		user.UserStatus = notifyUserStatusMes.UserStatus
	}

}
