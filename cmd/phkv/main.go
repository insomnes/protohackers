package main

import (
	"flag"
	"fmt"

	"github.com/insomnes/protohackers/pkg/kvstore"
)

const (
	defaultHost = "127.0.0.1"
	defaultPort = 9999
)

func main() {
	host := flag.String("host", defaultHost, "address to listen on")
	port := flag.Uint("port", defaultPort, "port to listen on 1-65535")
	flag.Parse()
	address := fmt.Sprintf("%s:%d", *host, *port)

	server := kvstore.NewKVServer(address)
	server.Run()
}
