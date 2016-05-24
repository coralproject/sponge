package source

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	uuidimported "github.com/pborman/uuid"
)

const (
	cfgLoggingLevel = "LOGGING_LEVEL"
)

var (
	m    Sourcer
	mm   MySQL
	mdb  MongoDB
	mp   PostgreSQL
	mapi API
)

// we save the strategy env to recover it on teardown
var oStrategy string

func init() {
	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.DEV
		}
		return ll
	}

	log.Init(os.Stderr, logLevel)
}

func setupMysql() {
	// Initialize logging
	logLevel := func() int {
		ll, err := cfg.Int(cfgLoggingLevel)
		if err != nil {
			return log.USER
		}
		return ll
	}
	log.Init(os.Stderr, logLevel)

	oStrategy = os.Getenv("STRATEGY_CONF")

	// MOCK STRATEGY CONF
	strategyConf := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/tests/strategy_nyt_test.json"
	e := os.Setenv("STRATEGY_CONF", strategyConf) // IS NOT REALLY SETTING UP THE VARIABLE environment FOR THE WHOLE PROGRAM :(
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}

	var ok bool

	u := uuidimported.New()

	s, e := Init(u)
	if e != nil {
		fmt.Printf("Error when initializing strategy, %v.\n", e)
	}
	m, e := New(s) // setting up new source
	if e != nil {
		fmt.Printf("Error when calling the function, %v.\n", e)
	}

	mm, ok = m.(MySQL)
	if !ok {
		fmt.Println("It should return a type MySQL")
	}
}

func setupMongo() {

	// Initialize logging
	logLevel := func() int {
		ll, err := cfg.Int(cfgLoggingLevel)
		if err != nil {
			return log.USER
		}
		return ll
	}
	log.Init(os.Stderr, logLevel)

	oStrategy = os.Getenv("STRATEGY_CONF")

	// MOCK STRATEGY CONF
	strategyConf := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/tests/strategy_wapo_test.json"
	e := os.Setenv("STRATEGY_CONF", strategyConf)
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}

	var ok bool

	u := uuidimported.New()
	s, e := Init(u)
	if e != nil {
		fmt.Printf("Error when initializing strategy, %v.\n", e)
	}

	m, e := New(s) // setting up new source
	if e != nil {
		fmt.Printf("Error when calling the function, %v.\n", e)
	}

	mdb, ok = m.(MongoDB)
	if !ok {
		fmt.Println("It should return a type MontoDB")
	}
}

func setupPostgreSQL() {

	// Initialize logging
	logLevel := func() int {
		ll, err := cfg.Int(cfgLoggingLevel)
		if err != nil {
			return log.USER
		}
		return ll
	}
	log.Init(os.Stderr, logLevel)

	oStrategy = os.Getenv("STRATEGY_CONF")

	// MOCK STRATEGY CONF
	strategyConf := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/tests/strategy_discourse_test.json"
	e := os.Setenv("STRATEGY_CONF", strategyConf) // IS NOT REALLY SETTING UP THE VARIABLE environment FOR THE WHOLE PROGRAM :(
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}

	var ok bool

	u := uuidimported.New()

	s, e := Init(u)
	if e != nil {
		fmt.Printf("Error when initializing strategy, %v.\n", e)
	}
	m, e := New(s) // setup new source
	if e != nil {
		fmt.Printf("Error when calling the function, %v.\n", e)
	}

	mp, ok = m.(PostgreSQL)
	if !ok {
		fmt.Println("It should return a type PostgreSQL")
	}
}

func mockAPI() string {

	// Initialization of stub server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var err error
		var d map[string]interface{}

		file, err := ioutil.ReadFile("../../tests/response.json")
		if err != nil {
			fmt.Printf("ERROR %v on setting up response in the test.", err)
			os.Exit(1)
		}

		json.Unmarshal(file, &d)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		if err == nil {
			w.WriteHeader(http.StatusOK)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(file)

		fmt.Fprintln(w, file)
	}))

	return server.URL
}

func setupAPI() {

	// Mock the API Server
	serverurl := mockAPI()

	// Initialize logging
	logLevel := func() int {
		ll, err := cfg.Int(cfgLoggingLevel)
		if err != nil {
			return log.USER
		}
		return ll
	}
	log.Init(os.Stderr, logLevel)

	oStrategy = os.Getenv("STRATEGY_CONF")

	// MOCK STRATEGY CONF
	strategyConf := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/tests/strategy_wapo_api_test.json"

	// write down serverurl into the strategyconf
	st := strategy
	content, err := ioutil.ReadFile(strategyConf)
	err = json.Unmarshal(content, &st)

	st.Credentials.Service.Endpoint = serverurl

	bst, err := json.Marshal(st)
	if err != nil {
		fmt.Println("Error when trying to marshall back the strategy file with server url of the mock server: ", err)
	}
	mode := os.FileMode(0777)
	err = ioutil.WriteFile(strategyConf, bst, mode)
	if err != nil {
		fmt.Println("Error when saving back the strategy file with the mock server: ", err)
	}

	e := os.Setenv("STRATEGY_CONF", strategyConf)
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}

	var ok bool

	u := uuidimported.New()
	s, e := Init(u)
	if e != nil {
		fmt.Printf("Error when initializing strategy, %v.\n", e)
	}

	m, e := New(s)
	if e != nil {
		fmt.Printf("Error when calling the function, %v.\n", e)
	}

	mapi, ok = m.(API)
	if !ok {
		fmt.Println("It should return a type API")
	}
	attributes := "scope:https://www.washingtonpost.com/lifestyle/style/carolyn-hax-stubborn-60-something-parent-refuses-to-see-a-doctor/2015/09/24/299ec776-5e2d-11e5-9757-e49273f05f65_story.html source:washpost.com itemsPerPage:100 sortOrder:reverseChronological"
	mapi.Connection = fmt.Sprintf("%s/v1/search?q=((%s))&appkey=dev.washpost.com", serverurl, attributes)

}

func teardown() {
	e := os.Setenv("STRATEGY_CONF", oStrategy)
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}
}

// NewSource returns a new MySQL struct
// Signature: New(d string) (Sourcer, error) {
// It depends on the credentials to get the connection string
func TestPostgreSQLNewSource(t *testing.T) {

	m, e := New("postgresql") // function being tested
	if e != nil {
		t.Fatalf("error when calling the function, %v.", e)
	}

	mp, ok := m.(PostgreSQL)
	if !ok {
		t.Fatalf("it should return a type PostgreSQL")
	}

	// m should have a valid connection string
	if mp.Connection == "" {
		t.Fatalf("connection string should not be nil")
	}

	// m should not have a database connection
	if mp.Database != nil {
		t.Error("database should be nil.")
	}
}

// NewSource returns a new mongodb struct
// Signature: New(d string) (Sourcer, error) {
// It depends on the credentials to get the connection string
func TestMongoDBNewSource(t *testing.T) {

	m, e := New("mongodb") // function being tested
	if e != nil {
		t.Fatalf("error when calling the function, %v.", e)
	}

	mdb, ok := m.(MongoDB)
	// it returns type MySQL
	if !ok {
		t.Fatalf("it should return a type MySQL")
	}

	// m should have a valid connection string
	if mdb.Connection == "" {
		t.Fatalf("connection string should not be nil")
	}

	// m should not have a database connection
	if mdb.Database != nil {
		t.Error("database should be nil.")
	}
}
