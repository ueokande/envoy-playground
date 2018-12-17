package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ueokande/envoy-playground/db/mysql"
	"github.com/ueokande/envoy-playground/web"
)

var (
	flgHTTP = flag.String("http", ":8080", "Listen port and address")

	flgMySQLAddr     = flag.String("mysql-addr", "127.0.0.1", "MySQL address")
	flgMySQLPort     = flag.Int("mysql-port", 3306, "MySQL port")
	flgMySQLDatabase = flag.String("mysql-database", "envoy-playground", "MySQL database name")

	mySQLUser     = os.Getenv("MYSQL_USER")
	mySQLPassword = os.Getenv("MYSQL_PASSWORD")
)

func run() error {
	if len(mySQLUser) == 0 {
		return errors.New("MYSQL_USER not set")
	}
	if len(mySQLPassword) == 0 {
		return errors.New("MYSQL_PASSWORD not set")
	}

	var conf = mysql.Conf{
		Address:  *flgMySQLAddr,
		Port:     *flgMySQLPort,
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Database: *flgMySQLDatabase,
	}

	if val := os.Getenv("MYSQL_USER"); len(val) > 0 {
		conf.User = val
	}
	if val := os.Getenv("MYSQL_PASSWORD"); len(val) > 0 {
		conf.Password = val
	}

	d, err := mysql.New(conf)
	if err != nil {
		return err
	}

	h := web.New(d)

	log.Printf("server started, http=%s, mysql-addr=%s, mysql-port=%d, mysql-database=%s",
		*flgHTTP, conf.Address, conf.Port, conf.Database)

	http.Handle("/", h)
	return http.ListenAndServe(":8080", nil)
}

func main() {
	flag.Parse()

	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
