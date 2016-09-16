package data

import (
	"database/sql"
	"fmt"
	"github.com/30x/apid"
	_ "github.com/mattn/go-sqlite3"
	"path"
	"sync"
)

const (
	configDataDriver = "data_driver"
	configDataSource = "data_source"
	configDataPath   = "data_path"

	commonDBID = "_apid_common_"
)

var log apid.LogService
var config apid.ConfigService

var dbMap = make(map[string]*sql.DB)
var dbMapSync sync.RWMutex

func CreateDataService() apid.DataService {
	config = apid.Config()
	log = apid.Log().ForModule("data")

	config.SetDefault(configDataDriver, "sqlite3")
	config.SetDefault(configDataSource, "file:%s?cache=shared&mode=rwc")
	config.SetDefault(configDataPath, "/var/tmp")

	return &dataService{}
}

type dataService struct {
}

func (d *dataService) DB() (*sql.DB, error) {
	return d.DBForID(commonDBID)
}

func (d *dataService) DBForID(id string) (db *sql.DB, err error) {

	dbMapSync.RLock()
	db = dbMap[id]
	dbMapSync.RUnlock()
	if db != nil {
		return
	}

	dbMapSync.Lock()
	defer dbMapSync.Unlock()

	db = dbMap[id]
	if db != nil {
		dbMapSync.Unlock()
		return
	}

	log.Info("LoadDB: ", id)

	dataPath := path.Join(config.GetString(configDataPath), id)
	dataSource := fmt.Sprintf(config.GetString(configDataSource), dataPath)
	db, err = sql.Open(config.GetString(configDataDriver), dataSource)
	if err != nil {
		log.Errorf("error loading db: %s", err)
		return
	}

	err = db.Ping()
	if err != nil {
		log.Errorf("error pinging db: %s", err)
		return
	}

	sqlString := "PRAGMA journal_mode=WAL;"
	_, err = db.Exec(sqlString)
	if err != nil {
		log.Errorf("error setting journal_mode: %s", err)
		return
	}

	sqlString = "PRAGMA foreign_keys = ON;"
	_, err = db.Exec(sqlString)
	if err != nil {
		log.Errorf("error enabling foreign_keys: %s", err)
		return
	}

	dbMap[id] = db
	return
}
