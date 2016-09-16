package apid

import "database/sql"

type DataService interface {
	DB() (*sql.DB, error)
	DBForID(id string) (db *sql.DB, err error)
}
