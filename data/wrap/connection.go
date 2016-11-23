package wrap

import (
	"github.com/30x/apid"
	"github.com/mattn/go-sqlite3"
	"database/sql/driver"
)

type wrapConn struct {
	*sqlite3.SQLiteConn
	log apid.LogService
}

func (c *wrapConn) Swap(cc *sqlite3.SQLiteConn) {
	c.SQLiteConn = cc
}

func (c *wrapConn) Prepare(query string) (driver.Stmt, error) {
	c.log.Debugf("begin prepare stmt: %s", query)

	stmt, err := c.SQLiteConn.Prepare(query)
	if err != nil {
		c.log.Errorf("prepare stmt failed: %s", err)
		return nil, err
	}

	c.log.Debug("end prepare stmt")
	s := stmt.(*sqlite3.SQLiteStmt)
	return &wrapStmt{s, c.log}, nil
}

func (c *wrapConn) Begin() (driver.Tx, error) {
	c.log.Debug("begin trans")

	tx, err := c.SQLiteConn.Begin()
	if err != nil {
		c.log.Errorf("begin trans failed: %s", err)
		return nil, err
	}

	c.log.Debug("end begin trans")
	t := tx.(*sqlite3.SQLiteTx)
	return &wrapTx{t, c.log}, nil
}

func (c *wrapConn) Close() (err error) {
	c.log.Debug("begin close")

	if err = c.SQLiteConn.Close(); err != nil {
		c.log.Errorf("close failed: %s", err)
		return
	}

	c.log.Debug("end close")
	return
}
