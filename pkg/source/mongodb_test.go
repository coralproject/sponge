package source

import "testing"

func TestMongoGetData(t *testing.T) {
	setupMongo()

	// Default Flags
	coralName := "comment"
	offset := 0
	limit := 9999999999
	orderby := ""

	// no error
	data, err := mdb.GetData(coralName, offset, limit, orderby)
	if err != nil {
		t.Fatalf("expected no error, got %s.", err)
	}

	// data should be []map[string]interface{}
	expectedlen := 0
	if len(data) != expectedlen { // this is a setup for the seed data
		t.Fatalf("expected %d, got %d", expectedlen, len(data))
	}
}

func TestMongoGetQueryData(t *testing.T) {
	t.Skip()
}

// Signature func (m MongoDB) GetTables() ([]string, error) {
func TestMongoGetTables(t *testing.T) {

	setupMongo()

	s, e := mdb.GetTables()
	if e != nil {
		t.Fatalf("expected no error, got %s.", e)
	}

	expectedLen := 3
	if len(s) != expectedLen {
		t.Fatalf("got %d, it should be %d", len(s), expectedLen)
	}

	teardown()
}
