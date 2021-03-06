package web

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	core "github.com/ueokande/envoy-playground"
	"github.com/ueokande/envoy-playground/blob"
	mockBlob "github.com/ueokande/envoy-playground/blob/mock"
	"github.com/ueokande/envoy-playground/db"
	mockDB "github.com/ueokande/envoy-playground/db/mock"
)

func TestUserIndex(t *testing.T) {
	ctx := context.Background()
	d := mockDB.New()
	d.AddUser(ctx, core.User{Login: "alice", Name: "Alice"})
	d.AddUser(ctx, core.User{Login: "bob", Name: "Bob"})

	h := New(d, mockBlob.New())

	r := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatal(`w.Result().StatusCode != 200: `, w.Code)
	}

	var us []core.User
	err := json.NewDecoder(w.Body).Decode(&us)
	if err != nil {
		t.Fatal(err)
	}

	sort.Slice(us, func(i, j int) bool { return us[i].Name < us[j].Name })
	if len(us) != 2 {
		t.Error("len(us) != 2: ", len(us))
	}
}

func TestUserGet(t *testing.T) {
	ctx := context.Background()
	d := mockDB.New()
	d.AddUser(ctx, core.User{Login: "alice", Name: "Alice"})
	h := New(d, mockBlob.New())

	r := httptest.NewRequest("GET", "/user/alice", nil).
		WithContext(context.WithValue(ctx, "login", "alice"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatal(`w.Result().StatusCode != 200: `, w.Code)
	}

	var u core.User
	err := json.NewDecoder(w.Body).Decode(&u)
	if err != nil {
		t.Fatal(err)
	}

	if u.Login != "alice" || u.Name != "Alice" {
		t.Error("unexpected user: ", u)
	}

	r = httptest.NewRequest("GET", "/user/ghost", nil).
		WithContext(context.WithValue(ctx, "login", "ghost"))
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 404 {
		t.Fatal(`w.Result().StatusCode != 404: `, w.Code)
	}
}

func TestUserAdd(t *testing.T) {
	ctx := context.Background()

	d := mockDB.New()
	h := New(d, mockBlob.New())

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(core.User{Login: "alice", Name: "Alice"})
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest("POST", "/users", b)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatal(`w.Result().StatusCode != 200: `, w.Code)
	}
	u, err := d.GetUser(ctx, "alice")
	if err != nil {
		t.Fatal(err)
	}
	if u.Login != "alice" || u.Name != "Alice" {
		t.Error("unexpected user: ", u)
	}

	b = new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(core.User{Login: "alice", Name: "Alice"})
	if err != nil {
		t.Fatal(err)
	}
	r = httptest.NewRequest("POST", "/users", b)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 409 {
		t.Fatal(`w.Result().StatusCode != 409: `, w.Code)
	}
}

func TestUserUpdate(t *testing.T) {
	ctx := context.Background()
	d := mockDB.New()
	d.AddUser(ctx, core.User{Login: "alice", Name: "Alice"})
	h := New(d, mockBlob.New())

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(core.User{Login: "alice", Name: "Alice 2nd"})
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest("PUT", "/user/alice", b).
		WithContext(context.WithValue(ctx, "login", "alice"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatal(`w.Result().StatusCode != 200: `, w.Code)
	}

	u, err := d.GetUser(ctx, "alice")
	if err != nil {
		t.Fatal(err)
	}
	if u.Login != "alice" || u.Name != "Alice 2nd" {
		t.Error("unexpected user: ", u)
	}

	b = new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(core.User{Login: "alice", Name: "Alice"})
	if err != nil {
		t.Fatal(err)
	}
	r = httptest.NewRequest("PUT", "/user/bob", b).
		WithContext(context.WithValue(ctx, "login", "bob"))
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Fatal(`w.Result().StatusCode != 400: `, w.Code)
	}

	b = new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(core.User{Login: "ghost", Name: "Ghost"})
	if err != nil {
		t.Fatal(err)
	}
	r = httptest.NewRequest("PUT", "/user/ghost", b).
		WithContext(context.WithValue(ctx, "login", "ghost"))
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 404 {
		t.Fatal(`w.Result().StatusCode != 404: `, w.Code)
	}
}

func TestUserDelete(t *testing.T) {
	ctx := context.Background()
	d := mockDB.New()
	b := mockBlob.New()
	h := New(d, b)

	d.AddUser(ctx, core.User{Login: "alice", Name: "Alice"})
	d.UpdatePhoto(ctx, "alice", "0000-0000")
	d.AddUser(ctx, core.User{Login: "bob", Name: "Bob"})
	b.Put(ctx, "0000-0000", strings.NewReader("raw:xxxxxxxx"))

	r := httptest.NewRequest("DELETE", "/user/alice", nil).
		WithContext(context.WithValue(ctx, "login", "alice"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatal(`w.Result().StatusCode != 200: `, w.Code)
	}

	_, err := d.GetUser(ctx, "alice")
	if err != db.ErrNotFound {
		t.Fatal(err)
	}
	_, err = d.GetPhoto(ctx, "alice")
	if err != db.ErrNotFound {
		t.Fatal(err)
	}
	_, err = b.Get(ctx, "0000-0000")
	if err != blob.ErrNotFound {
		t.Fatal(err)
	}

	r = httptest.NewRequest("DELETE", "/user/bob", nil).
		WithContext(context.WithValue(ctx, "login", "bob"))
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatal(`w.Result().StatusCode != 200: `, w.Code)
	}

	r = httptest.NewRequest("DELETE", "/user/alice", nil).
		WithContext(context.WithValue(ctx, "login", "bob"))
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 404 {
		t.Fatal(`w.Result().StatusCode != 404: `, w.Code)
	}
}
