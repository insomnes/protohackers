package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/insomnes/protohackers/pkg/config"
)

type HandlerFunc func(msg []byte, verbose bool) []byte

type Server struct {
	config.ServerConfig
	handlerFunc HandlerFunc
	addr        string
}

func NewServer(
	cfg config.ServerConfig,
	handlerFunc HandlerFunc,
) Server {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	return Server{
		ServerConfig: cfg,
		handlerFunc:  handlerFunc,
		addr:         addr,
	}
}

func (s *Server) Run() error {
	serverAddr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		err := fmt.Errorf("failed to listen: %w", err)
		return err
	}

	fmt.Println("Running server on", serverAddr)
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			err := fmt.Errorf("failed to accept: %w", err)
			return err
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	fmt.Println("Handling connection from", conn.RemoteAddr())
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		conn.SetReadDeadline(time.Now().Add(s.ReadTimeout))
		msg, err := reader.ReadBytes('\n')

		if err != nil && err.Error() == "EOF" {
			fmt.Println("Connection closed by client")
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading from conn:", err)
			return
		}
		resp := s.handlerFunc(msg, s.Verbose)
		// Nothing to do case
		if resp == nil {
			return
		}

		conn.SetWriteDeadline(time.Now().Add(s.WriteTimeout))
		if _, err := conn.Write(resp); err != nil {
			fmt.Fprintln(os.Stderr, "Error writing to conn:", err)
			return
		}
	}
}
