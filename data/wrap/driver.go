package wrap

import (
	"database/sql/driver"
	"strings"
	"github.com/30x/apid"
	"github.com/mattn/go-sqlite3"
)

type WrapDriver struct {
	driver.Driver
	Log apid.LogService
}

func (d WrapDriver) Open(dsn string) (driver.Conn, error) {

	internalDSN := strings.TrimPrefix(dsn, "dd:")
	internalCon, err := d.Driver.Open(internalDSN)
	if err != nil {
		return nil, err
	}

	c := internalCon.(*sqlite3.SQLiteConn)
	return &wrapConn{c, d.Log}, nil
}
