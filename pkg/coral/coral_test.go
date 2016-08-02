package coral_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/buth/pillar/pkg/model"
	"github.com/coralproject/sponge/pkg/coral"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When sending data to the coral system", func() {

	var (
		server *httptest.Server
	)

	BeforeEach(func() {
		setupEnvironment()

		// get the test strategy file
		setupStrategyConfEnv("strategy_coral_test.json")

		// mock server
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var err error

			switch r.RequestURI {
			case "/api/import/user":

				user := model.User{}
				err = json.NewDecoder(r.Body).Decode(&user)
			case "/api/import/asset":

				asset := model.Asset{}
				err = json.NewDecoder(r.Body).Decode(&asset)
			case "/api/import/comment":

				comment := model.Comment{}
				err = json.NewDecoder(r.Body).Decode(&comment)
			case "/api/import/index":

				index := model.Index{}
				err = json.NewDecoder(r.Body).Decode(&index)
			default:
				err = fmt.Errorf("Bad request") //errors.New("Bad request")
			}

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			if err == nil {
				w.WriteHeader(http.StatusOK)
			}
			w.Header().Set("Content-Type", "application/json")

			fmt.Fprintln(w, err)
		}))
		os.Setenv("PILLAR_URL", server.URL)

		fmt.Println("DEBUG server url ", server.URL)
		coral.Init(uuid)
	})

	AfterEach(func() {
		recoverEnvironment()
	})

	Context("when using wrong collection or table name", func() {
		var err error

		BeforeEach(func() {
			var data map[string]interface{}
			tableName := "wrongTable"

			err = coral.AddRow(data, tableName)
		})

		It("should return an error", func() {
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("when adding", func() {
		Context("a user", func() {

			var err error

			JustBeforeEach(func() {
				var newrow map[string]interface{}

				newrow, err = GetFixture("user.json")
				Expect(err).Should(BeNil(), "when loading fixtures")

				err = coral.AddRow(newrow, "users")
			})
			It("should not return error", func() {
				Expect(err).Should(BeNil())
			})
		})

		// Context("an asset", func() {
		// 	var err error
		// 	JustBeforeEach(func() {
		// 		newrow, e := GetFixture("assets.json")
		// 		Expect(e).Should(BeNil())
		// 		tableName := "assets"
		//
		// 		err = coral.AddRow(newrow, tableName)
		// 	})
		// 	It("should not return error", func() {
		//
		// 		Expect(err).Should(BeNil())
		// 	})
		// })

		// Context("comment", func() {
		// 	var err error
		// 	JustBeforeEach(func() {
		// 		newrow, e := GetFixture("comments.json")
		// 		Expect(e).Should(BeNil())
		//
		// 		tableName := "comments"
		//
		// 		err = coral.AddRow(newrow, tableName)
		// 	})
		// 	It("should not return an error", func() {
		// 		Expect(err).Should(BeNil())
		// 	})
		// })

		// Context("an index", func() {
		// 	var err error
		// 	JustBeforeEach(func() {
		// 		tableName := "comments"
		// 		err = coral.CreateIndex(tableName)
		// 	})
		// 	It("should not return an error", func() {
		// 		Expect(err).Should(BeNil())
		// 	})
		// })

		// Context("an index for a table that does not exist", func() {
		// 	var err error
		// 	JustBeforeEach(func() {
		// 		tableName := "itdoesnotexist"
		// 		err = coral.CreateIndex(tableName)
		// 	})
		// 	It("should return an error", func() {
		// 		Expect(err).Should(BeNil())
		// 	})
		// })
	})
})
