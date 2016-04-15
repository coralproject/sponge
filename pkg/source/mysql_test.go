/* package source_test is doing unit tests for the source package */
package source

import "testing"

//Signature: (m MySQL) GetTables() ([]string, error)
func TestGetTables(t *testing.T) {

	setupMysql()

	s, e := GetEntities()
	if e != nil {
		t.Fatalf("expected no error, got %s.", e)
	}

	expectedLen := 3
	if len(s) != expectedLen {
		t.Fatalf("got %d, it should be %d", len(s), expectedLen)
	}

	if s[0] != "assets" {
		t.Fatalf("got %s, it should be assets", s[0])
	}

	if s[2] != "comments" {
		t.Fatalf("got %s, it should be comments", s[0])
	}

	teardown()
}

// Signature: (m MySQL) GetData(coralTableName string, offset int, limit int, orderby string) ([]map[string]interface{}, error)
func TestGetData(t *testing.T) {

	setupMysql()

	// Default Flags
	coralName := "comments"
	options := &Options{offset: 0,
		limit:   9999999999,
		orderby: "createdate",
		query:   ""}

	// no error
	data, err := mm.GetData(coralName, options)
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
	coralName := "assets"
	options := &Options{
		offset:  0,
		limit:   9999999999,
		orderby: "",
		query:   "updatedate > 2013-12-12",
	}

	// no error
	data, err := mm.GetData(coralName, options)
	if err != nil {
		t.Fatalf("expected no error, got %s.", err)
	}

	// data should be []map[string]interface{}
	expectedlen := 9
	if len(data) != expectedlen { // this is a setup for the seed data
		t.Fatalf("exptected %d, got %d", expectedlen, len(data))
	}

}
