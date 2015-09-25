package db

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/libgo/hc"

	"github.com/megamsys/gulp/meta"
)

func init() {
	hc.AddChecker("Riak", healthCheck)
}

func healthCheck() error {
	conn, err := newConn("test")
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

//A global function which helps to avoid passing config of riak everywhere.
func newConn(bkt string) (*db.Storage, error) {
	log.Debug("New bucket: " + bkt)
	r, err := db.NewRiakDB(meta.MC.Riak, bkt)
	//r, err := db.NewRiakDB([]string{"localhost:8087"}, bkt)
	if err != nil {
		return nil, err
	}

	return r.Conn()
}

func Fetch(bkt string, key string, data interface{}) error {
	s, err := newConn(bkt)
	if err != nil {
		return err
	}
	defer s.Close()
	if err = s.FetchStruct(key, data); err != nil {
		return err
	}
	return nil
}

func FetchObject(bkt string, key string) (string, error) {
	s, err := newConn(bkt)
	if err != nil {
		return "", err
	}
	defer s.Close()
	out := &db.SomeObject{}
	if err = s.FetchObject(key, out); err != nil {
		return "", err
	}
	return out.Data, nil
}

func Store(bkt string, key string, data interface{}) error {
	s, err := newConn(bkt)
	if err != nil {
		return err
	}
	defer s.Close()
	if err = s.StoreStruct(key, data); err != nil {
		return err
	}
	return nil
}
