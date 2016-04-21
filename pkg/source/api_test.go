/* package source_test is doing unit tests for the source package */
package source

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	str "github.com/coralproject/sponge/pkg/strategy"
)

var (
	server  *httptest.Server
	path    string
	fakeStr str.Strategy
)

func mockAPI() string {

	// Initialization of stub server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var err error
		var d map[string]interface{}

		file, e := ioutil.ReadFile("../../tests/response.json")
		if e != nil {
			fmt.Printf("ERROR %v on setting up response in the test.", e)
			os.Exit(1)
		}

		json.Unmarshal(file, &d)

		// err = errors.New("Bad request")

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

func TestMain(m *testing.M) {

	setupAPI()

	code := m.Run()
	teardown()

	os.Exit(code)
}

func TestAPIGetData(t *testing.T) {

	data, _, err := mapi.GetWebServiceData()
	if err != nil {
		t.Fatalf("expected no error, got '%s'.", err)
	}

	if len(data) == 0 {
		t.Fatalf("expected some entry, got zero.")
	}

}

func TestGetAPIData(t *testing.T) {
	setupAPI()

	pageAfter := "1.399743732"

	// no error
	data, pageAfter1, err := mapi.GetFireHoseData(pageAfter)
	if err != nil {
		t.Fatalf("expected no error, got '%s'.", err)
	}

	// data should be []map[string]interface{}
	expectedlen := 2
	if len(data) != expectedlen { // this is a setup for the seed data
		t.Fatalf("expected %d, got %d", expectedlen, len(data))
	}

	expectedUser := "Duck504"
	if data[2]["actor.title"] != expectedUser {
		t.Fatalf("expected %s, got %s", expectedUser, data[2]["actor.title"])
	}

	if pageAfter1 == pageAfter {
		t.Fatalf("expected different pages %s and %s", pageAfter1, pageAfter)
	}

	data, _, err = mapi.GetFireHoseData(pageAfter1)
	if err != nil {
		t.Fatalf("expected no error, got '%s'.", err)
	}

	expectedlen = 4
	if len(data) != expectedlen { // this is a setup for the seed data
		t.Fatalf("expected %d, got %d", expectedlen, len(data))
	}

}
