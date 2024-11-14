package kvstore

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
)

type (
	Key   string
	Value string
)

const (
	ServerVersion = Value("KKVS 1.0")
	versionKey    = "version"
)

type QueryType int

const (
	Retrieve QueryType = iota
	Insert
	VersionReq
)

func (qt QueryType) String() string {
	return [...]string{"Retrieve", "Insert", "VersionReq"}[qt]
}

type Query struct {
	Type QueryType
	Key  Key
	Val  Value
	From net.UDPAddr
}

func (q Query) String() string {
	return fmt.Sprintf("Query: %s, %s, %s, %s", q.Type, q.Key, q.Val, q.From.String())
}

func NewQuery(qType QueryType, key string, val string, from net.UDPAddr) Query {
	return Query{
		Type: qType,
		Key:  Key(key),
		Val:  Value(val),
		From: from,
	}
}

func QueryFromBytes(b []byte, from net.UDPAddr) Query {
	var qType QueryType

	key, offset := extractKey(b)
	if offset != -1 {
		return NewQuery(Insert, key, string(b[offset:]), from)
	}

	if key == versionKey {
		qType = VersionReq
	} else {
		qType = Retrieve
	}
	return NewQuery(qType, key, "", from)
}

// Returning key and '=' offset (i + 1 to skip it) or -1 if not found
func extractKey(b []byte) (string, int) {
	offset := -1
	builder := strings.Builder{}
	for i, c := range b {
		if c == '=' {
			offset = i + 1
			break
		}
		builder.WriteByte(c)
	}
	return builder.String(), offset
}

type Response struct {
	Key Key
	Val Value
	To  net.UDPAddr
}

func (r Response) Bytes() []byte {
	return []byte(fmt.Sprintf("%s=%s", r.Key, r.Val))
}

func (r Response) String() string {
	return fmt.Sprintf("Response: %s, %s, %s", r.Key, r.Val, r.To.String())
}

type DB struct {
	storage map[Key]Value

	queries chan Query
}

func NewDB() *DB {
	return &DB{
		storage: make(map[Key]Value),
		queries: make(chan Query, EventChannelSize),
	}
}

func (db *DB) Run(ctx context.Context, results chan Response) {
	log.Println("Running DB")
	for {
		select {
		case <-ctx.Done():
			log.Println("DB shutting down")
			return
		case q := <-db.queries:
			db.handleQuery(q, results)
		}
	}
}

func (db *DB) QueueQuery(q Query) {
	if q.Type == Insert && q.Key == versionKey {
		return
	}
	db.queries <- q
}

func (db *DB) handleQuery(q Query, results chan Response) {
	switch q.Type {
	case Retrieve:
		val := db.retrieve(q.Key)
		results <- Response{Key: q.Key, Val: val, To: q.From}
	case Insert:
		if q.Key == versionKey {
			return
		}
		db.insert(q.Key, q.Val)
	case VersionReq:
		val := db.version()
		results <- Response{Key: q.Key, Val: val, To: q.From}
	}
}

func (db *DB) insert(key Key, val Value) {
	db.storage[key] = val
}

func (db *DB) retrieve(key Key) Value {
	return db.storage[key]
}

func (db *DB) version() Value {
	return ServerVersion
}
