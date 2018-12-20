package web

import (
	"context"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	core "github.com/ueokande/envoy-playground"
	"github.com/ueokande/envoy-playground/blob"
	mockBlob "github.com/ueokande/envoy-playground/blob/mock"
	"github.com/ueokande/envoy-playground/db"
	mockDB "github.com/ueokande/envoy-playground/db/mock"
)

func TestUserPhotoGet(t *testing.T) {
	ctx := context.Background()

	d := mockDB.New()
	d.AddUser(ctx, core.User{Login: "alice", Name: "Alice"})
	d.AddUser(ctx, core.User{Login: "bob", Name: "Bob"})
	d.UpdatePhoto(ctx, "alice", "0000-0000")

	b := mockBlob.New()
	b.Put(ctx, "0000-0000", strings.NewReader("raw:xxxxxxxx"))

	h := New(d, b)
	r := httptest.NewRequest("GET", "/user/alice/photo", nil).
		WithContext(context.WithValue(ctx, "login", "alice"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatal(`w.Code != 200: `, w.Code, w.Body.String())
	}
	if body := w.Body.String(); body != "raw:xxxxxxxx" {
		t.Error(`body != "raw:xxxxxxxx":`, body)
	}

	r = httptest.NewRequest("GET", "/user/bob/photo", nil).
		WithContext(context.WithValue(ctx, "login", "bob"))
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 404 {
		t.Fatal(`w.Code != 404: `, w.Code, w.Body.String())
	}

	r = httptest.NewRequest("GET", "/user/gone/photo", nil).
		WithContext(context.WithValue(ctx, "login", "gone"))
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 404 {
		t.Fatal(`w.Code != 404: `, w.Code, w.Body.String())
	}
}

func TestUserPhotoPut(t *testing.T) {
	ctx := context.Background()

	d := mockDB.New()
	d.AddUser(ctx, core.User{Login: "alice", Name: "Alice"})
	d.AddUser(ctx, core.User{Login: "bob", Name: "Bob"})
	d.UpdatePhoto(ctx, "alice", "0000-0000")

	b := mockBlob.New()
	b.Put(ctx, "0000-0000", strings.NewReader("raw:xxxxxxxx"))

	h := New(d, b)
	r := httptest.NewRequest("PUT", "/user/alice/photo", strings.NewReader("raw:yyyyyyyy")).
		WithContext(context.WithValue(ctx, "login", "alice"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatal(`w.Code != 200: `, w.Code, w.Body.String())
	}
	uuid, _ := d.GetPhoto(ctx, "alice")
	o, _ := b.Get(ctx, uuid)
	bytes, _ := ioutil.ReadAll(o)
	if string(bytes) != "raw:yyyyyyyy" {
		t.Error(`string(bytes) != "raw:yyyyyyyy":`, string(bytes))
	}

	r = httptest.NewRequest("PUT", "/user/bob/photo", strings.NewReader("raw:zzzzzzzz")).
		WithContext(context.WithValue(ctx, "login", "bob"))
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatal(`w.Code != 200: `, w.Code, w.Body.String())
	}
	uuid, _ = d.GetPhoto(ctx, "bob")
	o, _ = b.Get(ctx, uuid)
	bytes, _ = ioutil.ReadAll(o)
	if string(bytes) != "raw:zzzzzzzz" {
		t.Error(`string(bytes) != "raw:zzzzzzzz":`, string(bytes))
	}

	r = httptest.NewRequest("PUT", "/user/gone/photo", strings.NewReader("raw:zzzzzzzz")).
		WithContext(context.WithValue(ctx, "login", "gone"))
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 404 {
		t.Fatal(`w.Code != 404: `, w.Code, w.Body.String())
	}
}

func TestUserPhotoDelete(t *testing.T) {
	ctx := context.Background()

	d := mockDB.New()
	d.AddUser(ctx, core.User{Login: "alice", Name: "Alice"})
	d.AddUser(ctx, core.User{Login: "bob", Name: "Bob"})
	d.UpdatePhoto(ctx, "alice", "0000-0000")

	b := mockBlob.New()
	b.Put(ctx, "0000-0000", strings.NewReader("raw:xxxxxxxx"))

	h := New(d, b)
	r := httptest.NewRequest("DELETE", "/user/alice/photo", nil).
		WithContext(context.WithValue(ctx, "login", "alice"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatal(`w.Code != 200: `, w.Code, w.Body.String())
	}
	_, err := d.GetPhoto(ctx, "alice")
	if err != db.ErrNotFound {
		t.Error("err != db.ErrNotFound", err)
	}
	_, err = b.Get(ctx, "0000-0000")
	if err != blob.ErrNotFound {
		t.Error("err != blob.ErrNotFound", err)
	}

	r = httptest.NewRequest("DELETE", "/user/bob/photo", nil).
		WithContext(context.WithValue(ctx, "login", "bob"))
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 404 {
		t.Fatal(`w.Code != 404: `, w.Code, w.Body.String())
	}

	r = httptest.NewRequest("DELETE", "/user/gone/photo", nil).
		WithContext(context.WithValue(ctx, "login", "gone"))
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 404 {
		t.Fatal(`w.Code != 404: `, w.Code, w.Body.String())
	}
}
