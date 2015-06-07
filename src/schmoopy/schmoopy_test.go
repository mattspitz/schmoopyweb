package schmoopy

import (
	"io/ioutil"
	"os"
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

func (s *SchmoopySuite) SetUpTest(c *C) {
	// Create random file for storing sqlDb
	db_f, err := ioutil.TempFile("/tmp", "schmoopy_test_db")
	c.Assert(err, IsNil)

	s.dbFilename = db_f.Name()
	db_f.Close()
	// Just need the filename; initializing the database expects the file not to be present!
	os.Remove(s.dbFilename)

	InitializeDb(s.dbFilename)

	server, err := NewSchmoopyServer(s.dbFilename, "<addr>", "<templateDir>", "<staticDir>")
	c.Assert(err, IsNil)

	s.server = server.(*schmoopyServer)
}

func (s *SchmoopySuite) TearDownTest(c *C) {
	os.Remove(s.dbFilename)
}

func (s *SchmoopySuite) TestAdd(c *C) {
	sch, err := s.server.fetchSchmoopy("missing")
	c.Assert(sch, IsNil)
	c.Assert(err, IsNil)

	err = s.server.addSchmoopy("present", "with-url")
	c.Assert(err, IsNil)

	sch, err = s.server.fetchSchmoopy("present")
	c.Assert(*sch, DeepEquals, schmoopy{
		Name: "present",
		ImageUrls: map[string]struct{}{
			"with-url": struct{}{},
		},
	})
	c.Assert(err, IsNil)

	err = s.server.addSchmoopy("present", "with-another-url")
	c.Assert(err, IsNil)

	sch, err = s.server.fetchSchmoopy("present")
	c.Assert(*sch, DeepEquals, schmoopy{
		Name: "present",
		ImageUrls: map[string]struct{}{
			"with-url":         struct{}{},
			"with-another-url": struct{}{},
		},
	})
	c.Assert(err, IsNil)

	// adding the same URL does nothing
	err = s.server.addSchmoopy("present", "with-another-url")
	c.Assert(err, IsNil)

	sch, err = s.server.fetchSchmoopy("present")
	c.Assert(*sch, DeepEquals, schmoopy{
		Name: "present",
		ImageUrls: map[string]struct{}{
			"with-url":         struct{}{},
			"with-another-url": struct{}{},
		},
	})
	c.Assert(err, IsNil)
}

func (s *SchmoopySuite) TestRemoveMissing(c *C) {
	err := s.server.addSchmoopy("one", "url1")
	c.Assert(err, IsNil)

	schmoopys, err := s.server.fetchAllSchmoopys()
	c.Assert(err, IsNil)

	c.Assert(schmoopys, HasLen, 1)

	c.Assert(*schmoopys["one"], DeepEquals, schmoopy{
		Name: "one",
		ImageUrls: map[string]struct{}{
			"url1": struct{}{},
		},
	})

	err = s.server.removeSchmoopy("one", "nope")
	c.Assert(err, IsNil)

	schmoopys, err = s.server.fetchAllSchmoopys()
	c.Assert(err, IsNil)

	c.Assert(schmoopys, HasLen, 1)

	c.Assert(*schmoopys["one"], DeepEquals, schmoopy{
		Name: "one",
		ImageUrls: map[string]struct{}{
			"url1": struct{}{},
		},
	})
}

func (s *SchmoopySuite) TestAddRemove(c *C) {
	schmoopys, err := s.server.fetchAllSchmoopys()
	c.Assert(err, IsNil)

	c.Assert(schmoopys, HasLen, 0)

	err = s.server.addSchmoopy("one", "url1")
	c.Assert(err, IsNil)
	err = s.server.addSchmoopy("two", "url2")
	c.Assert(err, IsNil)
	err = s.server.addSchmoopy("three", "url3")
	c.Assert(err, IsNil)
	err = s.server.addSchmoopy("three", "url4")
	c.Assert(err, IsNil)

	schmoopys, err = s.server.fetchAllSchmoopys()
	c.Assert(err, IsNil)

	c.Assert(schmoopys, HasLen, 3)

	c.Assert(*schmoopys["one"], DeepEquals, schmoopy{
		Name: "one",
		ImageUrls: map[string]struct{}{
			"url1": struct{}{},
		},
	})
	c.Assert(*schmoopys["two"], DeepEquals, schmoopy{
		Name: "two",
		ImageUrls: map[string]struct{}{
			"url2": struct{}{},
		},
	})
	c.Assert(*schmoopys["three"], DeepEquals, schmoopy{
		Name: "three",
		ImageUrls: map[string]struct{}{
			"url3": struct{}{},
			"url4": struct{}{},
		},
	})

	// Removing the last URL removes the whole schmoopy
	c.Assert(s.server.removeSchmoopy("two", "url2"), IsNil)

	schmoopys, err = s.server.fetchAllSchmoopys()
	c.Assert(err, IsNil)

	c.Assert(schmoopys, HasLen, 2)

	c.Assert(*schmoopys["one"], DeepEquals, schmoopy{
		Name: "one",
		ImageUrls: map[string]struct{}{
			"url1": struct{}{},
		},
	})
	c.Assert(*schmoopys["three"], DeepEquals, schmoopy{
		Name: "three",
		ImageUrls: map[string]struct{}{
			"url3": struct{}{},
			"url4": struct{}{},
		},
	})

	// Removing not-the-last URL doesn't remove the schmoopy
	c.Assert(s.server.removeSchmoopy("three", "url3"), IsNil)

	schmoopys, err = s.server.fetchAllSchmoopys()
	c.Assert(err, IsNil)

	c.Assert(schmoopys, HasLen, 2)

	c.Assert(*schmoopys["one"], DeepEquals, schmoopy{
		Name: "one",
		ImageUrls: map[string]struct{}{
			"url1": struct{}{},
		},
	})
	c.Assert(*schmoopys["three"], DeepEquals, schmoopy{
		Name: "three",
		ImageUrls: map[string]struct{}{
			"url4": struct{}{},
		},
	})
}

func (s *SchmoopySuite) TestFetchSchmoopys(c *C) {
	err := s.server.addSchmoopy("one", "url1")
	c.Assert(err, IsNil)
	err = s.server.addSchmoopy("two", "url2")
	c.Assert(err, IsNil)
	err = s.server.addSchmoopy("three", "url3")
	c.Assert(err, IsNil)
	err = s.server.addSchmoopy("three", "url4")
	c.Assert(err, IsNil)

	schmoopys, err := dbFetchSchmoopys(s.server.conn, []string{"one", "three"})

	c.Assert(schmoopys, HasLen, 2)

	c.Assert(*schmoopys["one"], DeepEquals, schmoopy{
		Name: "one",
		ImageUrls: map[string]struct{}{
			"url1": struct{}{},
		},
	})
	c.Assert(*schmoopys["three"], DeepEquals, schmoopy{
		Name: "three",
		ImageUrls: map[string]struct{}{
			"url3": struct{}{},
			"url4": struct{}{},
		},
	})
}
