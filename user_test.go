package core

import "testing"

func TestUserValidate(t *testing.T) {
	valids := []User{
		{Login: "alice", Name: "Alice in Wonderland"},
		{Login: "alice", Name: "ありす"},
	}
	for _, u := range valids {
		err := u.Validate()
		if err != nil {
			t.Error("err != nil:", err)
		}

	}

	for _, login := range []string{"", "bob bob", "m@ster", "ありす"} {
		u := User{Login: login, Name: "Alice in Wonderland"}
		err := u.Validate()
		if err == nil {
			t.Error("valid Login:", login)
		}
	}
	for _, name := range []string{"", "Looking\tGlass", "Looking\nGlass"} {
		u := User{Login: "alice", Name: name}
		err := u.Validate()
		if err == nil {
			t.Error("valid Name:", name)
		}
	}
}
