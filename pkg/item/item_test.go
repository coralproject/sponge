package item_test

import (
	"fmt"
	"testing"

	"github.com/coralproject/sponge/pkg/item"
	"github.com/coralproject/sponge/pkg/item/ifix"

	"github.com/ardanlabs/kit/cfg"
	//	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
)

// prefix is what we are looking to delete after the test.
const prefix = "ITEM_TEST_O"

func init() {
	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	tests.Init("CORAL")

	// Initialize MongoDB using the `tests.TestSession` as the name of the
	// master session.
	cfg := mongo.Config{
		Host:     cfg.MustString("MONGO_HOST"),
		AuthDB:   cfg.MustString("MONGO_AUTHDB"),
		DB:       cfg.MustString("MONGO_DB"),
		User:     cfg.MustString("MONGO_USER"),
		Password: cfg.MustString("MONGO_PASS"),
	}
	tests.InitMongo(cfg)
}

//==============================================================================

func TestUpsertCreateItem(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	// register the item types
	ifix.RegisterTypes("types.json")

	dataSets, err := ifix.Get("data.json")
	if err != nil {
		t.Fatalf("\t%s\tShould be able to load item data.json fixture", tests.Failed, err)
	}

	t.Logf("\t%s\tShould be able to create items from data.", tests.Success)
	for _, d := range *dataSets {
		fmt.Println(d)
		i, err := item.Create("coral_comment", 1, d)
		if err != nil {
			t.Fatalf("\t%s\tCould not create item from data : %v", tests.Failed, err)
		}
		fmt.Println(i)
	}

}
