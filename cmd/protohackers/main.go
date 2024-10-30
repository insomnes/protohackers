package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/insomnes/protohackers/pkg/echo"
)

func RunServer(addr string, connHandler func(conn net.Conn)) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		err := fmt.Errorf("failed to listen: %w", err)
		return err
	}

	fmt.Println("Running server on", addr)
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			err := fmt.Errorf("failed to accept: %w", err)
			return err
		}

		go connHandler(conn)
	}
}

func GetHandler(handle string) (func(conn net.Conn), error) {
	switch handle {
	case "echo":
		return echo.EchoHandler, nil
	default:
		return nil, fmt.Errorf("unknown handler: %s", handle)
	}
}

func main() {
	host := flag.String("host", "localhost", "address to listen on")
	port := flag.Uint("port", 9999, "port to listen on")
	handler := flag.String("handler", "echo", "handler to use")
	flag.Parse()

	handlerFunc, err := GetHandler(*handler)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting handler: %w", err)
		os.Exit(1)
	}

	fullAddr := fmt.Sprintf("%s:%d", *host, *port)
	if err := RunServer(fullAddr, handlerFunc); err != nil {
		fmt.Fprintln(os.Stderr, "Error running server: %w", err)
		os.Exit(1)
	}
	fmt.Println("Server stopped")
}
