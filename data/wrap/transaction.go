package wrap

import (
	"github.com/30x/apid"
	"github.com/mattn/go-sqlite3"
)

type wrapTx struct {
	*sqlite3.SQLiteTx
	log apid.LogService
}

func (tx *wrapTx) Commit() (err error) {
	tx.log.Debug("begin commit")

	if err = tx.SQLiteTx.Commit(); err != nil {
		tx.log.Errorf("failed commit: %s", err)
		return
	}

	tx.log.Debug("end commit")
	return
}

func (tx *wrapTx) Rollback() (err error) {
	tx.log.Debug("begin rollback")

	if err = tx.SQLiteTx.Rollback(); err != nil {
		tx.log.Errorf("failed rollback: %s", err)
	}

	tx.log.Debug("end rollback")
	return
}
