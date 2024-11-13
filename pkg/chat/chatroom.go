package chat

import (
	"context"
	"fmt"
	"log"
	"strings"
)

type ChatRoom struct {
	users map[string]*User

	join     chan *Guest
	messages chan Message
	userFail chan UserError
}

func NewChatRoom() ChatRoom {
	return ChatRoom{
		users:    make(map[string]*User),
		join:     make(chan *Guest, EventChannelSize),
		messages: make(chan Message, EventChannelSize),
		userFail: make(chan UserError, EventChannelSize),
	}
}

func (cr *ChatRoom) Run(ctx context.Context) {
	log.Println("ChatRoom started")
	chatRoomCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			log.Println("ChatRoom stopped")
			return
		case guest := <-cr.join:
			cr.handleGuest(chatRoomCtx, guest)
		case message := <-cr.messages:
			cr.handleMessage(message)
		case err := <-cr.userFail:
			cr.handleUserError(err)
		}
	}
}

func (cr *ChatRoom) AddGuest(guest *Guest) {
	cr.join <- guest
}

func (cr *ChatRoom) handleGuest(ctx context.Context, guest *Guest) {
	log.Printf("Guest %s joined, checking name [%s]", guest.Conn.RemoteAddr(), guest.Name)

	if _, present := cr.users[guest.Name]; present {
		guest.Reject("Name already taken. Sorry.")
		return
	}
	user := NewUser(guest.Conn, guest.Name)

	go user.Run(ctx, cr.messages, cr.userFail)

	cr.notifyAboutNewUser(&user)
	cr.users[user.Name] = &user
}

func (cr *ChatRoom) notifyAboutNewUser(user *User) {
	cr.handleMessage(Message{Text: fmt.Sprintf("%s joined", user.Name)})

	sb := strings.Builder{}
	sb.WriteString("* Users in chatroom: ")
	for name := range cr.users {
		sb.WriteString(name)
		sb.WriteString(" ")
	}
	user.Send(sb.String())
}

func (cr *ChatRoom) handleMessage(message Message) {
	log.Println("Message:", message.String())
	for _, user := range cr.users {
		if user.Name == message.From {
			continue
		}
		user.Send(message.String())
	}
}

func (cr *ChatRoom) handleUserError(err UserError) {
	delete(cr.users, err.Name)
	cr.handleMessage(Message{Text: fmt.Sprintf("%s left", err.Name)})
}

type Message struct {
	From string
	Text string
}

func (m Message) String() string {
	if m.From == "" {
		return fmt.Sprintf("* %s", m.Text)
	}
	return fmt.Sprintf("[%s] %s", m.From, m.Text)
}
