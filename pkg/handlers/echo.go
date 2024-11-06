package handlers

import (
	"fmt"
	"net"

	"github.com/insomnes/protohackers/pkg/server"
)

type EchoHandler struct{}

func (e *EchoHandler) GetReader(conn net.Conn) server.MsgReader {
	reader := NewFullReader(conn)
	return &reader
}

func (e *EchoHandler) GetMsgHandler(conn net.Conn, verbose bool) server.MsgHandler {
	return &EchoMsgHandler{verbose: verbose}
}

type EchoMsgHandler struct {
	verbose bool
}

func (em *EchoMsgHandler) HandleMessage(msg []byte) ([]byte, error) {
	if em.verbose {
		fmt.Print("Echoing message: ", string(msg))
	}
	return msg, nil
}
