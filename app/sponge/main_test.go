package tests

import (
	"testing"

	"github.com/ardanlabs/kit/tests"
)

func init() {
	tests.Init("SPONGE")
}

//==============================================================================

// TestMain helps to clean up the test data.
func TestMain(m *testing.M) {
	// db := db.NewMGO()
	// defer db.CloseMGO()
	//
	// query.GenerateTestData(db)
	// defer query.DropTestData()
	//
	// loadQuery(db, "basic.json")
	// loadQuery(db, "basic_var.json")
	// defer query.RemoveTestSets(db)
	//
	// m.Run()
}
