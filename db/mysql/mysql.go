package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ueokande/envoy-playground/db"
)

type Conf struct {
	Address  string
	Port     int
	User     string
	Password string
	Database string
}

func New(c Conf) (db.Interface, error) {
	src := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", c.User, c.Password, c.Address, c.Port, c.Database)
	db, err := sql.Open("mysql", src)
	if err != nil {
		return nil, err
	}
	return &impl{db: db}, nil
}

type impl struct {
	db *sql.DB
}

func (i *impl) Close() error {
	return i.db.Close()
}
