package echo

import (
	"fmt"
)

func EchoHandler(msg []byte) []byte {
	fmt.Println("Echoing message:", string(msg))
	return msg
}
