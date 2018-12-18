package main

import (
	"database/sql"
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

	dbUser     = os.Getenv("MYSQL_USER")
	dbPassword = os.Getenv("MYSQL_PASSWORD")
)

func run() error {
	if len(dbUser) == 0 {
		dbUser = "root"
	}

	src := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		dbUser, dbPassword, *flgMySQLAddr, *flgMySQLPort, *flgMySQLDatabase)
	db, err := sql.Open("mysql", src)
	if err != nil {
		return err
	}

	d := mysql.New(db)
	h := web.New(d)

	log.Printf("server started, http=%s, mysql-addr=%s, mysql-port=%d, mysql-database=%s",
		*flgHTTP, *flgMySQLAddr, *flgMySQLPort, *flgMySQLDatabase)

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
