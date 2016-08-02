package source_test

// CHALLENGE
// HOW TO CREATE AND REMOVE A DATABASE FOR TESTING (or mock mysql)

// import (
// 	"database/sql"
// 	"os"
//
// 	"github.com/coralproject/sponge/pkg/source"
// 	_ "github.com/go-sql-driver/mysql"
// 	"github.com/pborman/uuid"
//
// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// )
//
// const MYSQLCONNECTION = "gaba:gabita@/coral_test"
//
// // CREATE MySQL TEST database
// func createTestMysqlDB() {
//
// 	// testMySQLDB := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/tests/fixtures/mysql_comments.sql"
// 	// file, e := ioutil.ReadFile(testMySQLDB)
// 	// if e != nil {
// 	// 	Expect(e).NotTo(HaveOccurred())
// 	// }
//
// 	db, e := sql.Open("mysql", MYSQLCONNECTION)
// 	if e != nil {
// 		Expect(e).NotTo(HaveOccurred())
// 	}
//
// 	if e = db.Ping(); e != nil {
// 		Expect(e).NotTo(HaveOccurred())
// 	}
//
// 	defer db.Close()
// }
//
// func removeMysqlDB() {
//
// 	db, e := sql.Open("mysql", MYSQLCONNECTION)
// 	if e != nil {
// 		Expect(e).NotTo(HaveOccurred())
// 	}
//
// 	if e = db.Ping(); e != nil {
// 		Expect(e).NotTo(HaveOccurred())
// 	}
//
// 	defer db.Close()
//
// 	_, e = db.Exec("DROP DATABASE IF EXISTS %s", DATABASE)
// 	if e != nil {
// 		Expect(e).NotTo(HaveOccurred())
// 	}
// }
//
// var _ = Describe("Getting Data", func() {
//
// 	var (
// 		mdb  source.MySQL
// 		data []map[string]interface{}
// 	)
//
// 	BeforeEach(func() {
//
// 		// MOCK STRATEGY CONF WITH MySQL CONFIG FILE
// 		strategyConf := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/tests/strategy_mysql_test.json"
// 		if e := os.Setenv("STRATEGY_CONF", strategyConf); e != nil {
// 			Expect(e).NotTo(HaveOccurred())
// 		}
//
// 		u := uuid.New()
// 		if _, e := source.Init(u); e != nil {
// 			Expect(e).NotTo(HaveOccurred())
// 		}
//
// 		s, e := source.New("mysql")
// 		if e != nil {
// 			Expect(e).NotTo(HaveOccurred())
// 		}
// 		var ok bool
// 		if mdb, ok = s.(source.MySQL); !ok {
// 			Expect(ok).To(BeTrue(), "The source's implementation should be MySQL")
// 		}
//
// 		// Default Flags
// 		coralName := "comments"
// 		options := &source.Options{
// 			Offset:  0,
// 			Limit:   9999999999,
// 			Orderby: "",
// 			Query:   "",
// 		}
//
// 		if data, e = mdb.GetData(coralName, options); e != nil {
// 			Expect(e).NotTo(HaveOccurred())
// 		}
// 	})
//
// 	Describe("from external MySQL database", func() {
// 		Context("with a valid strategy file", func() {
// 			It("should get back the count of records we expect", func() {
// 				Expect(len(data)).To(Equal(10))
// 			})
// 			It("should be data we expect related to the comment", func() {
// 				Expect(data[0]["object.content"]).To(Equal("Comment1"))
// 			})
// 			It("should be data we expect related to the comment", func() {
// 				Expect(data[0]["actor.title"]).To(Equal("Socks_Friend"))
// 			})
// 		})
// 	})
// })
