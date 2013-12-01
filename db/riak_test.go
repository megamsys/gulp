
package db

import (
//	"github.com/globocom/config"
	"launchpad.net/gocheck"
//	"reflect"
//	"sync"
	"testing"
//	"time"
"log"
)


func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct{}

var _ = gocheck.Suite(&S{})

func (s *S) SetUpSuite(c *gocheck.C) {
	ticker.Stop()
}

/*func (s *S) TearDownSuite(c *gocheck.C) {
	storage, err := Open("127.0.0.1:8098", "megam_storage_test")
	c.Assert(err, gocheck.IsNil)
	defer storage.coder_client.Close()
//	storage.session.DB("megam_storage_test").DropDatabase()
}

func (s *S) TearDownTest(c *gocheck.C) {
	conn = make(map[string]*session)
}

func (s *S) TestOpenConnectsToTheDatabase(c *gocheck.C) {
	storage, err := Open("127.0.0.1:8098", "megam_storage_test")
	c.Assert(err, gocheck.IsNil)
	defer storage.Close()
	_, err = storage.coder_client.Ping()
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestOpenCopiesConnection(c *gocheck.C) {
	storage, err := Open("127.0.0.1:8098", "megam_storage_test")
	c.Assert(err, gocheck.IsNil)
	defer storage.Close()
	storage2, err := Open("127.0.0.1:8098", "megam_storage_test")
	c.Assert(err, gocheck.IsNil)
	c.Assert(storage.coder_client, gocheck.Not(gocheck.Equals), storage2.coder_client)
}

func (s *S) TestOpenReconnects(c *gocheck.C) {
	storage, err := Open("127.0.0.1:8098", "megam_storage_test")
	c.Assert(err, gocheck.IsNil)
	storage.Close()
	storage, err = Open("127.0.0.1:8098", "megam_storage_test")
	c.Assert(err, gocheck.IsNil)
	_, err = storage.coder_client.Ping()
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestOpenConnectionRefused(c *gocheck.C) {
	storage, err := Open("127.0.0.1:27018", "megam_storage_test")
	c.Assert(storage, gocheck.IsNil)
	c.Assert(err, gocheck.NotNil)
}
*/
func (s *S) TestClose(c *gocheck.C) {
	defer func() {
		r := recover()
		c.Check(r, gocheck.IsNil)
	}()
	storage, err := Open("127.0.0.1:8098", "megam_storage_test")
	c.Assert(err, gocheck.IsNil)
	storage.Close()
	_, err = storage.coder_client.Ping()
	log.Printf(" --> %s", err)
	c.Check(err, gocheck.NotNil)
}
/*
func (s *S) TestConn(c *gocheck.C) {
	config.Set("riak:url", "127.0.0.1:8098")
	defer config.Unset("riak:url")
	config.Set("bucket:name", "megam_storage_test")
	defer config.Unset("bucket:name")
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
