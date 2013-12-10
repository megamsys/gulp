package db

import (
	"github.com/globocom/config"
	"launchpad.net/gocheck"
	//	"reflect"
	//	"sync"
	"testing"
	//	"time"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct{}

var _ = gocheck.Suite(&S{})

var addr = []string{"127.0.0.1:8098"}

const bkt = "nodes"

/*
func (s *S) SetUpSuite(c *gocheck.C) {
	ticker.Stop()
}

func (s *S) TearDownSuite(c *gocheck.C) {
	storage, err := Open(addr, bkt)
	c.Assert(err, gocheck.IsNil)
	defer storage.coder_client.Close()
}

func (s *S) TearDownTest(c *gocheck.C) {
	conn = make(map[string]*session)
}

func (s *S) TestOpenConnectsToTheDatabase(c *gocheck.C) {
	storage, err := Open(addr, bkt)
	c.Assert(err, gocheck.IsNil)
	defer storage.Close()
	_, err = storage.coder_client.Ping()
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestOpenReconnects(c *gocheck.C) {
	storage, err := Open(addr, bkt)
	c.Assert(err, gocheck.IsNil)
	storage.Close()
	storage, err = Open(addr, bkt)
	c.Assert(err, gocheck.IsNil)
	_, err = storage.coder_client.Ping()
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestOpenConnectionRefused(c *gocheck.C) {
	storage, err := Open([]string{"127.0.0.1:68098"}, bkt)
	c.Assert(storage, gocheck.IsNil)
	c.Assert(err, gocheck.NotNil)
}

func (s *S) TestClose(c *gocheck.C) {
	defer func() {
		r := recover()
		c.Check(r, gocheck.IsNil)
	}()

	storage, err := Open(addr, bkt)
	storage.Close()
	c.Assert(err, gocheck.IsNil)
	c.Assert(storage, gocheck.NotNil)
	_, err = storage.coder_client.Ping()
	c.Check(err, gocheck.NotNil)
}
*/
func (s *S) TestConn(c *gocheck.C) {
	config.Set("riak:url", "127.0.0.1:8087")
	defer config.Unset("riak:url")
	config.Set("bucket:name", "nodes")
	defer config.Unset("bucket:name")
	storage, err := Conn()
	c.Assert(err, gocheck.IsNil)
	c.Assert(storage, gocheck.NotNil)
	_, err = storage.coder_client.Ping()
	c.Check(err, gocheck.IsNil)
}

/*func (s *S) TestUsers(c *gocheck.C) {
	storage, _ := Open("127.0.0.1:27017", "megam_storage_test")
	defer storage.Close()
	users := storage.Users()
	usersc := storage.Collection("users")
	c.Assert(users, gocheck.DeepEquals, usersc)
	c.Assert(users, HasUniqueIndex, []string{"email"})
}
*/

/*
func (s *S) TestRetire(c *gocheck.C) {
	defer func() {
		if r := recover(); !c.Failed() && r == nil {
			c.Errorf("Should panic in ping, but did not!")
		}
	}()
	Open("127.0.0.1:8098", "megam_storage_test")
	sess := conn["127.0.0.1:8098"]
	sess.used = sess.used.Add(-1 * 2 * period)
	conn["127.0.0.1:8098"] = sess
	var ticker time.Ticker
	ch := make(chan time.Time, 1)
	ticker.C = ch
	ch <- time.Now()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		retire(&ticker)
		wg.Done()
	}()
	close(ch)
	wg.Wait()
	_, ok := conn["127.0.0.1:8098"]
	c.Check(ok, gocheck.Equals, false)
	sess.s.Ping()
}
*/
