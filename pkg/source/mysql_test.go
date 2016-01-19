/* package source_test is doing unit tests for the source package */
package source

import (
	"fmt"
	"os"
	"testing"
)

var m Sourcer
var mm MySQL

func setup() {

	//setup environment variable
	// strategy File
	// create mysql db

	// MOCK STRATEGY CONF
	strategyConf := "../../tests/strategy_test.json"
	e := os.Setenv("STRATEGY_CONF", strategyConf) // IS NOT REALLY SETTING UP THE VARIABLE environment FOR THE WHOLE PROGRAM :(
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}

	var ok bool

	m, e = New("mysql") // function being tested
	if e != nil {
		fmt.Printf("Error when calling the function, %v.\n", e)
	}

	mm, ok = m.(MySQL)
	if !ok {
		fmt.Println("It should return a type MySQL")
	}
}

//Signature: (m MySQL) GetTables() ([]string, error)
func TestGetTables(t *testing.T) {

	setup()

	s, e := mm.GetTables()
	if e != nil {
		t.Fatalf("expected no error, got %s.", e)
	}

	if len(s) != 4 { // 4 is in the seeds when creating the test strategy file
		t.Fatalf("expected 4, got %s", len(s))
	}
}

// Signature: (m MySQL) GetData(coralTableName string, offset int, limit int, orderby string) ([]map[string]interface{}, error)
func TestGetData(t *testing.T) {

	setup()

	// Default Flags
	coralName := "comment"
	offset := 0
	limit := 9999999999
	orderby := ""

	// no error
	data, err := mm.GetData(coralName, offset, limit, orderby)
	if err != nil {
		t.Fatalf("expected no error, got %s.", err)
	}

	// data should be []map[string]interface{}
	expectedlen := 24999
	if len(data) != expectedlen { // this is a setup for the seed data
		t.Fatalf("exptected %d, got %d", expectedlen, len(data))
	}

}

// Signature: (m MySQL) GetQueryData(coralTableName string, offset int, limit int, orderby string, ids []string) ([]map[string]interface{}, error)
// func TestGetQueryData(t *testing.T) {
// 	setup()
//
// 	// Default Flags
// 	coralTableName := "comment"
// 	offset := 0
// 	limit := 9999999999
// 	orderby := ""
//
// 	ids := []string{"16570043", "16570056", "16570088", "16570101", "16570134"}
//
// 	// no error
// 	data, err := mm.GetQueryData(coralTableName, offset, limit, orderby, ids)
// 	if err != nil {
// 		t.Fatalf("it should have no error %s.", err)
// 	}
//
// 	// data should be []map[string]interface{}
// 	if len(data) != 5 { // this is a setup for the seed data
// 		t.Fatalf("got %d, it should be 5", len(data))
// 	}
// }
