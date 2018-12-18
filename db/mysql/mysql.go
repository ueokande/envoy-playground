package mysql

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/ueokande/envoy-playground/db"
)

func New(db *sql.DB) db.Interface {
	return &impl{db: db}
}

type impl struct {
	db *sql.DB
}

func (i *impl) Close() error {
	return i.db.Close()
}

func isNotFound(err error) bool {
	return err == sql.ErrNoRows
}

func isConflict(err error) bool {
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}
	return mysqlErr.Number == 1062
}
