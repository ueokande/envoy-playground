package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/ueokande/envoy-playground/client"
)

var count int = 0
var mu sync.Mutex

func nextID() int {
	mu.Lock()
	defer mu.Unlock()

	count++
	return count
}

func BenchmarkUserAdd(b *testing.B) {
	ctx := context.Background()

	for _, num := range []int{1, 2, 4, 8, 16, 32} {
		b.Run(fmt.Sprintf("CPU%d", num), func(b *testing.B) {
			c := client.New("http://172.19.0.5", http.DefaultClient)

			us, err := c.ListUsers(ctx)
			if err != nil {
				b.Fatal(err)
			}
			for _, u := range us {
				err := c.DeleteUser(ctx, u.Login)
				if err != nil {
					b.Fatal(err)
				}
			}

			b.SetParallelism(num)
			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					login := fmt.Sprintf("user%d", nextID())
					name := strings.ToUpper(login)
					_, err := c.AddUser(ctx, login, name)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}

}

/*
func TestUserAdd(t *testing.T) {
	ctx := context.Background()
	c := client.New("http://172.18.0.5", http.DefaultClient)

	for _, login := range []string{"alice", "bobbob", "carol"} {
		name := strings.ToUpper(login)
		u, err := c.AddUser(ctx, login, name)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("Added user", u)
	}
	us, err := c.ListUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("current users", us)

	for _, login := range []string{"alice", "bobbob", "carol"} {
		photo := strings.NewReader("raw:xxxxxxxxxxxxxx_" + login)
		err := c.UpdateUserPhoto(ctx, login, photo)
		if err != nil {
			t.Fatal(err)
		}
	}

	for _, login := range []string{"alice", "bobbob", "carol"} {
		var b bytes.Buffer
		err := c.GetUserPhoto(ctx, &b, login)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("user's photo", b.String())
	}

	for _, login := range []string{"alice", "bobbob", "carol"} {
		photo := strings.NewReader("raw:yyyyyyyyyyyyyy_" + login)
		err := c.UpdateUserPhoto(ctx, login, photo)
		if err != nil {
			t.Fatal(err)
		}
	}

	for _, login := range []string{"alice", "bobbob", "carol"} {
		var b bytes.Buffer
		err := c.GetUserPhoto(ctx, &b, login)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("user's photo", b.String())
	}

}
*/
