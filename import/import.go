/*
Package data implements a way to get data from external sources and push it into our local db.

External possible sources:
* MySQL
* API?

Our local DB:
* MongoDb


A process that
  - mantain the date of the last record
  - keep control of any update to the db
  - push any updates to localdb
	- get from configuration how frequently has to update localdb

*/
package data // To Do: Look for a better name

import (
	"log"

	"github.com/coralproject/mod-data-import/import/source"
)

// Strategy structure has information about the source and destinty of the data
// To Do: look for a better name (we also have config.Strategy)
type Strategy struct {
	LocalDB
	Source
}

// NewStrategy initialize the data for the new strurcture
func NewStrategy() (*Strategy, Error) {

	ndb, err := localDB.NewLocalDB()
	if err != nil {
		log.Fatal("Error when trying to create new strategy. ", err)
	}

	ns, err := source.NewSource()
	if err != nil {
		log.Fatal("Error when trying to create new strategy. ", err)
	}
	s := Strategy{ndb, ns}

	return &s, err
}

// Import copy the data from external source to local database
func (s *Strategy) Import() Error {

	err := s.LocalDB.push(s.Source.GetNewData())
	if err != nil {
		log.Fatal("Error when importing external source into local database.")
	}
	return err
}

func main() {
	strategy, err := NewStrategy()
	if err != nil {
		log.Fatal("Error when creating new Strategy")
	}
}
