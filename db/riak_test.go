
package db

import (
	"github.com/globocom/config"
	"launchpad.net/gocheck"
//	"reflect"
	"sync"
	"testing"
	"time"
)


func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct{}

var _ = gocheck.Suite(&S{})

func (s *S) SetUpSuite(c *gocheck.C) {
	ticker.Stop()
}

func (s *S) TearDownSuite(c *gocheck.C) {
	storage, err := Open("127.0.0.1:27017", "megam_storage_test")
	c.Assert(err, gocheck.IsNil)
	defer storage.coder_client.Close()
//	storage.session.DB("megam_storage_test").DropDatabase()
}

func (s *S) TearDownTest(c *gocheck.C) {
	conn = make(map[string]*session)
}

func (s *S) TestOpenConnectsToTheDatabase(c *gocheck.C) {
	storage, err := Open("http://localhost:8098/riak/", "megam_storage_test")
	c.Assert(err, gocheck.IsNil)
	defer storage.Close()
	_, err = storage.coder_client.Ping()
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestOpenCopiesConnection(c *gocheck.C) {
	storage, err := Open("http://localhost:8098/riak/", "megam_storage_test")
	c.Assert(err, gocheck.IsNil)
	defer storage.Close()
	storage2, err := Open("http://localhost:8098/riak/", "megam_storage_test")
	c.Assert(err, gocheck.IsNil)
	c.Assert(storage.coder_client, gocheck.Not(gocheck.Equals), storage2.coder_client)
}

func (s *S) TestOpenReconnects(c *gocheck.C) {
	storage, err := Open("http://localhost:8098/riak/", "megam_storage_test")
	c.Assert(err, gocheck.IsNil)
	storage.Close()
	storage, err = Open("http://localhost:8098/riak/", "megam_storage_test")
	c.Assert(err, gocheck.IsNil)
	_, err = storage.coder_client.Ping()
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestOpenConnectionRefused(c *gocheck.C) {
	storage, err := Open("127.0.0.1:27018", "megam_storage_test")
	c.Assert(storage, gocheck.IsNil)
	c.Assert(err, gocheck.NotNil)
}

func (s *S) TestClose(c *gocheck.C) {
	defer func() {
		r := recover()
		c.Check(r, gocheck.NotNil)
	}()
	storage, err := Open("http://localhost:8098/riak/", "megam_storage_test")
	c.Assert(err, gocheck.IsNil)
	storage.Close()
	_, err = storage.coder_client.Ping()
	c.Check(err, gocheck.NotNil)
}

func (s *S) TestConn(c *gocheck.C) {
	config.Set("database:url", "http://localhost:8098/riak/")
	defer config.Unset("database:url")
	config.Set("database:name", "megam_storage_test")
	defer config.Unset("database:name")
	storage, err := Conn()
	c.Assert(err, gocheck.IsNil)
	defer storage.Close()
	_,  err = storage.coder_client.Ping()
	c.Assert(err, gocheck.IsNil)
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


func (s *S) TestRetire(c *gocheck.C) {
	defer func() {
		if r := recover(); !c.Failed() && r == nil {
			c.Errorf("Should panic in ping, but did not!")
		}
	}()
	Open("http://localhost:8098/riak/", "megam_storage_test")
	sess := conn["http://localhost:8098/riak/"]
	sess.used = sess.used.Add(-1 * 2 * period)
	conn["http://localhost:8098/riak/"] = sess
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
	_, ok := conn["http://localhost:8098/riak/"]
	c.Check(ok, gocheck.Equals, false)
	sess.s.Ping()
}

