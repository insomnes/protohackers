package chat

import (
	"bufio"
	"fmt"
	"net"
	"time"
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

	if _, err := fmt.Fprintln(conn, "Welcome to chat! What is your name?"); err != nil {
		fmt.Fprintln(conn, err.Error())
		conn.Close()
		return
	}
	go r.receiveUserName(conn)
}

func (r *Receiver) receiveUserName(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(1 * time.Minute))
	reader := bufio.NewReader(conn)
	userName, err := reader.ReadString('\n')
	if err != nil {
		conn.Close()
		return
	}
	userName = userName[:len(userName)-1]
	r.AddUserToChat(conn, userName)
}

func (r *Receiver) AddUserToChat(conn net.Conn, userName string) error {
	notifyError := make(chan error, 0)
	r.actQ <- func() {
		notifyError <- r.addUserToChat(conn, userName)
	}
	return <-notifyError
}

func (r *Receiver) addUserToChat(conn net.Conn, userName string) error {
	if err := validateUserName(userName); err != nil {
		return err
	}
	return r.chatRoom.AddUser(conn, userName)
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
