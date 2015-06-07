package schmoopy

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

func (s *schmoopyServer) renderTemplate(w http.ResponseWriter, templateFn string, data map[string]interface{}) {
	t, err := template.ParseFiles(path.Join(s.templateDir, templateFn))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *schmoopyServer) indexHandler(w http.ResponseWriter, r *http.Request) {
	schmoopys, err := s.fetchAllSchmoopys()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// randomly shuffle the schmoopys
	schmoopySlice := []*schmoopy{}
	for _, v := range schmoopys {
		schmoopySlice = append(schmoopySlice, v)
	}
	for i := range schmoopySlice {
		j := rand.Intn(i + 1)
		schmoopySlice[i], schmoopySlice[j] = schmoopySlice[j], schmoopySlice[i]
	}

	data := map[string]interface{}{
		"schmoopys": schmoopySlice,
	}
	s.renderTemplate(w, "index.html", data)
}

/* Handles viewing Schmoopys or redirecting to upload if it doesn't exist. */
func (s *schmoopyServer) schmoopyHandler(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["schmoopy"]
	schmoopy, err := s.fetchSchmoopy(name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if schmoopy != nil {
		data := map[string]interface{}{
			"schmoopy": schmoopy,
		}
		s.renderTemplate(w, "schmoopy.html", data)
	} else {
		data := map[string]interface{}{
			"name": name,
		}
		s.renderTemplate(w, "create.html", data)
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
	if name == "" || imageUrl == "" {
		writeApiResponse(w, errors.New("Must provide non-empty 'name' and 'imageUrl'"))
		return
	}

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
