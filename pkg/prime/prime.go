package prime

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
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

var (
	malformedResponse []byte = []byte("{}\n")
	falseResponse     []byte = []byte(`{"method":"isPrime","prime":false}` + "\n")
	trueResponse      []byte = []byte(`{"method":"isPrime","prime":true}` + "\n")
)

type PrimeHandler struct{}

func (ph *PrimeHandler) HandleMessage(msg []byte, verbose bool, remote string) ([]byte, error) {
	if len(msg) > 1 && verbose {
		fmt.Print("Prime message: ", string(msg))
	}

	input, err := parseInput(msg)
	if err != nil {
		switch err.Error() {
		case "invalid":
			fmt.Fprintln(os.Stderr, "Invalid json")
			return malformedResponse, nil
		case "not-prime":
			fmt.Fprintln(os.Stderr, "Invalid method")
			return malformedResponse, nil
		case "float":
			fmt.Fprintln(os.Stderr, "Invalid number")
			return falseResponse, nil
		default:
			fmt.Fprintln(os.Stderr, "Other error", err)
			return malformedResponse, nil
		}
	}
	inputIsPrime := isPrime(*input.Number)

	if !inputIsPrime {
		if verbose {
			fmt.Println("Got non-prime number:", *input.Number)
		}
		return falseResponse, nil
	}
	if verbose {
		fmt.Println("Got prime number:", *input.Number)
	}
	return trueResponse, nil
}
