package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"testing"
)

func newDB() (*sql.DB, error) {
	address := "127.0.0.1"
	port := 3306
	user := "root"
	password := ""
	database := "test"

	if val := os.Getenv("MYSQL_ADDR"); len(val) > 0 {
		address = val
	}
	if val := os.Getenv("MYSQL_PORT"); len(val) > 0 {
		var err error
		port, err = strconv.Atoi(val)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	if val := os.Getenv("MYSQL_USER"); len(val) > 0 {
		user = val
	}
	if val := os.Getenv("MYSQL_PASSWORD"); len(val) > 0 {
		password = val
	}
	if val := os.Getenv("MYSQL_DATABASE"); len(val) > 0 {
		database = val
	}

	src := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		user, password, address, port, database)
	return sql.Open("mysql", src)
}

func initDB() error {
	db, err := newDB()
	if err != nil {
		return err
	}
	_, err = db.Exec(`TRUNCATE user`)
	return err
}

func TestMain(m *testing.M) {
	err := initDB()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}
