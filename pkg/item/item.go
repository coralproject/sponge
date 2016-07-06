package item

import (
	"errors"
	"fmt"
	"time"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"

	gc "github.com/patrickmn/go-cache"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Contains the name of Mongo collections.
const (
	Collection        = "coral_items"
	CollectionHistory = "coral_items_history"
)

// Set of error variables.
var (
	ErrNotFound = errors.New("Item Not found")
)

// =============================================================================

// c contans a cache of set values. The cache will maintain items for one
// second and then marked as expired. This is a very small cache so the
// gc time will be every hour.

const (
	expiration = time.Second
	cleanup    = time.Hour
)

var cache = gc.New(expiration, cleanup)

// Items are trasparently created or updated depending on thier existence
func Upsert(context interface{}, db *db.DB, item *Item) error {

	// validate our item
	if err := item.Validate(); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	// We need to know if this is a new set.
	var new bool
	if _, err := GetById(context, db, item.Id); err != nil {
		if err != ErrNotFound {
			log.Error(context, "Upsert", err, "Completed")
			return err
		}

		new = true
	}

	// Insert or update the item.
	f := func(c *mgo.Collection) error {
		q := bson.M{"_id": item.Id}
		log.Dev(context, "Upsert", "MGO : db.%s.upsert(%s, %s)", c.Name, mongo.Query(q), mongo.Query(item))
		_, err := c.Upsert(q, item)
		return err
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	if new {
		// historical code
	}

	// if the item isn't new it may be in various caches
	//   flush the whole cache
	if !new {
		cache.Flush()
	}

	return nil
}

// GetById retrieves an item by its id.
func GetById(context interface{}, db *db.DB, id bson.ObjectId) (*Item, error) {
	log.Dev(context, "GetById", "Started : Id[%s]", id.Hex())

	var item Item

	// check if the item is in the cache
	key := "item-" + id.Hex()
	if v, found := cache.Get(key); found {
		item := v.(Item)
		log.Dev(context, "GetById", "Completed : CACHE : Item[%+v]", &item)
		return &item, nil
	}

	// query the database for the item
	f := func(c *mgo.Collection) error {
		q := bson.M{"_id": id}
		log.Dev(context, "GetById", "MGO : db.%s.findOne(%s)", c.Name, mongo.Query(q))
		return c.Find(q).One(&item)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetById", err, "Completed")
		return nil, err
	}

	// set the cache: TODO, caching based on type params
	cache.Set(key, item, gc.DefaultExpiration)

	log.Dev(context, "GetById", "Completed : Item[%+v]", &item)
	return &item, nil
}
