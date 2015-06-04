package schmoopy

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"code.google.com/p/go-sqlite/go1/sqlite3"
	"github.com/gorilla/mux"
)

type SchmoopyServer interface {
	ListenAndServe() error
}

type schmoopyServer struct {
	conn *sqlite3.Conn
	addr string
}

func InitializeDb(
	dbFilename string,
) error {
	if _, err := os.Stat(dbFilename); err == nil {
		log.Fatal(fmt.Errorf("File already exists: %v", dbFilename))
	}

	c, err := sqlite3.Open(dbFilename)
	if err != nil {
		return err
	}

	sql := "CREATE TABLE schmoopys(schmoopy STRING, imageUrl STRING); " +
		"CREATE INDEX schmoopys_schmoopys ON schmoopys(schmoopy)"

	if err = c.Exec(sql); err != nil {
		return err
	}

	if err = c.Close(); err != nil {
		return err
	}

	return nil
}

func NewSchmoopyServer(
	dbFilename string,
	addr string,
) (SchmoopyServer, error) {
	conn, err := sqlite3.Open(dbFilename)
	if err != nil {
		return nil, err
	}

	return &schmoopyServer{
		conn: conn,
		addr: addr,
	}, nil
}

func (s *schmoopyServer) ListenAndServe() error {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/schmoopy").
		Methods("POST").
		Subrouter()
	api.HandleFunc("/add", addHandler)
	api.HandleFunc("/remove", removeHandler)

	r.HandleFunc("/{schmoopy}", schmoopyHandler)
	r.HandleFunc("/", mainHandler)

	return http.ListenAndServe(s.addr, r)
}
