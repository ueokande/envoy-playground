package minio

import (
	"context"
	"io/ioutil"
	"strings"
	"testing"

	minio "github.com/minio/minio-go"
	"github.com/ueokande/envoy-playground/blob"
)

func TestGet(t *testing.T) {
	ctx := context.Background()

	m, err := initMinit(t)
	if err != nil {
		t.Fatal(err)
	}

	bucket := strings.ToLower(t.Name())
	src := "hello-wonderland"
	_, err = m.PutObject(bucket, "first", strings.NewReader(src), int64(len(src)), minio.PutObjectOptions{})
	if err != nil {
		t.Fatal(err)
	}

	c := impl{c: m, bucket: bucket}
	o, err := c.Get(ctx, "first")
	if err != nil {
		t.Fatal(err)
	}
	defer o.Close()

	b, err := ioutil.ReadAll(o)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != src {
		t.Error("string(b) != src: ", string(b))
	}

	o, err = c.Get(ctx, "gone")
	if err != blob.ErrNotFound {
		t.Error("err != blob.ErrNotFound:", err)
	}
}

func TestPut(t *testing.T) {
	ctx := context.Background()

	m, err := initMinit(t)
	if err != nil {
		t.Fatal(err)
	}

	bucket := strings.ToLower(t.Name())
	c := impl{c: m, bucket: bucket}

	for name, blob := range map[string]string{"first": "alice", "second": "bob"} {
		err := c.Put(ctx, name, strings.NewReader(blob), int64(len(blob)))
		if err != nil {
			t.Fatal(err)
		}
	}

	o1, err := m.GetObjectWithContext(ctx, bucket, "first", minio.GetObjectOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer o1.Close()

	b, err := ioutil.ReadAll(o1)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "alice" {
		t.Error(`string(b) != "alice": `, string(b))
	}

	o2, err := m.GetObjectWithContext(ctx, bucket, "second", minio.GetObjectOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer o2.Close()

	b, err = ioutil.ReadAll(o2)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "bob" {
		t.Error(`string(b) != "bob": `, string(b))
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()

	m, err := initMinit(t)
	if err != nil {
		t.Fatal(err)
	}

	bucket := strings.ToLower(t.Name())
	src := "hello-wonderland"
	_, err = m.PutObject(bucket, "first", strings.NewReader(src), int64(len(src)), minio.PutObjectOptions{})
	if err != nil {
		t.Fatal(err)
	}

	c := impl{c: m, bucket: bucket}
	err = c.Delete(ctx, "first")
	if err != nil {
		t.Fatal(err)
	}

	err = c.Delete(ctx, "gone")
	if err != nil {
		t.Fatal(err)
	}
}
