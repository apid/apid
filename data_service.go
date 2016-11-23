package apid

import (
	"database/sql"
)

type DataService interface {
	DB() (DB, error)
	DBForID(id string) (db DB, err error)
}

type DB interface {
	Ping() (error)
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Begin() (*sql.Tx, error)

	//Close() error
	//Stats() sql.DBStats
	//Driver() driver.Driver
}
