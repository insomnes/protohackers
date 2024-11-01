package prime

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type Input struct {
	Method string `json:"method"`
	Number *uint  `json:"number"`
}

func isFloatError(err error) bool {
	if typeErr, ok := err.(*json.UnmarshalTypeError); ok &&
		strings.Contains(typeErr.Value, ".") {
		return true
	}
	return false
}

func isNumberError(err error) bool {
	if typeErr, ok := err.(*json.UnmarshalTypeError); ok &&
		strings.Contains(typeErr.Value, "number") {
		return true
	}
	return false
}

func parseInput(buffer []byte) (Input, error) {
	var input Input
	err := json.Unmarshal(buffer, &input)
	if err != nil {
		// Big numbers will give 0, and non-prime
		if !isNumberError(err) {
			return input, fmt.Errorf("invalid")
		}
		if isFloatError(err) {
			return input, fmt.Errorf("float")
		}
		fmt.Fprintln(os.Stderr, "Possibly too big number in:", err)
	}
	if input.Method != "isPrime" {
		return input, fmt.Errorf("not-prime")
	}
	if input.Number == nil {
		return input, fmt.Errorf("invalid")
	}
	return input, nil
}

func isPrime(n uint) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}

	var i uint
	for i = 5; i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

const ConnTO time.Duration = time.Second * 30

const (
	malformedResponse string = "{}\n"
	falseResponse     string = `{"method":"isPrime","prime":false}` + "\n"
	trueResponse      string = `{"method":"isPrime","prime":true}` + "\n"
)

func PrimeHandler(conn net.Conn) {
	fmt.Println("Handling connection from", conn.RemoteAddr())
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(ConnTO))
		msg, err := reader.ReadBytes(byte('\n'))

		if len(msg) > 1 {
			fmt.Println("Got message:", string(msg[:len(msg)-1]))
		}

		if err != nil && err.Error() == "EOF" {
			fmt.Println("Connection closed by client")
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading from conn: ", err)
			return
		}

		input, err := parseInput(msg)
		if err != nil {
			switch err.Error() {
			case "invalid":
				fmt.Fprintln(os.Stderr, "Invalid json")
				conn.Write([]byte(malformedResponse))
			case "not-prime":
				fmt.Fprintln(os.Stderr, "Invalid method")
				conn.Write([]byte(malformedResponse))
			case "float":
				fmt.Fprintln(os.Stderr, "Invalid number")
				conn.Write([]byte(falseResponse))
			default:
				fmt.Fprintln(os.Stderr, "Other error", err)
				conn.Write([]byte(malformedResponse))
			}
			continue
		}
		inputIsPrime := isPrime(*input.Number)

		conn.SetWriteDeadline(time.Now().Add(ConnTO))
		if inputIsPrime {
			fmt.Println("Got prime number:", *input.Number)
			conn.Write([]byte(trueResponse))
		} else {
			fmt.Println("Got non-prime number:", *input.Number)
			conn.Write([]byte(falseResponse))
		}
	}
}
