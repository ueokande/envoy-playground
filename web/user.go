package web

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	core "github.com/ueokande/envoy-playground"
	"github.com/ueokande/envoy-playground/db"
)

func (i *impl) handleUserIndex(w http.ResponseWriter, r *http.Request) {
	us, err := i.db.ListUsers(r.Context())
	if err != nil {
		i.renderMessage(w, 500, err.Error())
		return
	}
	err = json.NewEncoder(w).Encode(us)
	if err != nil {
		i.renderMessage(w, 500, err.Error())
		return
	}
}

func (i *impl) handleUserGet(w http.ResponseWriter, r *http.Request) {
	login := mux.Vars(r)["login"]
	u, err := i.db.GetUser(r.Context(), login)
	if err == db.ErrNotFound {
		i.renderMessage(w, 404, "user not found: "+login)
		return
	} else if err != nil {
		i.renderMessage(w, 500, err.Error())
		return
	}
	i.renderJson(w, u)
}

func (i *impl) handleUserAdd(w http.ResponseWriter, r *http.Request) {
	var u core.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		i.renderMessage(w, 500, err.Error())
		return
	}
	err = u.Validate()
	if err != nil {
		i.renderMessage(w, 400, err.Error())
		return
	}
	err = i.db.AddUser(r.Context(), u)
	if err == db.ErrConflict {
		i.renderMessage(w, 409, "user already exists: "+u.Login)
		return
	} else if err != nil {
		i.renderMessage(w, 500, err.Error())
		return
	}
	u2, err := i.db.GetUser(r.Context(), u.Login)
	if err != nil {
		i.renderMessage(w, 500, err.Error())
		return
	}
	i.renderJson(w, u2)
}

func (i *impl) handleUserUpdate(w http.ResponseWriter, r *http.Request) {
	login := mux.Vars(r)["login"]

	var u core.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		i.renderMessage(w, 500, err.Error())
		return
	}
	if login != u.Login {
		i.renderMessage(w, 400, "incorrect login field")
		return
	}
	err = u.Validate()
	if err != nil {
		i.renderMessage(w, 400, err.Error())
		return
	}
	err = i.db.UpdateUser(r.Context(), u)
	if err == db.ErrNotFound {
		i.renderMessage(w, 404, "use not found: "+login)
		return
	} else if err != nil {
		i.renderMessage(w, 500, err.Error())
		return
	}
	u2, err := i.db.GetUser(r.Context(), u.Login)
	if err != nil {
		i.renderMessage(w, 500, err.Error())
		return
	}
	i.renderJson(w, u2)
}

func (i *impl) handleUserDelete(w http.ResponseWriter, r *http.Request) {
	login := mux.Vars(r)["login"]
	err := i.db.RemoveUser(r.Context(), login)
	if err == db.ErrNotFound {
		i.renderMessage(w, 404, "use not found: "+login)
		return
	} else if err != nil {
		i.renderMessage(w, 500, err.Error())
		return
	}
	i.renderMessage(w, 200, "removed: "+login)
}
