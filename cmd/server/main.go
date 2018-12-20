package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/cybozu-go/log"
	"github.com/cybozu-go/well"
	"github.com/ueokande/envoy-playground/blob/minio"
	"github.com/ueokande/envoy-playground/db/mysql"
	"github.com/ueokande/envoy-playground/web"
)

var (
	flgHTTP = flag.String("http", ":8080", "Listen port and address")

	dbAddr     = os.Getenv("MYSQL_ADDR")
	dbPortStr  = os.Getenv("MYSQL_PORT")
	dbName     = os.Getenv("MYSQL_NAME")
	dbUser     = os.Getenv("MYSQL_USER")
	dbPassword = os.Getenv("MYSQL_PASSWORD")
	dbPort     int

	minioEndpoint  = os.Getenv("MINIO_ENDPOINT")
	minioAccessKey = os.Getenv("MINIO_ACCESS_KEY")
	minioSecretKey = os.Getenv("MINIO_SECRET_KEY")
	minioBucket    = os.Getenv("MINIO_BUCKET")
)

func validate() error {
	var err error
	if len(dbAddr) == 0 {
		return errors.New("MYSQL_ADDR not set")
	}
	dbPort, err = strconv.Atoi(dbPortStr)
	if err != nil {
		return errors.New("parsing MYSQL_PORT: " + err.Error())
	}
	if len(dbName) == 0 {
		return errors.New("MYSQL_NAME not set")
	}
	if len(dbUser) == 0 {
		return errors.New("MYSQL_USER not set")
	}
	if len(dbPassword) == 0 {
		return errors.New("MYSQL_PASSWORD not set")
	}
	if len(minioEndpoint) == 0 {
		return errors.New("MINIO_ENDPOINT not set")
	}
	if len(minioAccessKey) == 0 {
		return errors.New("MINIO_ACCESS_KEY not set")
	}
	if len(minioSecretKey) == 0 {
		return errors.New("MINIO_SECRET_KEY not set")
	}
	if len(minioBucket) == 0 {
		return errors.New("MINIO_BUCKET not set")
	}
	return nil
}

func run() error {
	err := validate()
	if err != nil {
		return err
	}

	src := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", dbUser, dbPassword, dbAddr, dbPort, dbName)
	db, err := sql.Open("mysql", src)
	if err != nil {
		return err
	}
	defer db.Close()

	d := mysql.New(db)
	b, err := minio.New(minio.Conf{
		Endpoint:  minioEndpoint,
		AccessKey: minioAccessKey,
		SecretKey: minioSecretKey,
		UseSSL:    false,
	}, minioBucket)
	if err != nil {
		return err
	}
	h := web.New(d, b)

	log.Info("starting server", map[string]interface{}{
		"http":           *flgHTTP,
		"mysql_addr":     dbAddr,
		"mysql_port":     dbPort,
		"mysql_name":     dbName,
		"minio_endpoint": minioEndpoint,
		"minio_bucket":   minioBucket,
	})

	logger := log.NewLogger()

	serv := &well.HTTPServer{
		Server: &http.Server{
			Addr:    *flgHTTP,
			Handler: h,
		},
		AccessLog: logger,
	}

	err = serv.ListenAndServe()
	if err != nil {
		return err
	}
	return well.Wait()
}

func main() {
	flag.Parse()
	well.LogConfig{}.Apply()

	err := run()
	if err != nil {
		log.ErrorExit(err)
	}
}
