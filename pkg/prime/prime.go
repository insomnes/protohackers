package prime

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
	"strings"
	"time"
)

type Input struct {
	Method string `json:"method"`
	Number int    `json:"number"`
}

func isFloatError(err error) bool {
	if typeErr, ok := err.(*json.UnmarshalTypeError); ok &&
		strings.Contains(typeErr.Value, "number") &&
		strings.Contains(typeErr.Value, ".") {
		return true
	}
	return false
}

func parseInput(buffer []byte) (Input, error) {
	var input Input
	err := json.Unmarshal(buffer, &input)
	if err != nil {
		if isFloatError(err) {
			return input, fmt.Errorf("float")
		}
		return input, fmt.Errorf("invalid")
	}
	if input.Method != "isPrime" {
		return input, fmt.Errorf("not-prime")
	}
	return input, nil
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	for i := 5; i <= int(math.Sqrt(float64(n))); i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

const ConnTO time.Duration = time.Second * 30

const (
	falseResponse string = `{"method":"isPrime","prime":false}`
	trueResponse  string = `{"method":"isPrime","prime":true}`
)

func PrimeHandler(conn net.Conn) {
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
			fmt.Fprintln(os.Stderr, "Error reading from conn: %w", err)
			return
		}
		conn.SetWriteDeadline(time.Now().Add(ConnTO))

		input, err := parseInput(buffer[:n])
		if err != nil {
			switch err.Error() {
			case "invalid":
				fmt.Fprintln(os.Stderr, "Invalid json")
				conn.Write([]byte("{}"))
			case "not-prime":
				fmt.Fprintln(os.Stderr, "Invalid method")
				conn.Write([]byte("{}"))
			case "float":
				fmt.Fprintln(os.Stderr, "Invalid number")
				conn.Write([]byte(falseResponse))
			}
			continue
		}
		if isPrime(input.Number) {
			conn.Write([]byte(trueResponse))
		} else {
			conn.Write([]byte(falseResponse))
		}
	}
}
