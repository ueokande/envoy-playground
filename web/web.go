package web

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ueokande/envoy-playground/blob"
	"github.com/ueokande/envoy-playground/db"
)

func New(db db.Interface, blob blob.Interface) http.Handler {
	r := mux.NewRouter()
	i := &impl{r: r, db: db, blob: blob}
	i.init()
	return i
}

type impl struct {
	r    *mux.Router
	db   db.Interface
	blob blob.Interface
}

func (i *impl) init() {
	i.r.HandleFunc("/users", i.handleUserIndex).Methods("GET")
	i.r.HandleFunc("/user/{login}", i.handleUserGet).Methods("GET")
	i.r.HandleFunc("/users", i.handleUserAdd).Methods("POST")
	i.r.HandleFunc("/user/{login}", i.handleUserUpdate).Methods("PUT")
	i.r.HandleFunc("/user/{login}", i.handleUserDelete).Methods("DELETE")

	i.r.HandleFunc("/user/{login}/photo", i.handleUserPhotoGet).Methods("GET")
	i.r.HandleFunc("/user/{login}/photo", i.handleUserPhotoPut).Methods("PUT")
	i.r.HandleFunc("/user/{login}/photo", i.handleUserPhotoDelete).Methods("DELETE")
}

func (i *impl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i.r.ServeHTTP(w, r)
}

func (i *impl) renderJson(w http.ResponseWriter, data interface{}) {
	err := json.NewEncoder(w).Encode(data)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		i.renderMessage(w, 500, err.Error())
		return
	}
}

func (i *impl) renderMessage(w http.ResponseWriter, status int, reason string) {
	msg := map[string]interface{}{
		"status": status,
		"reason": reason,
	}
	j, err := json.Marshal(msg)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(status)
	w.Write(j)
}

type Message struct {
	Status int    `json:"status"`
	Reason string `json:"reason"`
}
