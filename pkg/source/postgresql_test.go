/* package source_test is doing unit tests for the source package */
package source

import "testing"

//Signature: (m MySQL) GetTables() ([]string, error)
func TestPostgresGetTables(t *testing.T) {

	setupPostgreSQL()

	s, e := GetEntities()
	if e != nil {
		t.Fatalf("expected no error, got %s.", e)
	}

	expectedLen := 2
	if len(s) != expectedLen {
		t.Fatalf("got %d, it should be %d", len(s), expectedLen)
	}

	if s[0] != "assets" {
		t.Fatalf("got %s, it should be asset", s[0])
	}

	if s[1] != "users" {
		t.Fatalf("got %s, it should be users", s[1])
	}

	teardown()
}

// Signature: (m MySQL) GetData(coralTableName string, offset int, limit int, orderby string) ([]map[string]interface{}, error)
func TestPostgresGetData(t *testing.T) {

	setupPostgreSQL()

	// Default Flags
	coralName := "users"
	offset := 0
	limit := 9999999999
	var orderby string
	var query string

	// no error
	data,_, err := mp.GetData(coralName, offset, limit, orderby, query)
	if err != nil {
		t.Fatalf("expected no error, got %s.", err)
	}

	// data should be []map[string]interface{}
	expectedlen := 161
	if len(data) != expectedlen { // this is a setup for the seed data
		t.Fatalf("exptected %d, got %d", expectedlen, len(data))
	}
}
