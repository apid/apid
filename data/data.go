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

	log.Errorf("#### DB path %s", datapath)
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

func (d *dataService) InsertSnapshotDB(rawjson []byte) (err error) {

	log.Info("Cleanup old DB")
	dbname := path.Join(config.GetString(configDataPath), commonDBID)

	/* Clean up existing DB before getting new one */
	os.Remove(dbname)
	os.Remove(dbname + "-wal")
	os.Remove(dbname + "-shm")

	/* Create DB pointer */
	db, err := d.DB()
	if err != nil {
		log.Error("Unable to create DB")
		return err
	}

	snp, err := common.UnmarshalSnapshot(rawjson)
	if err != nil {
		log.Erro("Unable to Marshal Snapshot")
		return err
	}
	txn, err := db.Begin()
	if err != nil {
		log.Error("Unable to get Txn for DB")
		return err
	}

	for _, table := range snp.Tables {

		for _, row := range table.Rows {
			count := 0
			sqli := "INSERT INTO " + table.Name + " ("
			for rn, _ := range row {
				count++
				sqli += rn
				if count < len(row) {
					sqli += ","
				}
			}
			count = 0
			sqli += ") VALUES ("
			for _, rv := range row {
				count++
				switch rv.Type {
				case 1043:
					sqli += "'" + rv.Value + "'"
				default:
					sqli += rv.Value
				}
				if count < len(row) {
					sqli += ","
				}
			}
			sqli += ")"
			_, err = txn.Exec(sqli)
		}
	}
	txn.Commit()

	log.Info("Downloaded a new DB from Snapshot server")
	return err
}
