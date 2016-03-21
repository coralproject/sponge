/* package source_test is doing unit tests for the source package */
package source

import "testing"

//Signature: (m MySQL) GetTables() ([]string, error)
func TestGetTables(t *testing.T) {

	setupMysql()

	s, e := mm.GetTables()
	if e != nil {
		t.Fatalf("expected no error, got %s.", e)
	}

	expectedLen := 3
	if len(s) != expectedLen {
		t.Fatalf("got %d, it should be %d", len(s), expectedLen)
	}

	if s[0] != "asset" {
		t.Fatalf("got %s, it should be asset", s[0])
	}

	if s[2] != "comment" {
		t.Fatalf("got %s, it should be asset", s[0])
	}

	teardown()
}

// Signature: (m MySQL) GetData(coralTableName string, offset int, limit int, orderby string) ([]map[string]interface{}, error)
func TestGetData(t *testing.T) {

	setupMysql()

	// Default Flags
	coralName := "comment"
	offset := 0
	limit := 9999999999
	orderby := ""
	query := ""

	// no error
	data, err := mm.GetData(coralName, offset, limit, orderby, query)
	if err != nil {
		t.Fatalf("expected no error, got %s.", err)
	}

	// data should be []map[string]interface{}
	expectedlen := 99
	if len(data) != expectedlen { // this is a setup for the seed data
		t.Fatalf("exptected %d, got %d", expectedlen, len(data))
	}

}

// Signature: (m MySQL) GetData(coralTableName string, offset int, limit int, orderby string) ([]map[string]interface{}, error)
func TestQueryGetData(t *testing.T) {

	setupMysql()

	// Default Flags
	coralName := "asset"
	offset := 0
	limit := 9999999999
	orderby := ""
	query := "updatedate > 2013-12-12"

	// no error
	data, err := mm.GetData(coralName, offset, limit, orderby, query)
	if err != nil {
		t.Fatalf("expected no error, got %s.", err)
	}

	// data should be []map[string]interface{}
	expectedlen := 9
	if len(data) != expectedlen { // this is a setup for the seed data
		t.Fatalf("exptected %d, got %d", expectedlen, len(data))
	}

}
