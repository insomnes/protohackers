package handlers

import (
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
	return NewMeansMsgHandler(verbose, conn.RemoteAddr().String())
}

type MeansMsgHandler struct {
	verbose bool
	db      *BST
	remote  string
}

func NewMeansMsgHandler(verbose bool, remote string) *MeansMsgHandler {
	return &MeansMsgHandler{
		verbose: verbose,
		db:      &BST{root: nil, valCnt: 0},
		remote:  remote,
	}
}

func (mh *MeansMsgHandler) HandleMessage(msg []byte) ([]byte, error) {
	if mh.verbose {
		fmt.Printf("%s->Parsing message: %v\n", mh.remote, msg)
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
		if err.Error() == "invalid range" {
			return []byte{0, 0, 0, 0}, nil
		}
		return nil, err
	}
	fmt.Println(mh.remote, "->Querying", query.from, query.to)
	all := mh.db.Search(query.from, query.to)
	if len(all) == 0 {
		return []byte{0, 0, 0, 0}, nil
	}
	sum := 0
	for _, v := range all {
		sum += int(v)
	}
	if mh.verbose {
		fmt.Println(mh.remote, "->Sum:", sum)
	}
	mean := sum / len(all)
	if mh.verbose {
		fmt.Println("Mean:", mean)
	}
	buf := make([]byte, 4)
	buf[0] = byte(mean >> 24)
	buf[1] = byte(mean >> 16)
	buf[2] = byte(mean >> 8)
	buf[3] = byte(mean)

	return buf, nil
}

func (mh *MeansMsgHandler) handleInsert(msg []byte) ([]byte, error) {
	insert, err := parseInsert(msg)
	if err != nil {
		return nil, err
	}
	fmt.Println(mh.remote, "->Inserting", insert.ts, insert.value)
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
	root   *TreeNode
	valCnt uint32
}

func (b *BST) Insert(qts int32, val int32) {
	b.valCnt += 1
	valNode := &TreeNode{qts: qts, val: val}
	current := b.root
	if current == nil {
		b.root = valNode
		return
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
}

func (b *BST) Search(from, to int32) []int32 {
	result := make([]int32, 0)
	if b.root == nil {
		return make([]int32, 0)
	}
	toVisit := make([]*TreeNode, 0, b.valCnt)
	toVisit = append(toVisit, b.root)

	for len(toVisit) > 0 {
		current := toVisit[0]

		toVisit = toVisit[1:]
		if current.qts >= from && current.qts <= to {
			result = append(result, current.val)
		}
		if current.left != nil && current.qts > from {
			toVisit = append(toVisit, current.left)
		}
		if current.right != nil && current.qts < to {
			toVisit = append(toVisit, current.right)
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
