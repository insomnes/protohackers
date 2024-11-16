package mitm

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"regexp"
)

const (
	chatAddr      string = "chat.protohackers.com:16963"
	tonyWallet    string = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"
	walletPattern string = `(?:^|\s)(7[a-zA-Z0-9]{25,34})(?:$|[\s\n])`
)

var re *regexp.Regexp = regexp.MustCompile(walletPattern)

func tonyWalletFix(text string) string {
	return re.ReplaceAllString(text, tonyWallet)
}

func RunMitmProxy(ctx context.Context, conn net.Conn) {
	userCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	defer conn.Close()
	defer log.Printf("Mitm proxy for %s closed", conn.LocalAddr().String())

	userConn := NewMitmConn(conn, conn.LocalAddr().String())
	userUp := make(chan string, EventChannelSize)

	chatConn, err := createChatServerConn()
	if err != nil {
		log.Printf("Failed to create chat connection for %s: %v", userConn.Address, err)
		return
	}
	chatUp := make(chan string, EventChannelSize)

	fail := make(chan ConnError, EventChannelSize)

	userConn.Run(userCtx, userUp, fail)
	chatConn.Run(userCtx, chatUp, fail)

	log.Printf("Mitm proxy for %s<->%s started", userConn.Address, chatConn.Address)

	for {
		select {
		case <-ctx.Done():
			return
		case userMessage := <-userUp:
			chatConn.QueueSend(tonyWalletFix(userMessage))
		case chatMessage := <-chatUp:
			userConn.QueueSend(tonyWalletFix(chatMessage))
		case err := <-fail:
			if errors.Is(err.Err, net.ErrClosed) || errors.Is(err.Err, io.EOF) {
				return
			}
			log.Printf(
				"Err on conn for pair %s<->%s => %v",
				userConn.Address,
				chatConn.Address,
				err,
			)
			return
		}
	}
}

func createChatServerConn() (MitmConn, error) {
	srvTCPAddr, err := net.ResolveTCPAddr("tcp", chatAddr)
	if err != nil {
		return MitmConn{}, err
	}

	conn, err := net.DialTCP("tcp", nil, srvTCPAddr)
	if err != nil {
		return MitmConn{}, err
	}

	return NewMitmConn(conn, conn.LocalAddr().String()), nil
}
