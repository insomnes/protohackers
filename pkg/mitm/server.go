package mitm

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

type MitmServer struct {
	Address  string
	ChatAddr string
}

func NewMitmServer(address string, chatAddr string) MitmServer {
	return MitmServer{
		Address:  address,
		ChatAddr: chatAddr,
	}
}

func (ms *MitmServer) Run() {
	ln, err := net.Listen("tcp", ms.Address)
	if err != nil {
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer ln.Close()

		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println("Error accepting connection:", err)
				return
			}
			log.Printf("Connection from %s\n", conn.RemoteAddr())
			go RunMitmProxy(ctx, conn, ms.ChatAddr)
		}
	}()

	sig := <-sigChan
	log.Printf("Signal received: %v\n", sig)
	cancel()
	<-time.After(300 * time.Millisecond)
}
