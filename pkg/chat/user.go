package chat

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
)

type UserError struct {
	Addr string
	Name string
	Err  error
}

func (ue UserError) Error() string {
	return fmt.Sprintf("user [%s] (%s) error: %v", ue.Name, ue.Addr, ue.Err)
}

type User struct {
	Address string
	Name    string

	rxChan chan Message
	txChan chan string

	conn net.Conn
}

func NewUser(conn net.Conn, name string) User {
	return User{
		Address: conn.RemoteAddr().String(),
		Name:    name,
		rxChan:  make(chan Message, EventChannelSize),
		txChan:  make(chan string, EventChannelSize),
		conn:    conn,
	}
}

func (u *User) Run(ctx context.Context, messages chan<- Message, fail chan<- UserError) {
	userCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	defer u.conn.Close()

	userFail := make(chan UserError, EventChannelSize)
	go u.runRX(userFail, messages)
	go u.runTX(userCtx, userFail)

	select {
	case <-ctx.Done():
		log.Printf("External stop for [%s] (%s)\n", u.Name, u.Address)
		return
	case err := <-userFail:
		log.Println("Stopping", err.Error())
		fail <- err
		return
	}
}

func (u *User) Send(text string) {
	u.txChan <- text
}

func (u *User) runRX(fail chan<- UserError, messages chan<- Message) {
	reader := bufio.NewReader(u.conn)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("User can not read from %s: %v\n", u.Address, err)
			fail <- u.NewError(err)
			return
		}
		messages <- Message{
			From: u.Name,
			Text: text[:len(text)-1],
		}
	}
}

func (u *User) runTX(ctx context.Context, fail chan<- UserError) {
	for {
		select {
		case <-ctx.Done():
			return
		case text := <-u.txChan:
			_, err := u.conn.Write([]byte(text + "\n"))
			if err != nil {
				log.Printf("Error writing to %s: %v\n", u.Address, err)
				fail <- u.NewError(err)
				return
			}
		}
	}
}

func (u *User) NewError(err error) UserError {
	return UserError{
		Addr: u.Address,
		Name: u.Name,
		Err:  err,
	}
}
