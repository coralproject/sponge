/* package source_test is doing unit tests for the source package */
package source

import (
	"encoding/json"
	"errors"
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

func mockAPI() {

	// Initialization of stub server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var err error
		var d map[string]interface{}

		file, e := ioutil.ReadFile("tests/response.json")
		if e != nil {
			fmt.Println("ERROR on setting up response in the test.")
			os.Exit(1)
		}

		json.Unmarshal(file, &d)
		fmt.Println(d)

		// check that the row is what we want it to be
		switch r.RequestURI {
		case "/v1/search":

			d = make(map[string]interface{})
			// look at query parameters
			// setup d
			//q=((scope: https://www.washingtonpost.com/lifestyle/style/carolyn-hax-stubborn-60-something-parent-refuses-to-see-a-doctor/2015/09/24/299ec776-5e2d-11e5-9757-e49273f05f65_story.html source:washpost.com itemsPerPage: 100 sortOrder:reverseChronological ))&appkey=dev.washpost.com
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

		fmt.Fprintln(w, d)
	}))

	fmt.Println("DEBUG server url: ", server.URL)
}

func TestMain(m *testing.M) {

	setupAPI()

	code := m.Run()
	teardown()

	os.Exit(code)
}

func TestAPIGetData(t *testing.T) {

	coralTableName := ""
	offset := 0
	limit := 0
	orderby := ""
	q := ""

	data, _, err := mapi.GetData(coralTableName, offset, limit, orderby, q)
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
	data, finish, pageAfter1, err := mapi.GetAPIData(pageAfter)
	if err != nil {
		t.Fatalf("expected no error, got '%s'.", err)
	}

	// data should be []map[string]interface{}
	expectedlen := 200
	if len(data) != expectedlen { // this is a setup for the seed data
		t.Fatalf("expected %d, got %d", expectedlen, len(data))
	}

	expectedUser := "Duck504"
	if data[100]["actor.title"] != expectedUser {
		t.Fatalf("expected %s, got %s", expectedUser, data[100]["actor.title"])
	}

	if !finish {
		if pageAfter1 == pageAfter {
			t.Fatalf("expected different pages %s and %s", pageAfter1, pageAfter)
		}

		data, _, pageAfter2, err := mapi.GetAPIData(pageAfter1)
		if err != nil {
			t.Fatalf("expected no error, got '%s'.", err)
		}

		expectedlen := 200
		if len(data) != expectedlen { // this is a setup for the seed data
			t.Fatalf("expected %d, got %d", expectedlen, len(data))
		}

		if pageAfter1 == pageAfter2 {
			t.Fatalf("expected different pages %s and %s", pageAfter1, pageAfter2)
		}

	}
}
