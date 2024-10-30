package echo

import (
	"fmt"
	"net"
	"os"
	"time"
)

const ConnTO time.Duration = time.Second * 5

func EchoHandler(conn net.Conn) {
	fmt.Println("Handling connection from", conn.RemoteAddr())
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		conn.SetReadDeadline(time.Now().Add(ConnTO))
		n, err := conn.Read(buffer)

		if err != nil && err.Error() == "EOF" {
			fmt.Println("Connection closed by client")
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading from conn:", err)
			return
		}
		conn.SetWriteDeadline(time.Now().Add(ConnTO))
		conn.Write(buffer[:n])
	}
}
