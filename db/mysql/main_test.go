package mysql

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

var conf = Conf{
	Address:  "127.0.0.1",
	Port:     3306,
	User:     "root",
	Password: "",
	Database: "test",
}

func TestMain(m *testing.M) {
	if val := os.Getenv("MYSQL_ADDR"); len(val) > 0 {
		conf.Address = val
	}
	if val := os.Getenv("MYSQL_PORT"); len(val) > 0 {
		port, err := strconv.Atoi(val)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		conf.Port = port
	}
	if val := os.Getenv("MYSQL_USER"); len(val) > 0 {
		conf.User = val
	}
	if val := os.Getenv("MYSQL_PASSWORD"); len(val) > 0 {
		conf.Password = val
	}
	if val := os.Getenv("MYSQL_DATABASE"); len(val) > 0 {
		conf.Database = val
	}

	os.Exit(m.Run())
}
