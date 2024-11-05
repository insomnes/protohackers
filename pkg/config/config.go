package config

import (
	"flag"
	"log"
	"time"
)

const (
	DefaultHost         string = "localhost"
	DefaultPort         uint16 = 9999
	DefaultHandler      string = "echo"
	DefaultReadTimeout  uint   = 20
	DefaultWriteTimeout uint   = 20
)

type ServerConfig struct {
	Host         string
	Port         uint16
	Handler      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func ParseConfig() ServerConfig {
	host := flag.String("host", DefaultHost, "address to listen on")
	port := flag.Uint("port", uint(DefaultPort), "port to listen on 1-65535")
	handler := flag.String("handler", DefaultHandler, "handler to use")
	readTimeout := flag.Uint("read-timeout", DefaultReadTimeout, "read timeout")
	writeTimeout := flag.Uint("write-timeout", DefaultWriteTimeout, "write timeout")
	flag.Parse()

	if *port < 1 || *port > 65535 {
		log.Fatalf("Invalid port number: %d", *port)
	}
	if *readTimeout < 1 {
		log.Fatalf("Invalid read timeout: %d", *readTimeout)
	}
	if *writeTimeout < 1 {
		log.Fatalf("Invalid write timeout: %d", *writeTimeout)
	}

	return ServerConfig{
		Host:         *host,
		Port:         uint16(*port),
		Handler:      *handler,
		ReadTimeout:  time.Duration(*readTimeout) * time.Second,
		WriteTimeout: time.Duration(*writeTimeout) * time.Second,
	}
}
