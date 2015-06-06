package schmoopy

import (
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type SchmoopySuite struct {
	server     *schmoopyServer
	dbFilename string
}

var _ = Suite(&SchmoopySuite{})

func (s *SchmoopySuite) SetUpSuite(c *C) {
	// Create random file for storing sqlDb
	db_f, err := ioutil.TempFile("/tmp", "schmoopy_test_db")
	c.Assert(err, IsNil)

	s.dbFilename = db_f.Name()
	db_f.Close()
}

func (s *SchmoopySuite) TearDownSuite(c *C) {
	os.Remove(s.dbFile)
}

func (s *SchmoopySuite) SetUpTest(c *C) {
	InitializeDb(s.dbFilename)

	server, err := NewSchmoopyServer(s.dbFile, nil)
	c.Assert(err, IsNil)

	s.server = server.(*schmoopyServer)
}
