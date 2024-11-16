package mitm

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type ConnError struct {
	Addr string
	Err  error
}

func (ce ConnError) Error() string {
	return fmt.Sprintf("conn <%s> error: %v", ce.Addr, ce.Err)
}

type MitmConn struct {
	Address string
	conn    net.Conn

	tx chan string
}

func NewMitmConn(conn net.Conn, addr string) MitmConn {
	return MitmConn{
		Address: addr,
		conn:    conn,
		tx:      make(chan string, EventChannelSize),
	}
}

func (mc *MitmConn) Run(ctx context.Context, up chan<- string, fail chan<- ConnError) {
	go mc.runRX(up, fail)
	go mc.runTX(ctx, fail)
}

func (mc *MitmConn) QueueSend(message string) {
	mc.tx <- message
}

func (mc *MitmConn) runRX(up chan<- string, fail chan<- ConnError) {
	defer mc.conn.Close()
	reader := bufio.NewReader(mc.conn)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			if errors.Is(err, io.EOF) {
				log.Printf("Conn closed by %s\n", mc.Address)
			} else {
				log.Printf("Conn can not read from %s: %v\n", mc.Address, err)
			}
			fail <- mc.NewError(err)
			return
		}
		up <- text
	}
}

func (mc *MitmConn) runTX(ctx context.Context, fail chan<- ConnError) {
	defer mc.conn.Close()
	for {
		select {
		case <-ctx.Done():
			return
		case text := <-mc.tx:
			_, err := mc.conn.Write([]byte(text))
			if err != nil {
				log.Printf("Conn can not write to %s: %v\n", mc.Address, err)
				fail <- mc.NewError(err)
				return
			}
		}
	}
}

func (mc *MitmConn) NewError(err error) ConnError {
	return ConnError{
		Addr: mc.Address,
		Err:  err,
	}
}
