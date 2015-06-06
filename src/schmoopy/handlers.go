package schmoopy

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *schmoopyServer) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

/* Handles viewing Schmoopys or redirecting to upload if it doesn't exist. */
func (s *schmoopyServer) schmoopyHandler(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["schmoopy"]
	schmoopy, err := s.fetchSchmoopy(name)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("err: %v", err)))
	} else if schmoopy == nil {
		w.Write([]byte(fmt.Sprintf("missing: %v", name)))
	} else {
		w.Write([]byte(fmt.Sprintf("found: %v", schmoopy)))
	}
}

/* API for adding a new Schmoopy imageUrl */
func (s *schmoopyServer) addHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	imageUrl := r.FormValue("imageUrl")

	err := s.addSchmoopy(name, imageUrl)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("err: %v", err)))
	} else {
		w.Write([]byte(fmt.Sprintf("added %v for %v", imageUrl, name)))
	}
}

/* API for removing a Schmoopy imageUrl */
func (s *schmoopyServer) removeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("remove"))
}
