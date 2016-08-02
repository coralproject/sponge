package source_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/coralproject/sponge/pkg/source"
	"github.com/pborman/uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// CREATE Server for Test
func mockServer() *httptest.Server {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				var d map[string]interface{}
				file, err := ioutil.ReadFile(os.Getenv("$GOPATH") + "/tests/response.json")
				if err != nil {
					fmt.Printf("ERROR %v on setting up response in the test.", err)
				}
				if err := json.Unmarshal(file, &d); err != nil {
					fmt.Printf("ERROR %v on setting up response in the test.", err)
				}

				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, d)
			},
		),
	)
	return server
}

var _ = Describe("Getting Data", func() {

	var (
		mdb  source.API
		data []map[string]interface{}
	)

	BeforeEach(func() {

		// MOCK STRATEGY CONF WITH API CONFIG FILE
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
		if mdb, ok = s.(source.API); !ok {
			Expect(ok).To(BeTrue(), "The source's implementation should be Service")
		}

		// attributes := "scope:https://www.washingtonpost.com/lifestyle/style/carolyn-hax-stubborn-60-something-parent-refuses-to-see-a-doctor/2015/09/24/299ec776-5e2d-11e5-9757-e49273f05f65_story.html source:washpost.com itemsPerPage:100 sortOrder:reverseChronological"
		// connection := fmt.Sprintf("%s/v1/search?q=((%s))&appkey=dev.washpost.com", server.URL, attributes)
		// write down connection into strategyConf
		var d map[string]interface{}
		file, e := ioutil.ReadFile(strategyConf)
		if e != nil {
			fmt.Printf("ERROR %v on setting up the new strategy conf.", e)
		}
		if e = json.Unmarshal(file, &d); e != nil {
			fmt.Printf("ERROR %v on setting up the new strategy conf.", e)
		}
		d["Credentials"].(map[string]interface{})["service"].(map[string]interface{})["endpoint"] = server.URL
		b, e := json.Marshal(d)
		if e != nil {
			fmt.Printf("ERROR %v on setting up the new strategy conf.", e)
		}
		e = ioutil.WriteFile(strategyConf, b, 0777)
		if e != nil {
			fmt.Printf("ERROR %v on setting up the new strategy conf.", e)
		}

		data, _, e = mdb.GetFireHoseData("")
		if e != nil {
			Expect(e).NotTo(HaveOccurred())
		}
	})

	Describe("from external service", func() {
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
