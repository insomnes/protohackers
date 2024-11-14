package kvstore

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const EventChannelSize = 16

type KVServer struct {
	Address net.UDPAddr
	db      *DB
}

func NewKVServer(address string) KVServer {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal("Error resolving address: ", err)
	}
	db := NewDB()
	return KVServer{
		Address: *addr,
		db:      db,
	}
}

func (cs *KVServer) Run() {
	conn, err := net.ListenUDP("udp", &cs.Address)
	if err != nil {
		log.Fatal("Error listening: ", err)
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())

	results := make(chan Response, EventChannelSize)
	defer cancel()
	go cs.db.Run(ctx, results)
	go cs.readData(conn)
	go cs.processResults(ctx, conn, results)

	sig := <-sigChan
	log.Printf("Signal received: %v\n", sig)
	cancel()
	<-time.After(300 * time.Millisecond)
}

func (cs *KVServer) readData(conn *net.UDPConn) {
	log.Println("Reading data agent started")
	for {
		data := make([]byte, 1000)
		n, addr, err := conn.ReadFromUDP(data)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Println("Connection closed")
			} else {
				log.Println("Error reading data:", err)
			}
			break
		}

		go func() {
			cs.db.QueueQuery(QueryFromBytes(data[:n], *addr))
		}()
	}
}

func (cs *KVServer) processResults(ctx context.Context, conn *net.UDPConn, results chan Response) {
	log.Println("Results agent started")
	for {
		select {
		case <-ctx.Done():
			log.Println("Results agent shutting down")
			return
		case res := <-results:
			go sendResponse(conn, res)
		}
	}
}

func sendResponse(conn *net.UDPConn, res Response) {
	_, err := conn.WriteToUDP(res.Bytes(), &res.To)
	if err != nil {
		log.Printf("Error sending response to %s: %v", &res.To, err)
	}
}
