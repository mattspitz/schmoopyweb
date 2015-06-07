package schmoopy

import (
	"log"
	"net/http"

	"code.google.com/p/go-sqlite/go1/sqlite3"
	"github.com/gorilla/mux"
)

type SchmoopyServer interface {
	ListenAndServe() error
}

type schmoopyServer struct {
	conn         *sqlite3.Conn
	addr         string
	templateDir  string
	staticDir    string
	dbAccessChan chan *dbAccessEvent
}

func NewSchmoopyServer(
	dbFilename string,
	addr string,
	templateDir string,
	staticDir string,
) (SchmoopyServer, error) {
	conn, err := sqlite3.Open(dbFilename)
	if err != nil {
		return nil, err
	}

	s := &schmoopyServer{
		conn:         conn,
		addr:         addr,
		templateDir:  templateDir,
		staticDir:    staticDir,
		dbAccessChan: make(chan *dbAccessEvent),
	}

	go s.serializeDbAccess()
	return s, nil
}

func (s *schmoopyServer) ListenAndServe() error {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/schmoopy").
		Methods("POST").
		Subrouter()
	api.HandleFunc("/add", s.addHandler)
	api.HandleFunc("/remove", s.removeHandler)

	if s.staticDir != "" {
		fs := http.FileServer(http.Dir(s.staticDir))
		r.PathPrefix("/static/").
			Handler(http.StripPrefix("/static/", fs))
	}

	r.HandleFunc("/{schmoopy}", s.schmoopyHandler)
	r.HandleFunc("/", s.indexHandler)

	return http.ListenAndServe(s.addr, r)
}

type dbAccessType string

const (
	query  dbAccessType = "query"
	add                 = "add"
	remove              = "remove"
)

type dbAccessEvent struct {
	accessType        dbAccessType
	queryNames        []string
	addRemoveName     string
	addRemoveImageUrl string
	res               chan *dbAccessResponse // the result of the query
}

type dbAccessResponse struct {
	schmoopys map[string]*schmoopy
	err       error
}

func (s *schmoopyServer) serializeDbAccess() {
	for {
		select {
		case event := <-s.dbAccessChan:
			var response *dbAccessResponse

			switch event.accessType {
			case query:
				schmoopys, err := dbFetchSchmoopys(s.conn, event.queryNames)
				response = &dbAccessResponse{
					schmoopys: schmoopys,
					err:       err,
				}
				break

			case add:
				response = &dbAccessResponse{
					err: dbAddSchmoopy(s.conn, event.addRemoveName, event.addRemoveImageUrl),
				}
				break

			case remove:
				response = &dbAccessResponse{
					err: dbRemoveSchmoopy(s.conn, event.addRemoveName, event.addRemoveImageUrl),
				}
				break

			default:
				log.Fatal("Unhandled access type: %v", event.accessType)
			}

			event.res <- response
		}
	}
}

func (s *schmoopyServer) fetchSchmoopy(name string) (*schmoopy, error) {
	schmoopys, err := s.fetchSchmoopys([]string{name})
	if err != nil {
		return nil, err
	}

	return schmoopys[name], nil
}

func (s *schmoopyServer) fetchSchmoopys(names []string) (map[string]*schmoopy, error) {
	res := make(chan *dbAccessResponse, 1)
	s.dbAccessChan <- &dbAccessEvent{
		accessType: query,
		queryNames: names,
		res:        res,
	}

	select {
	case response := <-res:
		if response.err != nil {
			return nil, response.err
		}
		return response.schmoopys, nil
	}
}

func (s *schmoopyServer) fetchAllSchmoopys() (map[string]*schmoopy, error) {
	return s.fetchSchmoopys(nil)
}

func (s *schmoopyServer) addSchmoopy(name string, imageUrl string) error {
	res := make(chan *dbAccessResponse, 1)
	s.dbAccessChan <- &dbAccessEvent{
		accessType:        add,
		addRemoveName:     name,
		addRemoveImageUrl: imageUrl,
		res:               res,
	}

	select {
	case response := <-res:
		return response.err
	}
}

func (s *schmoopyServer) removeSchmoopy(name string, imageUrl string) error {
	res := make(chan *dbAccessResponse, 1)
	s.dbAccessChan <- &dbAccessEvent{
		accessType:        remove,
		addRemoveName:     name,
		addRemoveImageUrl: imageUrl,
		res:               res,
	}

	select {
	case response := <-res:
		return response.err
	}
}
