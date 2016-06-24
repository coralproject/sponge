package source_test

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"gopkg.in/mgo.v2"

	"github.com/coralproject/sponge/pkg/source"
	"github.com/pborman/uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const MONGOCONNECTION = "mongodb://localhost:27017"

// CREATE MongoDB TEST database
func createTestMongoDB() {

	testMongoDB := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/tests/fixtures/mongo_comments.json"
	file, e := ioutil.ReadFile(testMongoDB)
	if e != nil {
		Expect(e).NotTo(HaveOccurred())
	}
	var comments []map[string]interface{}
	if e := json.Unmarshal(file, &comments); e != nil {
		Expect(e).NotTo(HaveOccurred())
	}

	session, e := mgo.Dial(MONGOCONNECTION)
	if e != nil {
		Expect(e).NotTo(HaveOccurred(), MONGOCONNECTION)
	}
	defer session.Close()

	col := session.DB(DATABASE).C("comments")
	for _, c := range comments {
		if e := col.Insert(c); e != nil {
			Expect(e).NotTo(HaveOccurred())
		}
	}
}

func removeMongoDB() {

	session, e := mgo.Dial(MONGOCONNECTION)
	if e != nil {
		Expect(e).NotTo(HaveOccurred())
	}
	defer session.Close()

	db := session.DB(DATABASE)
	if db.DropDatabase(); e != nil {
		Expect(e).NotTo(HaveOccurred())
	}
}

var _ = Describe("Getting Data", func() {

	var (
		mdb  source.MongoDB
		data []map[string]interface{}
	)

	BeforeEach(func() {

		// MOCK STRATEGY CONF WITH MONGODB CONFIG FILE
		strategyConf := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/tests/strategy_mongo_test.json"
		if e := os.Setenv("STRATEGY_CONF", strategyConf); e != nil {
			Expect(e).NotTo(HaveOccurred())
		}

		u := uuid.New()
		if _, e := source.Init(u); e != nil {
			Expect(e).NotTo(HaveOccurred())
		}

		s, e := source.New("mongodb")
		if e != nil {
			Expect(e).NotTo(HaveOccurred())
		}
		var ok bool
		if mdb, ok = s.(source.MongoDB); !ok {
			Expect(ok).To(BeTrue(), "The source's implementation should be MongoDB")
		}

		// Default Flags
		coralName := "comments"
		options := &source.Options{
			Offset:  0,
			Limit:   9999999999,
			Orderby: "",
			Query:   "",
		}

		if data, e = mdb.GetData(coralName, options); e != nil {
			Expect(e).NotTo(HaveOccurred())
		}
	})

	Describe("from external Mongo database", func() {
		Context("with a valid strategy file", func() {
			It("should get back the count of records we expect", func() {
				Expect(len(data)).To(Equal(10))
			})
			It("should be data we expect related to the comment", func() {
				Expect(data[0]["object.content"]).To(Equal("Comment1"))
			})
			It("should be data we expect related to the comment", func() {
				Expect(data[0]["actor.title"]).To(Equal("Socks_Friend"))
			})
		})
	})
})
