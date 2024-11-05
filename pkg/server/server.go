package server

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/insomnes/protohackers/pkg/config"
)

type Handler interface {
	HandleMessage(msg []byte, verbose bool, remote string) ([]byte, error)
}

type Server struct {
	config.ServerConfig
	addr       string
	handler    Handler
	readerType ReaderType
}

func NewServer(
	cfg config.ServerConfig,
	handler Handler,
	readerType ReaderType,
) Server {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	return Server{
		ServerConfig: cfg,
		addr:         addr,
		handler:      handler,
		readerType:   readerType,
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
	remote := conn.RemoteAddr().String()
	fmt.Println("Handling connection from", remote)
	defer conn.Close()
	reader := NewMsgReader(conn, s.readerType)
	for {
		conn.SetReadDeadline(time.Now().Add(s.ReadTimeout))
		msg, err := reader.ReadMessage()
		if err != nil {
			switch err.Error() {
			case "EOF":
				fmt.Println("Connection closed by client")
			default:
				fmt.Fprintln(os.Stderr, "Error reading from conn:", err)
			}
			return
		}
		resp, err := s.handler.HandleMessage(msg, s.Verbose, remote)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error handling message:", err)
			return
		}

		conn.SetWriteDeadline(time.Now().Add(s.WriteTimeout))
		if _, err := conn.Write(resp); err != nil {
			fmt.Fprintln(os.Stderr, "Error writing to conn:", err)
			return
		}
	}
}
