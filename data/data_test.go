package data_test

import (
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
	"fmt"
)

const (
	count    = 2000
	setupSql = `
		CREATE TABLE IF NOT EXISTS test_1 (id INTEGER PRIMARY KEY, counter TEXT);
		CREATE TABLE IF NOT EXISTS test_2 (id INTEGER PRIMARY KEY, counter TEXT);
		DELETE FROM test_1;
		DELETE FROM test_2;
	`
)

var (
	tmpDir string
	r      *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

var _ = Describe("Events Service", func() {

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

	It("should be able to open a new datbase", func () {
		db, err := apid.Data().DBForID("test")
		Expect(err).NotTo(HaveOccurred())
		setup(db)

		var prod string
		rows, err := db.Query(`SELECT counter FROM test_2 LIMIT 5`)
		Expect(err).NotTo(HaveOccurred())
		defer rows.Close()
		var count = 0
		for rows.Next() {
			count++
			rows.Scan(&prod)
		}
		Expect(count).To(Equal(5))

		//db, err := apid.Data().DBForID("test", "someid")
	})

	It("should handle multi-threaded access", func(done Done) {
		db, err := apid.Data().DBForID("test")
		Expect(err).NotTo(HaveOccurred())
		setup(db)

		finished := make(chan struct{})

		go func() {
			for i := 0; i < count; i++ {
				write(db, i)
				randomSleep()
			}
			finished <- struct{}{}
		}()

		go func() {
			for i := 0; i < count; i++ {
				go func(i int) {
					read(db, i)
					finished <- struct{}{}
				}(i)
				randomSleep()
			}
		}()

		for i := 0; i < count+1; i++ {
			<-finished
		}

		close(done)
	}, 4)
})

func randomSleep() {
	time.Sleep(time.Duration(r.Intn(1)) * time.Millisecond)
}

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

func read(db apid.DB, i int) {
	var prod string
	rows, err := db.Query(`SELECT counter FROM test_2 LIMIT 5`)
	if err != nil {
		log.Fatalf("test_2 select failed. Exec error=%s", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&prod)
			fmt.Print("*")
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
