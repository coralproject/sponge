package source_test

import (
	"os"

	"github.com/coralproject/sponge/pkg/source"
	"github.com/pborman/uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//***************** MONGODB *****************//

var _ = Describe("Source MongoDB", func() {

	var (
		mdb source.MongoDB
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
		mdb, ok = s.(source.MongoDB)
		if !ok {
			Expect(ok).To(BeTrue(), "The source's implementation should be MongoDB")
		}
	})

	Describe("Creating New MongoDB Source", func() {
		Context("with a valid connection", func() {
			It("should not have a nil connection", func() {
				Expect(&mdb.Connection).NotTo(BeNil())
			})
			It("should not have a nil database", func() {
				Expect(&mdb.Database).NotTo(BeNil())
			})
		})
	})
})

//***************** MYSQL *****************//

var _ = Describe("Source MySQL", func() {

	var (
		mdb source.MySQL
	)

	BeforeEach(func() {

		// MOCK STRATEGY CONF WITH MYSQL CONFIG FILE
		strategyConf := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/tests/strategy_mysql_test.json"
		if e := os.Setenv("STRATEGY_CONF", strategyConf); e != nil {
			Expect(e).NotTo(HaveOccurred())
		}

		u := uuid.New()
		if _, e := source.Init(u); e != nil {
			Expect(e).NotTo(HaveOccurred())
		}

		s, e := source.New("mysql")
		if e != nil {
			Expect(e).NotTo(HaveOccurred())
		}
		var ok bool
		mdb, ok = s.(source.MySQL)
		if !ok {
			Expect(ok).To(BeTrue(), "The source's implementation should be MySQL")
		}
	})

	Describe("Creating New MySQL Source", func() {
		Context("with a valid connection", func() {
			It("should not have a nil connection", func() {
				Expect(&mdb.Connection).NotTo(BeNil())
			})
			It("should not have a nil database", func() {
				Expect(&mdb.Database).NotTo(BeNil())
			})
		})
	})
})

//***************** PostgreSQL *****************//

var _ = Describe("Source PostgreSQL", func() {

	var (
		mdb source.PostgreSQL
	)

	BeforeEach(func() {

		// MOCK STRATEGY CONF WITH POSTGRESQL CONFIG FILE
		strategyConf := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/tests/strategy_postgresql_test.json"
		if e := os.Setenv("STRATEGY_CONF", strategyConf); e != nil {
			Expect(e).NotTo(HaveOccurred())
		}

		u := uuid.New()
		if _, e := source.Init(u); e != nil {
			Expect(e).NotTo(HaveOccurred())
		}

		s, e := source.New("postgresql")
		if e != nil {
			Expect(e).NotTo(HaveOccurred())
		}
		var ok bool
		mdb, ok = s.(source.PostgreSQL)
		if !ok {
			Expect(ok).To(BeTrue(), "The source's implementation should be PostgreSQL")
		}
	})

	Describe("Creating New PostgreSQL Source", func() {
		Context("with a valid connection", func() {
			It("should not have a nil connection", func() {
				Expect(&mdb.Connection).NotTo(BeNil())
			})
			It("should not have a nil database", func() {
				Expect(&mdb.Database).NotTo(BeNil())
			})
		})
	})
})

//***************** Service *****************//

var _ = Describe("Source Service", func() {

	var (
		mdb source.API
	)

	BeforeEach(func() {

		// MOCK STRATEGY CONF WITH SERVICE CONFIG FILE
		strategyConf := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/tests/strategy_api_test.json"
		if e := os.Setenv("STRATEGY_CONF", strategyConf); e != nil {
			Expect(e).NotTo(HaveOccurred())
		}

		u := uuid.New()
		if _, e := source.Init(u); e != nil {
			Expect(e).NotTo(HaveOccurred())
		}

		s, e := source.New("service")
		if e != nil {
			Expect(e).NotTo(HaveOccurred())
		}
		var ok bool
		mdb, ok = s.(source.API)
		if !ok {
			Expect(ok).To(BeTrue(), "The source's implementation should be Service")
		}
	})

	Describe("Creating New Service Source", func() {
		Context("with a valid connection", func() {
			It("should not have a nil connection", func() {
				Expect(&mdb.Connection).NotTo(BeNil())
			})
		})
	})
})
