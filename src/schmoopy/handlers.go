package schmoopy

import (
	"net/http"

	"github.com/gorilla/mux"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

/* Handles viewing Schmoopys or redirecting to upload if it doesn't exist. */
func schmoopyHandler(w http.ResponseWriter, r *http.Request) {
	schmoopy := mux.Vars(r)["schmoopy"]
	w.Write([]byte(schmoopy))
}

/* API for adding a new Schmoopy imageUrl */
func addHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("add"))
}

/* API for removing a Schmoopy imageUrl */
func removeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("remove"))
}
