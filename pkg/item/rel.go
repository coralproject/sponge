package item

import (
	"errors"

	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/db"
)

//==============================================================================

var (
	ErrRelsTypesNotFound = errors.New("Could not retrieve RelTypes for item")
)

//==============================================================================

// Rel holds an item's relationship to another item
type Rel struct {
	Name string `bson:"n" json:"n"`   // Name of relationship
	Type string `bson:"t" json:"t"`   // Item Type of target
	Id   string `bson:"id" json:"id"` // Id of target (consider storing as native bson.Id?)
}

//==============================================================================

func getDatumByKey(k string, d interface{}) interface{} {

	// get the data as a map for searching
	m := d.(map[string]interface{})

	return m[k]

}

// GetRels looks up an item's relationships and returns them
func GetRels(context interface{}, db *db.DB, item *Item) (*[]Rel, error) {

	var rels []Rel

	// get the rel types for this item's type
	rts := Types[item.Type].Rels

	// for each reltype
	for _, rt := range rts {

		// find the foreign key value in the item data
		fkv := getDatumByKey(rt.Field, item.Data)

		// if there is not value, skip this rel
		if fkv == nil {
			continue
		}

		// create the field path for the foreign key field
		fkf := "d." + Types[item.Type].IdField

		// otherwise build a query
		var q = bson.M{"t": rt.Type, fkf: fkv}

		items, err := GetByQuery(context, db, q)
		if err != nil {
			// how should we handle not being able to look up related items?
			//  we probably don't want to prevent the insert without measures to recover the item
			//  although we would be causing data inconsistencies?  maybe a flag that
			//  relations need to be re-queried?
		}

		// for each item
		for _, i := range *items {

			// create the relationship
			r := Rel{
				Name: rt.Name,
				Type: i.Type,
				Id:   i.Id.Hex(),
			}

			// add it to the list
			rels = append(rels, r)
		}

	}

	return &rels, nil
}
