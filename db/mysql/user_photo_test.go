package mysql

import (
	"context"
	"testing"

	"github.com/ueokande/envoy-playground/db"
)

func TestGetPhoto(t *testing.T) {
	ctx := context.Background()

	sql, err := initDB()
	if err != nil {
		t.Fatal(err)
	}
	defer sql.Close()
	d := New(sql)

	_, err = sql.Exec(
		`INSERT INTO user(login, name) VALUES ('alice', 'ALICE'), ('bob', 'BOB')`,
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = sql.Exec(`INSERT INTO user_photo(user_id,uuid) SELECT id,'aaaa-bbbb' FROM user where login='alice'; `)
	if err != nil {
		t.Fatal(err)
	}

	uuid, err := d.GetPhoto(ctx, "alice")
	if err != nil {
		t.Fatal(err)
	}
	if uuid != "aaaa-bbbb" {
		t.Error(`uuid != "aaaa-bbbb": `, uuid)
	}

	_, err = d.GetPhoto(ctx, "bob")
	if err != db.ErrNotFound {
		t.Fatal(err)
	}
}

func TestUpdatePhoto(t *testing.T) {
	ctx := context.Background()

	sql, err := initDB()
	if err != nil {
		t.Fatal(err)
	}
	defer sql.Close()
	d := New(sql)

	_, err = sql.Exec(
		`INSERT INTO user(login, name) VALUES ('alice', 'ALICE'), ('bob', 'BOB')`,
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = sql.Exec(`INSERT INTO user_photo(user_id,uuid) SELECT id,'aaaa-bbbb' FROM user where login='alice'; `)
	if err != nil {
		t.Fatal(err)
	}

	err = d.UpdatePhoto(ctx, "alice", "cccc-dddd")
	if err != nil {
		t.Fatal(err)
	}
	err = d.UpdatePhoto(ctx, "bob", "xxxx-yyyy")
	if err != nil {
		t.Fatal(err)
	}

	uuid, err := d.GetPhoto(ctx, "alice")
	if err != err {
		t.Fatal(err)
	}
	if uuid != "cccc-dddd" {
		t.Error(`uuid != "cccc-dddd": `, uuid)
	}

	uuid, err = d.GetPhoto(ctx, "bob")
	if err != err {
		t.Fatal(err)
	}
	if uuid != "xxxx-yyyy" {
		t.Error(`uuid != "xxxx-yyyy": `, uuid)
	}
}

func TestRemovePhoto(t *testing.T) {
	ctx := context.Background()

	sql, err := initDB()
	if err != nil {
		t.Fatal(err)
	}
	defer sql.Close()
	d := New(sql)

	_, err = sql.Exec(
		`INSERT INTO user(login, name) VALUES ('alice', 'ALICE'), ('bob', 'BOB')`,
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = sql.Exec(`INSERT INTO user_photo(user_id,uuid) SELECT id,'aaaa-bbbb' FROM user where login='alice'; `)
	if err != nil {
		t.Fatal(err)
	}

	err = d.RemovePhoto(ctx, "alice")
	if err != nil {
		t.Fatal(err)
	}
	err = d.RemovePhoto(ctx, "bob")
	if err != db.ErrNotFound {
		t.Fatal(err)
	}
}
