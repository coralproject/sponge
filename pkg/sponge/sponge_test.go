package sponge_test

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/sponge/pkg/coral"
	"github.com/coralproject/sponge/pkg/fiddler"
	uuidimported "github.com/pborman/uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	. "github.com/coralproject/sponge/pkg/sponge"
)

var _ = Describe("Transform row of data", func() {

	var (
		oStrategy        string
		oPillarURL       string
		oPollingInterval string
	)

	BeforeEach(func() {

		// Save original enviroment variables
		oStrategy = os.Getenv("STRATEGY_CONF")
		oPillarURL = os.Getenv("PILLAR_URL")
		oPollingInterval = os.Getenv("POLLING_INTERVAL")

		logLevel := func() int {
			ll, err := cfg.Int("LOGGING_LEVEL")
			if err != nil {
				return log.NONE
			}
			return ll
		}

		log.Init(os.Stderr, logLevel)

		// MOCK STRATEGY CONF
		strategyConf := os.Getenv("GOPATH") + "src/github.com/coralproject/sponge/tests/strategy_sponge_test.json"
		e := os.Setenv("STRATEGY_CONF", strategyConf)
		if e != nil {
			fmt.Println("It could not setup the mock strategy conf variable")
		}

		// create and seed mysql foreign test database (coral_test , gaba, gabita )  - TO DO - get the credentials into a separate file
		create_and_seed_test_db()

		u := uuidimported.New()

		fiddler.Init(u)
		coral.Init(u)
	})

	AfterEach(func() {
		// recover the environment variables
		e := os.Setenv("STRATEGY_CONF", oStrategy)
		if e != nil {
			fmt.Println("It could not setup the mock strategy conf variable")
		}

		e = os.Setenv("PILLAR_URL", oPillarURL)
		if e != nil {
			fmt.Println("It could not setup the mock pillar url variable")
		}

		e = os.Setenv("POLLING_INTERVAL", oPollingInterval)
		if e != nil {
			fmt.Println("It could not setup the mock polling interval variable")
		}

		// tear down mysql foreign test database
	})

	Describe("import data", func() {
		Context("with a default option values", func() {
			It("should be valid", func() {

				limit := 999
				offset := 0
				orderby := ""
				query := ""
				types := ""
				importonlyfailed := false
				reportOnFailedRecords := false
				reportdbfile := ""
				timeWaiting := 5

				AddOptions(limit, offset, orderby, query, types, importonlyfailed, reportOnFailedRecords, reportdbfile, timeWaiting)

				Import()

				//	Expect().To(Equal())
			})

		})
	})

	Describe("create index", func() {
		Context("with a default option values", func() {
			It("should be valid", func() {

				limit := 999
				offset := 0
				orderby := ""
				query := ""
				types := ""
				importonlyfailed := false
				reportOnFailedRecords := false
				reportdbfile := ""
				timeWaiting := 5

				AddOptions(limit, offset, orderby, query, types, importonlyfailed, reportOnFailedRecords, reportdbfile, timeWaiting)

				Import()

				//	Expect().To(Equal())
			})
		})
	})

})

func create_and_seed_test_db() {
	db, err := sql.Open("mysql", "gaba:gabita@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE coral_test")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE coral_test")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE example ( id integer, data varchar(32) )")
	if err != nil {
		panic(err)
	}
}
