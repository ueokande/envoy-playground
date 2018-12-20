package web

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ueokande/envoy-playground/blob"
	"github.com/ueokande/envoy-playground/db"
)

func (i *impl) handleUserPhotoGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	login := mux.Vars(r)["login"]

	_, err := i.db.GetUser(ctx, login)
	if err == db.ErrNotFound {
		i.renderMessage(w, 404, "user not found: "+login)
		return
	} else if err != nil {
		i.renderMessage(w, 500, err.Error())
		return
	}

	uuid, err := i.db.GetPhoto(ctx, login)
	if err == db.ErrNotFound {
		i.renderMessage(w, 404, "photo not set: "+login)
		return
	} else if err != nil {
		i.renderMessage(w, 500, err.Error())
		return
	}

	b, err := i.blob.Get(ctx, uuid)
	if err == blob.ErrNotFound {
		i.renderMessage(w, 404, "blob not found: "+login)
		return
	} else if err != nil {
		i.renderMessage(w, 500, err.Error())
		return
	}

	_, err = io.Copy(w, b)

}

func (i *impl) handleUserPhotoPut(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	login := mux.Vars(r)["login"]

	_, err := i.db.GetUser(ctx, login)
	if err == db.ErrNotFound {
		i.renderMessage(w, 404, "user not found: "+login)
		return
	} else if err != nil {
		i.renderMessage(w, 500, err.Error())
	}

	id := blob.NewID()
	err = i.blob.Put(ctx, id, r.Body)
	if err != nil {
		i.renderMessage(w, 500, err.Error())
		return
	}

	err = i.db.UpdatePhoto(ctx, login, id)
	if err != nil {
		i.renderMessage(w, 500, err.Error())
		return
	}
	i.renderMessage(w, 200, "updated user photo: "+login)
}

func (i *impl) handleUserPhotoDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	login := mux.Vars(r)["login"]

	uuid, err := i.db.GetPhoto(ctx, login)
	if err == db.ErrNotFound {
		i.renderMessage(w, 404, "photo not set")
		return
	} else if err != nil {
		i.renderMessage(w, 500, err.Error())
		return
	}
	err = i.db.RemovePhoto(ctx, login)
	if err != nil {
		i.renderMessage(w, 500, err.Error())
		return
	}
	err = i.blob.Delete(ctx, uuid)
	if err != nil {
		i.renderMessage(w, 500, err.Error())
		return
	}
	i.renderMessage(w, 200, "photo removed")
}
