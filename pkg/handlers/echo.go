package handlers

import (
	"fmt"
)

type EchoHandler struct{}

func (e *EchoHandler) HandleMessage(msg []byte, verbose bool, remote string) ([]byte, error) {
	if verbose {
		fmt.Print("Echoing message: ", string(msg))
	}
	return msg, nil
}
