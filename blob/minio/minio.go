package minio

import (
	"context"
	"fmt"
	"io"

	minio "github.com/minio/minio-go"
	"github.com/ueokande/envoy-playground/blob"
)

type Conf struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
}

type impl struct {
	c      *minio.Client
	bucket string
}

func New(conf Conf, bucket string) (blob.Interface, error) {
	c, err := minio.New(conf.Endpoint, conf.AccessKey, conf.SecretKey, conf.UseSSL)
	if err != nil {
		return nil, err
	}
	return &impl{c: c, bucket: bucket}, nil
}

func (i *impl) Get(ctx context.Context, name string) (blob.Object, error) {
	o, err := i.c.GetObjectWithContext(ctx, i.bucket, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	_, err = o.Stat()
	if err != nil {
		defer o.Close()
		resp := minio.ToErrorResponse(err)
		fmt.Println("RESP", resp)
		if resp.Code == "NoSuchKey" {
			return nil, blob.ErrNotFound
		}
		return nil, err
	}
	return o, nil
}

func (i *impl) Put(ctx context.Context, name string, r io.Reader) error {
	_, err := i.c.PutObjectWithContext(ctx, i.bucket, name, r, -1, minio.PutObjectOptions{})
	return err
}

func (i *impl) Delete(ctx context.Context, name string) error {
	return i.c.RemoveObject(i.bucket, name)
}
