package processes

import (
	"ChatRoom/Go/client/model"
	"ChatRoom/Go/common/message"
	"ChatRoom/Go/common/userinfo"
	"fmt"
)

/// 客户端需要维护的map
var onlineUsers map[int]*userinfo.User = make(map[int]*userinfo.User, 10)
var CurUser model.CurUser

// 哭护短显示当前在线用户
func OutputOnlineUsers() {
	// 遍历onlineUsers
	fmt.Println("当前在线用户列表:")
	for id, _ := range onlineUsers {
		fmt.Printf("用户id:%d\n", id)
	}
}

// 编写一个方法  处理返回的 NotifyUsersStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {
	// 查询在线用户列表
	user, ok := onlineUsers[notifyUserStatusMes.UserID]
	if !ok {
		// 用户不在在线用户列表
		if notifyUserStatusMes.UserStatus != message.USER_OFFLINE {
			// 非下线通知 则更向在线用户列表添加用户
			onlineUsers[notifyUserStatusMes.UserID] = &userinfo.User{
				UserID:     notifyUserStatusMes.UserID,
				UserStatus: notifyUserStatusMes.UserStatus,
			}
			fmt.Printf("用户%d上线\n", notifyUserStatusMes.UserID)
		}

	} else {
		// 用户在在线用户列表
		if notifyUserStatusMes.UserStatus == message.USER_OFFLINE {
			// 用户下线 则将用户从在线用户列表中移除
			delete(onlineUsers, notifyUserStatusMes.UserID)
			fmt.Printf("用户%d下线\n", notifyUserStatusMes.UserID)
		} else {
			// 非下线通知 则更新用户状态
			user.UserStatus = notifyUserStatusMes.UserStatus
		}
	}
}
