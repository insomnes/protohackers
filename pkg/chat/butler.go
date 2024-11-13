package chat

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"time"
	"unicode"
)

type Butler struct {
	chatRoom    *ChatRoom
	namedGuests chan *Guest
	connections chan net.Conn
}

func NewButler(chatRoom *ChatRoom) Butler {
	return Butler{
		chatRoom:    chatRoom,
		namedGuests: make(chan *Guest, EventChannelSize),
		connections: make(chan net.Conn, EventChannelSize),
	}
}

func (b *Butler) Run(ctx context.Context) {
	log.Println("Butler started")
	butlerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			log.Println("Butler stopped")
			return
		case conn := <-b.connections:
			guest := NewGuest(conn)
			go guest.Greet(butlerCtx, b.namedGuests)
		case guest := <-b.namedGuests:
			b.chatRoom.AddGuest(guest)
		}
	}
}

func (b *Butler) AddConnection(conn net.Conn) {
	b.connections <- conn
}

type Guest struct {
	Name string
	Conn net.Conn
}

func NewGuest(conn net.Conn) *Guest {
	return &Guest{
		Name: "",
		Conn: conn,
	}
}

func (g *Guest) Greet(ctx context.Context, guests chan<- *Guest) {
	done := make(chan struct{}, 0)
	go g.greet(done)
	select {
	case <-ctx.Done():
		log.Printf("External stop for %s\n", g.Conn.RemoteAddr())
		g.Close()
		return
	case <-done:
		log.Printf("Accepted %s as %s\n", g.Conn.RemoteAddr(), g.Name)
		guests <- g
	}
}

func (g *Guest) greet(done chan<- struct{}) {
	err := g.send("Welcome! What is your name?")
	if err != nil {
		log.Printf("Error writing to %s: %v\n", g.Conn.RemoteAddr(), err)
		g.Close()
		return
	}

	g.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	reader := bufio.NewReader(g.Conn)
	name, err := reader.ReadString('\n')
	if err != nil {
		if err.Error() == "EOF" {
			log.Printf("%s closed connection\n", g.Conn.RemoteAddr())
		} else {
			log.Printf("Guest cannot read from %s: %v\n", g.Conn.RemoteAddr(), err)
		}
		g.Close()

		return
	}

	name = name[:len(name)-1]
	err = validateName(name)
	if err != nil {
		g.Reject(fmt.Sprintf("Invalid name %s: %v", name, err))
		return
	}

	g.Name = name
	g.Conn.SetReadDeadline(time.Time{})

	close(done)
}

func (g *Guest) Reject(reason string) {
	log.Printf("Rejecting %s: %s\n", g.Conn.RemoteAddr(), reason)
	g.send(reason)
	g.Close()
}

func (g *Guest) Close() {
	g.Conn.Close()
}

func (g *Guest) send(text string) error {
	_, err := g.Conn.Write([]byte(text + "\n"))
	if err != nil {
		log.Printf("Error writing to %s: %v\n", g.Conn.RemoteAddr(), err)
		g.Conn.Close()
		return err
	}

	return nil
}

func validateName(name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if len(name) > 20 {
		return fmt.Errorf("name cannot be longer than 32 characters")
	}
	for _, c := range name {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return fmt.Errorf("name can only contain letters and digits")
		}
	}
	return nil
}
