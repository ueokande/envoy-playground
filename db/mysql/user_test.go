package mysql

import (
	"context"
	"testing"

	core "github.com/ueokande/envoy-playground"
	"github.com/ueokande/envoy-playground/db"
)

func TestAddUser(t *testing.T) {
	ctx := context.Background()

	sql, err := newDB()
	if err != nil {
		t.Fatal(err)
	}
	defer sql.Close()
	d := New(sql)

	err = d.AddUser(ctx, core.User{Login: "alice", Name: "Alice"})
	if err != nil {
		t.Fatal(err)
	}
	err = d.AddUser(ctx, core.User{Login: "bob", Name: "Bob"})
	if err != nil {
		t.Fatal(err)
	}

	err = d.AddUser(ctx, core.User{Login: "alice", Name: "Alice2"})
	if err != db.ErrConflict {
		t.Fatal(err)
	}
	err = d.AddUser(ctx, core.User{Login: "ALICE", Name: "Alice3"})
	if err != db.ErrConflict {
		t.Fatal(err)
	}
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()

	sql, err := newDB()
	if err != nil {
		t.Fatal(err)
	}
	defer sql.Close()
	d := New(sql)

	err = d.AddUser(ctx, core.User{Login: "carol", Name: "Carol"})
	if err != nil {
		t.Fatal(err)
	}

	u, err := d.GetUser(ctx, "carol")
	if err != nil {
		t.Fatal(err)
	}
	if u.Login != "carol" || u.Name != "Carol" {
		t.Fatal("unexpected user", u)
	}

	u, err = d.GetUser(ctx, "ghost")
	if err != db.ErrNotFound {
		t.Fatal(err)
	}
}

func TestListUser(t *testing.T) {
	ctx := context.Background()

	sql, err := newDB()
	if err != nil {
		t.Fatal(err)
	}
	defer sql.Close()
	d := New(sql)

	err = d.AddUser(ctx, core.User{Login: "dan", Name: "Dan"})
	if err != nil {
		t.Fatal(err)
	}
	err = d.AddUser(ctx, core.User{Login: "erin", Name: "Erin"})
	if err != nil {
		t.Fatal(err)
	}

	us, err := d.ListUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, u := range us {
		if u.Login == "dan" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("dan not found")
	}
	found = false
	for _, u := range us {
		if u.Login == "erin" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("erin not found")
	}
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()

	sql, err := newDB()
	if err != nil {
		t.Fatal(err)
	}
	defer sql.Close()
	d := New(sql)

	err = d.AddUser(ctx, core.User{Login: "faythe", Name: "Faythe"})
	if err != nil {
		t.Fatal(err)
	}

	u, err := d.GetUser(ctx, "faythe")
	if err != nil {
		t.Fatal(err)
	}

	u.Name = "Faythe's mom"
	err = d.UpdateUser(ctx, *u)
	if err != nil {
		t.Fatal(err)
	}

	u, err = d.GetUser(ctx, "faythe")
	if err != nil {
		t.Fatal(err)
	}

	if u.Name != "Faythe's mom" {
		t.Error(`u.Name != "Faythe's mom"`, u.Name)
	}
}

func TestRemoveUser(t *testing.T) {
	ctx := context.Background()

	sql, err := newDB()
	if err != nil {
		t.Fatal(err)
	}
	defer sql.Close()
	d := New(sql)

	err = d.AddUser(ctx, core.User{Login: "Grace", Name: "grace"})
	if err != nil {
		t.Fatal(err)
	}

	err = d.RemoveUser(ctx, "grace")
	if err != nil {
		t.Fatal(err)
	}
	_, err = d.GetUser(ctx, "grace")
	if err != db.ErrNotFound {
		t.Fatal(err)
	}
	err = d.RemoveUser(ctx, "grace")
	if err != db.ErrNotFound {
		t.Fatal(err)
	}
}
