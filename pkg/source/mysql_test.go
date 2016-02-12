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

	// no error
	data, err := mm.GetData(coralName, offset, limit, orderby)
	if err != nil {
		t.Fatalf("expected no error, got %s.", err)
	}

	// data should be []map[string]interface{}
	expectedlen := 99
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
