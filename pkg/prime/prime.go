package prime

import (
	"encoding/json"
	"fmt"
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

var (
	malformedResponse []byte = []byte("{}\n")
	falseResponse     []byte = []byte(`{"method":"isPrime","prime":false}` + "\n")
	trueResponse      []byte = []byte(`{"method":"isPrime","prime":true}` + "\n")
)

func PrimeHandler(msg []byte) []byte {
	if len(msg) > 1 {
		fmt.Println("Prime message:", string(msg[:len(msg)-1]))
	}

	input, err := parseInput(msg)
	if err != nil {
		switch err.Error() {
		case "invalid":
			fmt.Fprintln(os.Stderr, "Invalid json")
			return malformedResponse
		case "not-prime":
			fmt.Fprintln(os.Stderr, "Invalid method")
			return malformedResponse
		case "float":
			fmt.Fprintln(os.Stderr, "Invalid number")
			return falseResponse
		default:
			fmt.Fprintln(os.Stderr, "Other error", err)
			return malformedResponse
		}
	}
	inputIsPrime := isPrime(*input.Number)

	if !inputIsPrime {
		fmt.Println("Got non-prime number:", *input.Number)
		return falseResponse
	}

	fmt.Println("Got prime number:", *input.Number)
	return trueResponse
}
