package main

import (
	"flag"
	"fmt"

	"github.com/insomnes/protohackers/pkg/mitm"
)

const (
	defaultChatAddr = "chat.protohackers.com:16963"
	defaultHost     = "127.0.0.1"
	defaultPort     = 9999
)

func main() {
	chatAddr := flag.String("chat", defaultChatAddr, "chat server address string")
	host := flag.String("host", defaultHost, "address to listen on")
	port := flag.Uint("port", defaultPort, "port to listen on 1-65535")
	flag.Parse()
	address := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Println("chat address:", *chatAddr)
	fmt.Println("address:", address)
	chatServer := mitm.NewMitmServer(address, *chatAddr)
	chatServer.Run()
}
