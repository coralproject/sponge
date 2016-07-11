package coral

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/sponge/pkg/strategy"
	. "github.com/onsi/ginkgo"
	uuidimported "github.com/pborman/uuid"
)

var _ = Describe("Testing with Ginkgo", func() {
	It("mockup server", func() {

		method := "POST"
		urlStr := server.URL + "/api/import/users"
		row := map[string]interface{}{"juan": 3, "pepe": "what"}
		juser, err := json.Marshal(row)
		payload := bytes.NewBuffer(juser)
		request, err := http.NewRequest(method, urlStr, payload)
		if err != nil {
			GinkgoT().Fatalf("expect not error and got one %s.", err)
		}
		request.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			GinkgoT().Fatalf("expect not error and got one %s.", err)
		}
		defer response.Body.Close()
	})
	It("add row wrong table", func() {

		var data map[string]interface{}

		tableName := "wrongTable"

		e := AddRow(data, tableName)
		if e == nil {
			GinkgoT().Fatal("expecting error and got none.")
		}
	})
	It("add user row", func() {

		newrow, e := GetFixture("users.json")
		if e != nil {
			GinkgoT().Fatalf("error with the test data: %s.", e)
		}

		tableName := "users"

		e = AddRow(newrow, tableName)
		if e != nil {
			GinkgoT().Fatalf("expecting not error but got one %v.", e)
		}
	})
	It("add asset row", func() {

		newrow, e := GetFixture("assets.json")
		if e != nil {
			GinkgoT().Fatalf("error with the test data: %s.", e)
		}

		tableName := "assets"

		e = AddRow(newrow, tableName)
		if e != nil {
			GinkgoT().Fatalf("expecting not error but got one %v.", e)
		}
	})
	It("add comment row", func() {

		newrow, e := GetFixture("comments.json")
		if e != nil {
			GinkgoT().Fatalf("error with the test data: %s.", e)
		}

		tableName := "comments"

		e = AddRow(newrow, tableName)
		if e != nil {
			GinkgoT().Fatalf("expecting not error but got one %v.", e)
		}
	})
	It("create index", func() {

		tableName := "comments"

		e := CreateIndex(tableName)
		if e != nil {
			GinkgoT().Fatalf("expecting not error but got one %v.", e)
		}
	})
	It("create index error", func() {

		tableName := "itdoesnotexist"

		e := CreateIndex(tableName)
		if e == nil {
			GinkgoT().Fatalf("expecting an error but got none.")
		}
	})
})
var (
	server  *httptest.Server
	path    string
	fakeStr strategy.Strategy

	oStrategy string
	oPillar   string
)

func setup() {

	oStrategy = os.Getenv("STRATEGY_CONF")
	oPillar = os.Getenv("PILLAR_URL")

	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.NONE
		}
		return ll
	}

	log.Init(os.Stderr, logLevel)

	strategyConf := os.Getenv("GOPATH") + "/test/strategy_coral_test.json"
	e := os.Setenv("STRATEGY_CONF", strategyConf)
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}

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
			err = errors.New("Bad request")
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

	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/test/fixtures/"

	os.Setenv("PILLAR_URL", server.URL)

	u := uuidimported.New()

	Init(u)
}

func teardown() {

	e := os.Setenv("STRATEGY_CONF", oStrategy)
	if e != nil {
		fmt.Println("It could not setup the strategy conf enviroment variable back.")
	}
	e = os.Setenv("PILLAR_URL", oPillar)
	if e != nil {
		fmt.Println("It could not setup the pillar home environment variable back.")
	}
}

func TestMain(m *testing.M) {

	setup()
	code := m.Run()
	teardown()

	os.Exit(code)
}

func GetFixture(fileName string) (map[string]interface{}, error) {
	file, err := os.Open(path + fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var qs map[string]interface{}
	err = json.Unmarshal(content, &qs)
	if err != nil {
		return nil, err
	}

	return qs, nil
}
