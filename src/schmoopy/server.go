package schmoopy

import (
	"net/http"

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
	api.HandleFunc("/add", s.addHandler)
	api.HandleFunc("/remove", s.removeHandler)

	r.HandleFunc("/{schmoopy}", s.schmoopyHandler)
	r.HandleFunc("/", s.indexHandler)

	return http.ListenAndServe(s.addr, r)
}
