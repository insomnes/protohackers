package chat

import (
	"bufio"
	"fmt"
	"net"
)

//     Supervisor
//    / del    \ error
// ChatRoom -> User

type User struct {
	id       string
	conn     net.Conn
	chatRoom *ChatRoom

	ActorBase
}

func NewUserOut(conn net.Conn, userId string, chatRoom *ChatRoom) User {
	return User{
		id:        userId,
		conn:      conn,
		chatRoom:  chatRoom,
		ActorBase: *NewActorBase("User"),
	}
}

func (u *User) Run() {
	go u.receive()
	u.ActorBase.Run()
}

func (u *User) Send(msg string) {
	u.actQ <- func() {
		u.send(msg)
	}
}

func (u *User) send(msg string) {
	if _, err := fmt.Fprintln(u.conn, msg); err != nil {
		u.chatRoom.RemoveUser(u.id)
		u.Stop()
	}
}

func (u *User) receive() {
	reader := bufio.NewReader(u.conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		u.chatRoom.Broadcast(ChatMessage{From: u.id, Text: msg[:len(msg)-1]})
	}
	fmt.Printf("User %s left\n", u.id)
	u.chatRoom.RemoveUser(u.id)
	u.Stop()
}

func (u *User) Stop() {
	if u.stopped {
		return
	}

	u.StopActor()
	u.conn.Close()

	fmt.Printf("User %s conn closed\n", u.id)
}
