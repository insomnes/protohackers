package server

import (
	"bufio"
	"net"
)

type ReaderType string

const (
	ReaderTypeBuff = ReaderType("buff")
	ReaderTypeLine = ReaderType("line")
)

type MsgReader struct {
	rtype  ReaderType
	reader *bufio.Reader
}

func NewMsgReader(conn net.Conn, rtype ReaderType) MsgReader {
	return MsgReader{
		rtype:  rtype,
		reader: bufio.NewReader(conn),
	}
}

func (r *MsgReader) ReadMessage() ([]byte, error) {
	switch r.rtype {
	case ReaderTypeBuff:
		return readBuff(r.reader)
	case ReaderTypeLine:
		return readLine(r.reader)
	default:
		panic("unknown reader type")
	}
}

func readLine(r *bufio.Reader) ([]byte, error) {
	return r.ReadBytes('\n')
}

func readBuff(r *bufio.Reader) ([]byte, error) {
	msg := make([]byte, 4096)
	_, err := r.Read(msg)
	return msg, err
}
