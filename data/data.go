package data

import (
	"database/sql"
	"fmt"
	"github.com/30x/apid"
	"github.com/30x/apid/data/wrap"
	"github.com/mattn/go-sqlite3"
	"os"
	"path"
	"sync"
	"runtime"
)

const (
	configDataDriverKey = "data_driver"
	configDataSourceKey = "data_source"
	configDataPathKey   = "data_path"

	commonDBID      = "common"
	commonDBVersion = "base"

	defaultTraceLevel = "warn"
)

var log, dbTraceLog apid.LogService
var config apid.ConfigService

var dbMap = make(map[string]*sql.DB)
var dbMapSync sync.RWMutex

func CreateDataService() apid.DataService {
	config = apid.Config()
	log = apid.Log().ForModule("data")

	// we don't want to trace normally
	config.SetDefault("DATA_TRACE_LOG_LEVEL", defaultTraceLevel)
	dbTraceLog = apid.Log().ForModule("data_trace")

	config.SetDefault(configDataDriverKey, "sqlite3")
	config.SetDefault(configDataSourceKey, "file:%s")
	config.SetDefault(configDataPathKey, "sqlite")

	return &dataService{}
}

type dataService struct {
}

func (d *dataService) DB() (apid.DB, error) {
	return d.dbVersionForID(commonDBID, commonDBVersion)
}

func (d *dataService) DBForID(id string) (apid.DB, error) {
	if id == commonDBID {
		return nil, fmt.Errorf("reserved ID: %s", id)
	}
	return d.dbVersionForID(id, commonDBVersion)
}

func (d *dataService) DBVersion(version string) (apid.DB, error) {
	if version == commonDBVersion {
		return nil, fmt.Errorf("reserved version: %s", version)
	}
	return d.dbVersionForID(commonDBID, version)
}

func (d *dataService) DBVersionForID(id, version string) (apid.DB, error) {
	if id == commonDBID {
		return nil, fmt.Errorf("reserved ID: %s", id)
	}
	if version == commonDBVersion {
		return nil, fmt.Errorf("reserved version: %s", version)
	}
	return d.dbVersionForID(id, version)
}

// will set DB to close and delete when no more references
func (d *dataService) ReleaseDB(id, version string) {
	versionedID := VersionedDBID(id, version)

	dbMapSync.Lock()
	defer dbMapSync.Unlock()

	db := dbMap[versionedID]
	if db != nil {
		dbMap[versionedID] = nil
		log.Errorf("SETTING FINALIZER")
		finalizer := Delete(versionedID)
		runtime.SetFinalizer(db, finalizer)
	}

	return
}

func (d *dataService) dbVersionForID(id, version string) (db *sql.DB, err error) {

	versionedID := VersionedDBID(id, version)

	dbMapSync.RLock()
	db = dbMap[versionedID]
	dbMapSync.RUnlock()
	if db != nil {
		return
	}

	dbMapSync.Lock()
	defer dbMapSync.Unlock()

	db = dbMap[versionedID]
	if db != nil {
		return
	}

	dataPath := DBPath(versionedID)

	if err = os.MkdirAll(path.Dir(dataPath), 0700); err != nil {
		return
	}

	log.Infof("LoadDB: %s", dataPath)
	dataSource := fmt.Sprintf(config.GetString(configDataSourceKey), dataPath)

	wrappedDriverName := "dd:" + config.GetString(configDataDriverKey)
	dataDriver := wrap.WrapDriver{&sqlite3.SQLiteDriver{}, dbTraceLog}
	func() {
		// just ignore the "registered twice" panic
		defer func() {
			recover()
		}()
		sql.Register(wrappedDriverName, &dataDriver)
	}()

	db, err = sql.Open(wrappedDriverName, dataSource)
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

	dbMap[versionedID] = db
	return
}

func Delete(versionedID string) interface{} {
	return func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Errorf("error closing DB: %v", err)
		}
		dataDir := path.Dir(DBPath(versionedID))
		err = os.RemoveAll(dataDir)
		if err != nil {
			log.Errorf("error removing DB files: %v", err)
		}
		delete(dbMap, versionedID)
	}
}

func VersionedDBID(id, version string) string {
	return path.Join(id, version)
}

func DBPath(id string) string {
	storagePath := config.GetString("local_storage_path")
	relativeDataPath := config.GetString(configDataPathKey)
	return path.Join(storagePath, relativeDataPath, id, "sqlite3")
}