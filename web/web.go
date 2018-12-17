package web

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ueokande/envoy-playground/db"
)

func New(db db.Interface) http.Handler {
	return &impl{db: db}
}

type impl struct {
	db db.Interface
}

func (i *impl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/users":
		i.handleUserIndex(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/user/"):
		i.handleUserGet(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/users":
		i.handleUserAdd(w, r)
	case r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/user/"):
		i.handleUserUpdate(w, r)
	case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/user/"):
		i.handleUserDelete(w, r)
	}
	i.renderMessage(w, 404, "path not found")
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
