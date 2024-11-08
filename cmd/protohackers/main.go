package main

import (
	"fmt"
	"os"

	"github.com/insomnes/protohackers/pkg/config"
	"github.com/insomnes/protohackers/pkg/handlers"
	"github.com/insomnes/protohackers/pkg/server"
)

var handlerMap = map[string]server.ConnHandler{
	"echo":  &handlers.EchoHandler{},
	"prime": &handlers.PrimeHandler{},
	"means": &handlers.MeansHandler{},
}

func main() {
	cfg := config.ParseConfig()
	handler, ok := handlerMap[cfg.Handler]
	if !ok {
		fmt.Fprintln(os.Stderr, "Unknown handler:", cfg.Handler)
		os.Exit(1)
	}

	server := server.NewServer(cfg, handler)

	if err := server.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error running server:", err)
		os.Exit(1)
	}
	fmt.Println("Server stopped")
}
