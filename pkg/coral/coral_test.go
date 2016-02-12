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
	"github.com/coralproject/pillar/server/model"
	"github.com/coralproject/sponge/pkg/strategy"
)

var (
	server  *httptest.Server
	path    string
	fakeStr strategy.Strategy

	oStrategy  string
	oPillarURL string
)

func setup() {

	// Save original enviroment variables
	oStrategy = os.Getenv("STRATEGY_CONF")
	oPillarURL = os.Getenv("PILLAR_URL")

	// Initialize log
	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.DEV
		}
		return ll
	}

	log.Init(os.Stderr, logLevel)

	// Initialization of stub server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var err error

		// check that the row is what we want it to be
		switch r.RequestURI {
		case "/api/import/user": // if user, the payload should be a user kind of payload
			// decode the user
			user := model.User{}
			err = json.NewDecoder(r.Body).Decode(&user)
		case "/api/import/asset": // if asset, the payload should be an asset kind of payload
			// decode the asset
			asset := model.Asset{}
			err = json.NewDecoder(r.Body).Decode(&asset)
		case "/api/import/comment": // if comment, the payload should be a comment kind of payload
			// decode the comment
			comment := model.Comment{}
			err = json.NewDecoder(r.Body).Decode(&comment)
		case "/api/import/index":
			// decode the index
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
	defer server.Close()

	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/tests/fixtures/"

	// Mock strategy configuration
	strategyConf := "../../tests/strategy_test.json"
	e := os.Setenv("STRATEGY_CONF", strategyConf) // IS NOT REALLY SETTING UP THE VARIABLE environment FOR THE WHOLE PROGRAM :(
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}

	// Mock pillar url
	os.Setenv("PILLAR_URL", server.URL)

	// Initialize coral
	Init()
}

func teardown() {

	// recover the environment variables

	os.Setenv("STRATEGY_CONF", oStrategy)
	os.Setenv("PILLAR_URL", oPillarURL)
}

func TestMain(m *testing.M) {

	setup()
	code := m.Run()
	teardown()

	os.Exit(code)
}

// We need to test that the mockup server is working by itself
func TestMockupServer(t *testing.T) {

	method := "POST"
	urlStr := server.URL + "/api/import/user"
	row := map[string]interface{}{"juan": 3, "pepe": "what"}
	juser, err := json.Marshal(row)
	payload := bytes.NewBuffer(juser)
	request, err := http.NewRequest(method, urlStr, payload)
	if err != nil {
		t.Fatalf("expect not error and got one %s.", err)
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("expect not error and got one %s.", err)
	}
	defer response.Body.Close()

}

// GetFixture retrieves a query record from the filesystem for testing.
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

// Signature: AddRow(data []byte, tableName string) error
func TestAddRowWrongTable(t *testing.T) {

	var data []byte

	tableName := "wrongTable"

	e := AddRow(data, tableName)
	if e == nil {
		t.Fatal("expecting error and got none.")
	}

}

// test that is sent to the right collection if it is User
func TestAddUserRow(t *testing.T) {

	// Test Data
	newrow, e := GetFixture("users.json")
	if e != nil {
		t.Fatalf("error with the test data: %s.", e)
	}

	var data []byte
	data, e = json.Marshal(newrow)
	if e != nil {
		t.Fatalf("error with the test data: %s.", e)
	}

	tableName := "user"

	e = AddRow(data, tableName)
	if e != nil {
		t.Fatalf("expecting not error but got one %v.", e)
	}
}

// test that is sent to the right collection if it is Asset
func TestAddAssetRow(t *testing.T) {

	// Test Data
	newrow, e := GetFixture("assets.json")
	if e != nil {
		t.Fatalf("error with the test data: %s.", e)
	}

	var data []byte
	data, e = json.Marshal(newrow)
	if e != nil {
		t.Fatalf("error with the test data: %s.", e)
	}

	tableName := "asset"

	e = AddRow(data, tableName)
	if e != nil {
		t.Fatalf("expecting not error but got one %v.", e)
	}
}

// test that is sent to the right collection if it is Comment
func TestAddCommentRow(t *testing.T) {

	// Test Data
	newrow, e := GetFixture("comments.json")
	if e != nil {
		t.Fatalf("error with the test data: %s.", e)
	}

	var data []byte
	data, e = json.Marshal(newrow)
	if e != nil {
		t.Fatalf("error with the test data: %s.", e)
	}

	tableName := "comment"

	e = AddRow(data, tableName)
	if e != nil {
		t.Fatalf("expecting not error but got one %v.", e)
	}

}

// test the request on create index
func TestCreateIndex(t *testing.T) {

	tableName := "comment"

	e := CreateIndex(tableName)
	if e != nil {
		t.Fatalf("expecting not error but got one %v.", e)
	}
}

func TestCreateIndexError(t *testing.T) {

	tableName := "itdoesnotexist"

	e := CreateIndex(tableName)
	if e == nil {
		t.Fatalf("expecting an error but got none.")
	}
}
