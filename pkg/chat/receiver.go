package chat

import (
	"bufio"
	"fmt"
	"net"
	"unicode"
)

type Receiver struct {
	chatRoom *ChatRoom

	ActorBase
}

func NewReceiver(chatRoom *ChatRoom) *Receiver {
	return &Receiver{
		chatRoom:  chatRoom,
		ActorBase: *NewActorBase("Receiver"),
	}
}

func (r *Receiver) CheckConnection(conn net.Conn) {
	r.actQ <- func() {
		r.checkConnection(conn)
	}
}

func (r *Receiver) checkConnection(conn net.Conn) {
	fmt.Printf("Checking connection for %s\n", conn.RemoteAddr())
	var err error
	defer func() {
		if err != nil {
			fmt.Fprintln(conn, err.Error())
			conn.Close()
		}
	}()

	if _, err := fmt.Fprintln(conn, "Welcome to chat! What is your name?"); err != nil {
		return
	}

	var userName string
	reader := bufio.NewReader(conn)
	userName, err = reader.ReadString('\n')
	if err != nil {
		return
	}
	userName = userName[:len(userName)-1]

	if err = validateUserName(userName); err != nil {
		return
	}

	if err = r.chatRoom.AddUser(conn, userName); err != nil {
		return
	}
}

func (r *Receiver) Stop() {
	if r.stopped {
		return
	}
	r.StopActor()
	fmt.Printf("Receiver stopped\n")
}

func validateUserName(userName string) error {
	if userName == "" {
		return fmt.Errorf("Username can't be empty")
	}
	if len(userName) > 20 {
		return fmt.Errorf("Username can't be longer than 20 characters")
	}
	for _, r := range userName {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return fmt.Errorf("Username can contain only letters and digits")
		}
	}
	return nil
}
