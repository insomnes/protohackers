package main

import (
	"fmt"
	"os"

	"github.com/insomnes/protohackers/pkg/config"
	"github.com/insomnes/protohackers/pkg/echo"
	"github.com/insomnes/protohackers/pkg/prime"
	"github.com/insomnes/protohackers/pkg/server"
)

var handlers = map[string]server.Handler{
	"echo":  &echo.EchoHandler{},
	"prime": &prime.PrimeHandler{},
}

var readerTypes = map[string]server.ReaderType{
	"echo":  server.ReaderTypeBuff,
	"prime": server.ReaderTypeLine,
}

func main() {
	cfg := config.ParseConfig()
	handler, ok := handlers[cfg.Handler]
	if !ok {
		fmt.Fprintln(os.Stderr, "Unknown handler:", cfg.Handler)
		os.Exit(1)
	}
	readerType, ok := readerTypes[cfg.Handler]
	if !ok {
		fmt.Fprintln(os.Stderr, "Unknown reader type for handler:", cfg.Handler)
		os.Exit(1)
	}

	server := server.NewServer(cfg, handler, readerType)

	if err := server.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error running server:", err)
		os.Exit(1)
	}
	fmt.Println("Server stopped")
}
