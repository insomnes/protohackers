package chat

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const EventChannelSize = 16

type ChatServer struct {
	Address string
}

func NewChatServer(address string) ChatServer {
	return ChatServer{
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
	butler := NewButler(&chatRoom)

	ctx, cancel := context.WithCancel(context.Background())
	go chatRoom.Run(ctx)
	go butler.Run(ctx)

	go func() {
		defer ln.Close()

		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println("Error accepting connection:", err)
				return
			}
			log.Printf("Connection from %s\n", conn.RemoteAddr())
			butler.AddConnection(conn)
		}
	}()

	sig := <-sigChan
	log.Printf("Signal received: %v\n", sig)
	cancel()
	<-time.After(1 * time.Second)
}
