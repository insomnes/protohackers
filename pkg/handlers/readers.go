package handlers

import (
	"bufio"
	"io"
	"net"
)

type FullReader struct {
	conn net.Conn
}

func NewFullReader(conn net.Conn) FullReader {
	return FullReader{conn: conn}
}

func (f *FullReader) ReadMessage() ([]byte, error) {
	buf := make([]byte, 4096)
	n, err := f.conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

type LineReader struct {
	bufr *bufio.Reader
}

func NewLineReader(conn net.Conn) LineReader {
	return LineReader{bufr: bufio.NewReader(conn)}
}

func (l *LineReader) ReadMessage() ([]byte, error) {
	return l.bufr.ReadBytes('\n')
}

type NBytesReader struct {
	conn net.Conn
	buf  []byte
	n    int
}

func NewNBytesReader(conn net.Conn, n int) NBytesReader {
	buf := make([]byte, n)
	return NBytesReader{buf: buf, n: n}
}

func (n *NBytesReader) ReadMessage() ([]byte, error) {
	if _, err := io.ReadFull(n.conn, n.buf); err != nil {
		return nil, err
	}
	return n.buf, nil
}
