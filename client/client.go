package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	core "github.com/ueokande/envoy-playground"
)

func New(addr string, c *http.Client) Interface {
	return &impl{
		addr: addr,
		http: c,
	}
}

type impl struct {
	addr string
	http *http.Client
}

func (i *impl) ListUsers(ctx context.Context) ([]core.User, error) {
	var us []core.User
	err := i.requestJson(ctx, "GET", "/users", nil, &us)
	if err != nil {
		return nil, err
	}
	return us, nil
}

func (i *impl) GetUser(ctx context.Context, login string) (*core.User, error) {
	var u core.User
	err := i.requestJson(ctx, "GET", "/user/"+login, nil, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (i *impl) AddUser(ctx context.Context, login, name string) (*core.User, error) {
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(core.User{Login: login, Name: name})
	if err != nil {
		return nil, err
	}
	var ret core.User
	err = i.requestJson(ctx, "POST", "/users", &b, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (i *impl) UpdateUser(ctx context.Context, u core.User) (*core.User, error) {
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(u)
	if err != nil {
		return nil, err
	}
	var ret core.User
	err = i.requestJson(ctx, "POST", "/user/"+u.Login, &b, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (i *impl) DeleteUser(ctx context.Context, login string) error {
	var msg map[string]interface{}
	return i.requestJson(ctx, "DELETE", "/user/"+login, nil, &msg)
}

func (i *impl) UpdateUserPhoto(ctx context.Context, login string, r io.Reader) error {
	var msg map[string]interface{}
	return i.requestJson(ctx, "PUT", "/user/"+login+"/photo", r, &msg)
}

func (i *impl) GetUserPhoto(ctx context.Context, w io.Writer, login string) error {
	url := i.addr + "/user/" + login + "/photo"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	resp, err := i.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)
	return err
}

func (i *impl) DeleteUserPhoto(ctx context.Context, login string) error {
	var msg map[string]interface{}
	return i.requestJson(ctx, "DELETE", "/user/"+login+"/photo", nil, &msg)
}

func (i *impl) requestJson(ctx context.Context, method, path string, body io.Reader, data interface{}) error {
	url := i.addr + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	resp, err := i.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if 400 <= resp.StatusCode && resp.StatusCode < 600 {
		var msg map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&msg)
		if err != nil {
			return err
		}
		return fmt.Errorf("%d: %s: %s", resp.StatusCode, resp.Status, msg["reason"])
	}

	return json.NewDecoder(resp.Body).Decode(&data)
}
