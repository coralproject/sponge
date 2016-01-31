/* package source_test is doing unit tests for the source package */
package source

import "testing"

// Signature: (m MySQL) GetData(coralTableName string, offset int, limit int, orderby string) ([]map[string]interface{}, error)
func TestGetMongoData(t *testing.T) {

	setup("mongo")

	// Default Flags
	coralName := "comment"
	offset := 0
	limit := 9999999999
	orderby := ""

	// no error
	_, err := mmongo.GetData(coralName, offset, limit, orderby)
	if err != nil {
		t.Fatalf("expected no error, got %s.", err)
	}

	// data should be []map[string]interface{}
	// expectedlen := 99
	// if len(data) != expectedlen { // this is a setup for the seed data
	// 	t.Fatalf("exptected %d, got %d", expectedlen, len(data))
	// }

	teardown()

}

// Signature: (m MySQL) GetQueryData(coralTableName string, offset int, limit int, orderby string, ids []string) ([]map[string]interface{}, error)
func TestGetMongoQueryData(t *testing.T) {
	setupMongo()

	// Default Flags
	coralTableName := "comment"
	offset := 0
	limit := 9999999999
	orderby := ""

	ids := []string{"16570043", "16570056", "16570088", "16570101", "16570134"}

	// no error
	data, err := mmongo.GetQueryData(coralTableName, offset, limit, orderby, ids)
	if err != nil {
		t.Fatalf("it should have no error %s.", err)
	}

	// data should be []map[string]interface{}
	if len(data) != 5 { // this is a setup for the seed data
		t.Fatalf("got %d, it should be 5", len(data))
	}
	teardown()
}
