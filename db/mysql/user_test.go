package mysql

import (
	"context"
	"sort"
	"testing"

	core "github.com/ueokande/envoy-playground"
)

func testUserCRUD(t *testing.T) {
	ctx := context.Background()

	db, err := New(conf)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Add
	err = db.AddUser(ctx, core.User{Login: "alice", Name: "Alice"})
	if err != nil {
		t.Fatal(err)
	}
	err = db.AddUser(ctx, core.User{Login: "bob", Name: "Bob"})
	if err != nil {
		t.Fatal(err)
	}

	// List
	us, err := db.ListUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(us) != 2 {
		t.Fatal(`len(us) != 2`, len(us))
	}
	sort.Slice(us, func(i, j int) bool { return us[i].Name < us[j].Name })

	// Get
	u, err := db.GetUser(ctx, us[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if u.Login != "alice" || u.Name != "Alice" {
		t.Fatal("unexpected user", u)
	}
	u, err = db.GetUser(ctx, us[1].ID)
	if err != nil {
		t.Fatal(err)
	}
	if u.Login != "bob" || u.Name != "Bob" {
		t.Fatal("unexpected user", u)
	}

	// Update
	us[1].Name = "Bobby"
	err = db.UpdateUser(ctx, *us[1])
	if err != nil {
		t.Fatal(err)
	}
	u, err = db.GetUser(ctx, us[1].ID)
	if err != nil {
		t.Fatal(err)
	}
	if u.Login != "bob" || u.Name != "Bobby" {
		t.Fatal("unexpected user", u)
	}

	// Remove
	err = db.RemoveUser(ctx, us[1].ID)
	if err != nil {
		t.Fatal(err)
	}
	us, err = db.ListUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(us) != 1 {
		t.Fatal(`len(us) != 2`, len(us))
	}
	if us[0].Login != "alice" || us[0].Name != "Alice" {
		t.Fatal("unexpected user", us[0])
	}
}

func TestUser(t *testing.T) {
	t.Run("CRUD", testUserCRUD)
}
