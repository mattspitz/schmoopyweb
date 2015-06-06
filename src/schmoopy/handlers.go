package schmoopy

import (
	"encoding/json"
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

type apiResponse struct {
	Ok  bool   `json:"ok,string"`
	Err string `json:"err,omitempty"`
}

func writeApiResponse(w http.ResponseWriter, err error) {
	response := apiResponse{
		Ok: err == nil,
	}
	if err != nil {
		response.Err = err.Error()
	}

	s, err := json.Marshal(response)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("err: %v", err)))
	} else {
		w.Write(s)
	}
}

/* API for adding a new Schmoopy imageUrl */
func (s *schmoopyServer) addHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	imageUrl := r.FormValue("imageUrl")

	err := s.addSchmoopy(name, imageUrl)
	writeApiResponse(w, err)
}

/* API for removing a Schmoopy imageUrl */
func (s *schmoopyServer) removeHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	imageUrl := r.FormValue("imageUrl")

	err := s.removeSchmoopy(name, imageUrl)
	writeApiResponse(w, err)
}
