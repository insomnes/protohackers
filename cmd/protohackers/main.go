package main

import (
	"fmt"
	"os"

	"github.com/insomnes/protohackers/pkg/config"
	"github.com/insomnes/protohackers/pkg/handlers"
	"github.com/insomnes/protohackers/pkg/server"
)

var handlerMap = map[string]handlers.Handler{
	"echo":  &handlers.EchoHandler{},
	"prime": &handlers.PrimeHandler{},
}

var readerMap = map[string]server.ReaderType{
	"echo":  server.ReaderTypeBuff,
	"prime": server.ReaderTypeLine,
	"means": server.ReaderTypeNineBytes,
}

func main() {
	cfg := config.ParseConfig()
	handler, ok := handlerMap[cfg.Handler]
	if !ok {
		fmt.Fprintln(os.Stderr, "Unknown handler:", cfg.Handler)
		os.Exit(1)
	}
	readerType, ok := readerMap[cfg.Handler]
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
