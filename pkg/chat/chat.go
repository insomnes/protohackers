package chat

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type ChatServer struct {
	Address string
}

func NewChatServer(address string) *ChatServer {
	return &ChatServer{
		Address: address,
	}
}

func (cs *ChatServer) Run() {
	ln, err := net.Listen("tcp", cs.Address)
	if err != nil {
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	chatRoom := NewChatRoom()
	receiver := NewReceiver(chatRoom)
	director := NewDirector(chatRoom, receiver)
	director.Run()
	defer director.Stop()

	go func() {
		defer ln.Close()

		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err)
				return
			}
			fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr())
			receiver.CheckConnection(conn)
		}
	}()

	sig := <-sigChan
	fmt.Printf("Signal received: %v\n", sig)
}
