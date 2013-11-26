// Package db encapsulates connection with Riak.
//
// The function Open dials to Riak and returns a connection (represented by
// the Storage type). It manages an internal pool of connections, and
// reconnects in case of failures. That means that you should not store
// references to the connection, but always call Open.
package db

import (
	"github.com/globocom/config"
	"github.com/mrb/riakpbc"
	"fmt"
	"sync"
	"time"
)

var (
	conn   = make(map[string]*session) // pool of connections
	mut    sync.RWMutex                // for pool thread safety
	ticker *time.Ticker                // for garbage collection
)

const (
	DefaultRiakURL    = "http://localhost:8098/riak/"
	DefaultBucketName = "nodes"
)

const period time.Duration = 7 * 24 * time.Hour

type session struct {
	s    *riakpbc.Client
	used time.Time
}

// Storage holds the connection with the bucket name.
type Storage struct {
	coder_client *riakpbc.Client
	bktname      string
}

func open(addr, bucketname string) (*Storage, error) {
	// Alternative marshallers can be built from this interface.
	coder := riakpbc.NewCoder("json", riakpbc.JsonMarshaller, riakpbc.JsonUnmarshaller)
	riakCoder := riakpbc.NewClientWithCoder([]string{addr}, coder)

	err := riakCoder.Dial()
	if err != nil {
		return nil, err
	}

	storage := &Storage{coder_client: riakCoder, bktname: bucketname}
	mut.Lock()
	conn[addr] = &session{s: riakCoder, used: time.Now()}
	mut.Unlock()
	return storage, nil
}

// Open dials to the Riak database, and return the connection (represented
// by the type Storage).
//
// addr is a MongoDB connection URI, and bktname is the name of the bucket.
//
// This function returns a pointer to a Storage, or a non-nil error in case of
// any failure.
func Open(addr, bktname string) (storage *Storage, err error) {
	defer func() {
		if r := recover(); r != nil {
			storage, err = open(addr, bktname)
		}
	}()
	mut.RLock()
	if session, ok := conn[addr]; ok {
		mut.RUnlock()
		if _, err = session.s.Ping(); err == nil {
			mut.Lock()
			session.used = time.Now()
			conn[addr] = session
			mut.Unlock()
			return &Storage{session.s, bktname}, nil
		}
		return open(addr, bktname)
	}
	mut.RUnlock()
	return open(addr, bktname)
}

// Conn reads the megam config and calls Open to get a database connection.
//
// Most megam packages should probably use this function. Open is intended for
// use when supporting more than one database.
func Conn() (*Storage, error) {
	url, _ := config.GetString("riak:url")
	if url == "" {
		url = DefaultRiakURL
	}
	bktname, _ := config.GetString("bucket:name")
	if bktname == "" {
		bktname = DefaultBucketName
	}
	return Open(url, bktname)
}

// Close closes the storage, releasing the connection.
func (s *Storage) Close() {
	s.coder_client.Close()
}

// FetchStruct stores a struct  as JSON
//   eg: data := ExampleData{
//        Field1: "ExampleData1",
//        Field2: 1,
//   }
// So the send can pass in 	out := &ExampleData{}
// Apps returns the apps collection from MongoDB.
func (s *Storage) FetchStruct(key string, out interface{}) error {
	if _, err := s.coder_client.FetchStruct(s.bktname, key, out); err != nil {
		return fmt.Errorf("Convert fetched JSON to the Struct, and return it failed: %s", err)
	}
	//TO-DO:
	//we need to return the fetched json -> to struct interface
	return nil
}

// StoreStruct returns the apps collection from MongoDB.
func (s *Storage) StoreStruct(key string, data interface{}) error {
	if _, err := s.coder_client.StoreStruct(s.bktname, key, &data); err != nil {
		return fmt.Errorf("Convert fetched JSON to the Struct, and return it failed: %s", err)
	}
	return nil
}

func init() {
	ticker = time.NewTicker(time.Hour)
	go retire(ticker)
}

// retire retires old connections
func retire(t *time.Ticker) {
	for _ = range t.C {
		now := time.Now()
		var old []string
		mut.RLock()
		for k, v := range conn {
			if now.Sub(v.used) >= period {
				old = append(old, k)
			}
		}
		mut.RUnlock()
		mut.Lock()
		for _, c := range old {
			conn[c].s.Close()
			delete(conn, c)
		}
		mut.Unlock()
	}
}
