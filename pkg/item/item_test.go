package item_test

import (
	"testing"

	"github.com/coralproject/sponge/pkg/item"
	"github.com/coralproject/sponge/pkg/item/ifix"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
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
		DB:       cfg.MustString("MONGO_DB") + "_test",
		User:     cfg.MustString("MONGO_USER"),
		Password: cfg.MustString("MONGO_PASS"),
	}
	tests.InitMongo(cfg)
}

//==============================================================================

func TestUpsertCreateItem(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	// get db connection
	//  should this be moved to a shared testing package?
	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	// register the item types
	ifix.RegisterTypes("types.json")

	dataSets, err := ifix.Get("data.json")
	if err != nil {
		t.Fatalf("\t%s\tShould be able to load item data.json fixture", tests.Failed, err)
	}

	t.Logf("\t%s\tCreate, save and update items.", tests.Success)
	for _, d := range *dataSets {

		i, err := item.Create("coral_comment", 1, d)
		if err != nil {
			t.Fatalf("\t%s\tCould not create item from data: %v", tests.Failed, err)
		}

		_, err = item.Create("how many aphorisms do you remember? what do you take from them?", 1, d)
		if err == nil {
			t.Fatalf("\t%s\tShould not be able to create with unregistered type: %v", tests.Failed, err)
		}

		t.Logf("\t%s\tShould be able to insert items.", tests.Success)
		err = item.Upsert(tests.Context, db, &i)
		if err != nil {
			t.Fatalf("\t%s\tCould not upsert (insert) item: %v", tests.Failed, err)
		}

		// bump the version
		i.Version = 2
		t.Logf("\t%s\tShould be able to update items.", tests.Success)
		err = item.Upsert(tests.Context, db, &i)
		if err != nil {
			t.Fatalf("\t%s\tCould not upsert (update) item: %v", tests.Failed, err)
		}

		i2, err := item.GetById(tests.Context, db, i.Id)
		if err != nil {
			t.Fatalf("\t%s\tCould not GetById item: %v", tests.Failed, err)
		}

		if i.Version != i2.Version {
			t.Fatalf("\t%s\tDid not see update of version", tests.Failed)

		}

	}

}
