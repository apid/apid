package data_test

import (
	"fmt"
	"github.com/30x/apid"
	"github.com/30x/apid/factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
	"github.com/30x/apid/data"
	"database/sql"
)

const (
	count    = 2000
	setupSql = `
		CREATE TABLE test_1 (id INTEGER PRIMARY KEY, counter TEXT);
		CREATE TABLE test_2 (id INTEGER PRIMARY KEY, counter TEXT);`
)

var (
	tmpDir string
	r      *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

var _ = Describe("Data Service", func() {

	BeforeSuite(func() {
		apid.Initialize(factory.DefaultServicesFactory())

		var err error
		config := apid.Config()
		tmpDir, err = ioutil.TempDir("", "apid_test")
		Expect(err).NotTo(HaveOccurred())
		config.Set("local_storage_path", tmpDir)
	})

	AfterSuite(func() {
		os.RemoveAll(tmpDir)
	})

	It("should not allow reserved id or version", func() {
		_, err := apid.Data().DBForID("common")
		Expect(err).To(HaveOccurred())

		_, err = apid.Data().DBVersion("base")
		Expect(err).To(HaveOccurred())

		_, err = apid.Data().DBVersionForID("common", "base")
		Expect(err).To(HaveOccurred())
	})

	It("should be able to change versions of a datbase", func() {
		var versions []string
		var dbs []apid.DB

		for i := 0; i < 2; i++ {
			version := time.Now().String()
			db, err := apid.Data().DBVersionForID("test", version)
			Expect(err).NotTo(HaveOccurred())
			setup(db)
			versions = append(versions, version)
			dbs = append(dbs, db)
		}

		for _, db := range dbs {
			var numRows int
			err := db.QueryRow(`SELECT count(*) FROM test_2`).Scan(&numRows)
			Expect(err).NotTo(HaveOccurred())
			Expect(numRows).To(Equal(count))
		}
	})

	It("should be able to release a database", func() {
		db, err := apid.Data().DBVersionForID("release", "version")
		Expect(err).NotTo(HaveOccurred())
		setup(db)
		id := data.VersionedDBID("release", "version")
		sqlDB := db.(*sql.DB)
		Expect(sqlDB.Stats().OpenConnections).To(Equal(1))
		// run finalizer
		data.Delete(id).(func(db *sql.DB))(sqlDB)
		Expect(sqlDB.Stats().OpenConnections).To(Equal(0))
		Expect(data.DBPath(id)).ShouldNot(BeAnExistingFile())
	})

	It("should handle multi-threaded access", func(done Done) {
		db, err := apid.Data().DBForID("test")
		Expect(err).NotTo(HaveOccurred())
		setup(db)
		finished := make(chan struct{})

		go func() {
			for i := 0; i < count; i++ {
				write(db, i)
			}
			finished <- struct{}{}
		}()

		go func() {
			for i := 0; i < count; i++ {
				go func() {
					read(db)
					finished <- struct{}{}
				}()
				time.Sleep(time.Duration(r.Intn(2)) * time.Millisecond)
			}
		}()

		for i := 0; i < count+1; i++ {
			<-finished
		}

		close(done)
	}, 10)
})

func setup(db apid.DB) {
	_, err := db.Exec(setupSql)
	if err != nil {
		log.Fatal(err)
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < count; i++ {
		_, err := tx.Exec("INSERT INTO test_2 (counter) VALUES (?);", strconv.Itoa(i))
		if err != nil {
			log.Fatalf("filling up test_2 table failed. Exec error=%s", err)
		}
	}
	tx.Commit()
}

func read(db apid.DB) {
	var counter string
	rows, err := db.Query(`SELECT counter FROM test_2 LIMIT 5`)
	if err != nil {
		log.Fatalf("test_2 select failed. Exec error=%s", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&counter)
			//fmt.Print("*")
		}
	}
	fmt.Print(".")
}

func write(db apid.DB, i int) {
	tx, err := db.Begin()
	defer tx.Rollback()
	if err != nil {
		log.Fatalf("Write failed. Exec error=%s", err)
	}
	prep, err := tx.Prepare("INSERT INTO test_1 (counter) VALUES ($1);")
	_, err = tx.Stmt(prep).Exec(strconv.Itoa(i))
	if err != nil {
		log.Fatalf("Write failed. Exec error=%s", err)
	}
	prep.Close()
	tx.Commit()
	fmt.Print("+")
}
