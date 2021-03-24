package main

import (
	"ChatRoom/Web/handlers"
	"ChatRoom/Web/router"
)

func main() {
	handlers.UserChin = make(chan handlers.UserConn, 5)
	handlers.DialogList = make(map[int][]string, 5)
	r := router.SetupRouter()

	go func() {
		for {
			userconn := <-handlers.UserChin
			go handlers.Server(userconn.UserID, userconn.Conn)
		}

	}()

	r.Run(":9090")
}
