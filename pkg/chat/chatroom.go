package chat

import (
	"fmt"
	"net"
	"strings"
)

type ChatMessage struct {
	From string
	Text string
}

func (cm ChatMessage) String() string {
	if cm.From == "" {
		return fmt.Sprintf("* %s", cm.Text)
	}
	return fmt.Sprintf("[%s] %s", cm.From, cm.Text)
}

type ChatRoom struct {
	ActorBase
	users map[string]*User
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		ActorBase: *NewActorBase("ChatRoom"),
		users:     make(map[string]*User),
	}
}

func (c *ChatRoom) AddUser(conn net.Conn, userId string) error {
	notifyError := make(chan error, 0)
	// Broadcast the message first, then add the user. This is to not send
	// the message to the user that is joining.
	c.actQ <- func() {
		if _, userExists := c.users[userId]; userExists {
			notifyError <- fmt.Errorf("User %s already exists", userId)
			return
		}
		notifyError <- nil
		text := fmt.Sprintf("%s joined", userId)
		c.broadcast(ChatMessage{From: "", Text: text})
		user := NewUser(conn, userId, c)
		c.addUser(&user)
	}
	return <-notifyError
}

func (c *ChatRoom) addUser(user *User) {
	go user.Run()
	var builder strings.Builder
	builder.WriteString("* Users in chat:")
	i := 0
	for _, user := range c.users {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmt.Sprintf(" %s", user.id))
		i++
	}
	user.Send(builder.String())
	c.users[user.id] = user
	fmt.Printf("User %s added\n", user.id)
}

func (c *ChatRoom) RemoveUser(userId string) {
	// Delete the user first, then broadcast the message. This is to not send
	// the message to the user which has left.
	c.actQ <- func() {
		if _, userExists := c.users[userId]; !userExists {
			return
		}
		delete(c.users, userId)
		text := fmt.Sprintf("%s left", userId)
		c.broadcast(ChatMessage{From: "", Text: text})
	}
}

func (c *ChatRoom) Broadcast(msg ChatMessage) {
	c.actQ <- func() {
		c.broadcast(msg)
	}
}

func (c *ChatRoom) broadcast(msg ChatMessage) {
	text := msg.String()
	for _, user := range c.users {
		if user.id == msg.From {
			continue
		}
		user.Send(text)
	}
}

func (c *ChatRoom) Stop() {
	if c.stopped {
		return
	}
	for _, user := range c.users {
		user.Stop()
	}
	c.StopActor()

	fmt.Printf("ChatRoom stopped\n")
}
