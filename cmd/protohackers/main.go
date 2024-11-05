package main

import (
	"fmt"
	"os"

	"github.com/insomnes/protohackers/pkg/config"
	"github.com/insomnes/protohackers/pkg/echo"
	"github.com/insomnes/protohackers/pkg/prime"
	"github.com/insomnes/protohackers/pkg/server"
)

func GetHandlerFunction(handle string) (server.HandlerFunc, error) {
	switch handle {
	case "echo":
		return echo.EchoHandler, nil
	case "prime":
		return prime.PrimeHandler, nil
	default:
		return nil, fmt.Errorf("unknown handler: %s", handle)
	}
}

func main() {
	cfg := config.ParseConfig()
	handlerFunc, err := GetHandlerFunction(cfg.Handler)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting handler:", err)
		os.Exit(1)
	}

	server := server.NewServer(cfg, handlerFunc)

	if err := server.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error running server:", err)
		os.Exit(1)
	}
	fmt.Println("Server stopped")
}
