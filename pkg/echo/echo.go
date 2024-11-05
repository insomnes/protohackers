package echo

import (
	"fmt"
)

func EchoHandler(msg []byte, verbose bool) []byte {
	if verbose {
		fmt.Print("Echoing message: ", string(msg))
	}
	return msg
}
