package minio

import (
	"os"
	"strings"
	"testing"

	minio "github.com/minio/minio-go"
)

func initMinit(t *testing.T) (*minio.Client, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	c, err := minio.New(endpoint, accessKey, secretKey, false)
	if err != nil {
		return nil, err
	}

	bucket := strings.ToLower(t.Name())

	exists, err := c.BucketExists(bucket)
	if err != nil {
		t.Log("1")
		return nil, err
	}
	if !exists {
		err = c.MakeBucket(bucket, "")
		if err != nil {
			return nil, err
		}
	}

	done := make(chan struct{})
	defer close(done)
	for o := range c.ListObjects(bucket, "", true, done) {
		c.RemoveObject(bucket, o.Key)
	}
	return c, nil
}
