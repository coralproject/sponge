/* package source_test is doing unit tests for the source package */
package source

import (
	"fmt"
	"os"
	"testing"
)

var m Sourcer
var mm MySQL

var oStrategy string

func setup() {

	//setup environment variable
	// strategy File
	// create mysql db

	oStrategy = os.Getenv("STRATEGY_CONF")

	// MOCK STRATEGY CONF
	strategyConf := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/tests/strategy_test.json"
	e := os.Setenv("STRATEGY_CONF", strategyConf) // IS NOT REALLY SETTING UP THE VARIABLE environment FOR THE WHOLE PROGRAM :(
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}

	var ok bool

	m, e = New("mysql") // function being tested
	if e != nil {
		fmt.Printf("error when calling the function, %v.\n", e)
	}

	mm, ok = m.(MySQL)
	if !ok {
		fmt.Println("it should return a type MySQL")
	}
}

func teardown() {
	e := os.Setenv("STRATEGY_CONF", oStrategy)
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}
}

//Signature: (m MySQL) GetTables() ([]string, error)
func TestGetTables(t *testing.T) {

	setup()

	s, e := mm.GetTables()
	if e != nil {
		t.Fatalf("error when getting tables, %v.", e)
	}

	expectedLen := 3
	if len(s) != expectedLen {
		t.Fatalf("got %d, it should be %d", len(s), expectedLen)
	}

	teardown()
}

// Signature: (m MySQL) GetData(coralTableName string, offset int, limit int, orderby string) ([]map[string]interface{}, error)
func TestGetData(t *testing.T) {
	setup()

	// Default Flags
	coralTableName := "comment"
	offset := 0
	limit := 9999999999
	orderby := ""

	// no error
	data, err := mm.GetData(coralTableName, offset, limit, orderby)
	if err != nil {
		t.Fatalf("it should have no error %s.", err)
	}

	expectedLen := 24999
	// data should be []map[string]interface{}
	if len(data) != expectedLen { // this is a setup for the seed data
		t.Fatalf("got %d, it should be %d", len(data), expectedLen)
	}

}

// Signature: (m MySQL) GetQueryData(coralTableName string, offset int, limit int, orderby string, ids []string) ([]map[string]interface{}, error)
func TestGetQueryData(t *testing.T) {
	setup()

	// Default Flags
	coralTableName := "comment"
	offset := 0
	limit := 9999999999
	orderby := ""

	ids := []string{"16570043", "16570056", "16570088", "16570101", "16570134"}

	// no error
	data, err := mm.GetQueryData(coralTableName, offset, limit, orderby, ids)
	if err != nil {
		t.Fatalf("it should have no error %s.", err)
	}

	// data should be []map[string]interface{}
	if len(data) != 5 { // this is a setup for the seed data
		t.Fatalf("got %d, it should be 5", len(data))
	}
}
