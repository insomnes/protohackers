package handlers

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/insomnes/protohackers/pkg/server"
)

type MeansHandler struct{}

func (mh *MeansHandler) GetReader(conn net.Conn) server.MsgReader {
	reader := NewNBytesReader(conn, 9)
	return &reader
}

func (mh *MeansHandler) GetMsgHandler(conn net.Conn, verbose bool) server.MsgHandler {
	return NewMeansMsgHandler(verbose)
}

type MeansMsgHandler struct {
	verbose bool
	db      *BST
}

func NewMeansMsgHandler(verbose bool) *MeansMsgHandler {
	return &MeansMsgHandler{
		verbose: verbose,
		db:      &BST{root: nil},
	}
}

func (mh *MeansMsgHandler) HandleMessage(msg []byte) ([]byte, error) {
	if mh.verbose {
		fmt.Printf("Parsing message: %v\n", msg)
	}
	msgType := msg[0]
	switch msgType {
	case byte('Q'):
		return mh.handleQuery(msg[1:])
	case byte('I'):
		return mh.handleInsert(msg[1:])
	default:
		return nil, fmt.Errorf("invalid message type")
	}
}

func (mh *MeansMsgHandler) handleQuery(msg []byte) ([]byte, error) {
	query, err := parseQuery(msg)
	if err != nil {
		return nil, err
	}
	all := mh.db.Search(query.from, query.to)
	if len(all) == 0 {
		return []byte{0, 0, 0, 0}, nil
	}
	sum := int32(0)
	for _, v := range all {
		sum += v
	}
	mean := sum / int32(len(all))
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, mean)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (mh *MeansMsgHandler) handleInsert(msg []byte) ([]byte, error) {
	insert, err := parseInsert(msg)
	if err != nil {
		return nil, err
	}
	mh.db.Insert(insert.ts, insert.value)
	return nil, nil
}

type Query struct {
	from int32
	to   int32
}

func parseQuery(in []byte) (Query, error) {
	data, err := parseData(in)
	if err != nil {
		return Query{}, err
	}
	if data[0] > data[1] {
		return Query{}, fmt.Errorf("invalid range")
	}
	return Query{from: data[0], to: data[1]}, nil
}

type Insert struct {
	ts    int32
	value int32
}

func parseInsert(in []byte) (Insert, error) {
	data, err := parseData(in)
	if err != nil {
		return Insert{}, err
	}
	return Insert{ts: data[0], value: data[1]}, nil
}

func parseBigEndian(in []byte) (int32, error) {
	if len(in) != 4 {
		return 0, fmt.Errorf("invalid length for int32 big endian")
	}
	return int32(in[0])<<24 | int32(in[1])<<16 | int32(in[2])<<8 | int32(in[3]), nil
}

func parseData(in []byte) ([2]int32, error) {
	var data [2]int32
	if len(in) != 8 {
		return data, fmt.Errorf("invalid length")
	}
	from, err := parseBigEndian(in[0:4])
	if err != nil {
		return data, err
	}
	to, err := parseBigEndian(in[4:8])
	if err != nil {
		return data, err
	}
	data[0] = from
	data[1] = to
	return data, nil
}

type BST struct {
	root *TreeNode
}

func (b *BST) Insert(qts int32, val int32) error {
	valNode := &TreeNode{qts: qts, val: val}
	current := b.root
	if current == nil {
		b.root = valNode
		return nil
	}
	for {
		if qts < current.qts {
			if current.left == nil {
				current.left = valNode
				break
			}
			current = current.left
		} else {
			if current.right == nil {
				current.right = valNode
				break
			}
			current = current.right
		}
	}
	return nil
}

func (b *BST) Search(from, to int32) []int32 {
	current := b.root
	result := make([]int32, 0)
	for current != nil {
		if current.qts >= from && current.qts <= to {
			result = append(result, current.val)
		}
		if current.qts > from {
			current = current.left
		} else {
			current = current.right
		}
	}
	return result
}

type TreeNode struct {
	qts   int32
	val   int32
	left  *TreeNode
	right *TreeNode
}