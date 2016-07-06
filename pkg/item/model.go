package item

import (
	"errors"
	"time"

	"gopkg.in/bluesuncorp/validator.v8"
	"gopkg.in/mgo.v2/bson"
)

const (
	DefaultVersion = 1
)

//==============================================================================

// validate is used to perform model field validation.
var validate *validator.Validate

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
}

//==============================================================================

// Result contains Items returned via a read operation and metadata
type Result struct {
	Items []Item    `bson:"items" json:"items"`
	Date  time.Time `bson:"date" json:"date"`
}

//==============================================================================

// ItemData is what an Item can hold
//  Should be the intersection of the db and transport protocols supported
type ItemData interface{}

//==============================================================================

// An Item is data, properties and behavior wrapped in the thinnest
//  practical wrapper: Id, Type and Version
// this will be high volume so db and json field names are truncated
type Item struct {
	Id      bson.ObjectId `bson:"_id" json:"id"`
	Type    string        `bson:"t" json:"t"` // ItemType.Name
	Version int           `bson:"v" json:"v"`
	Data    ItemData      `bson:"d" json:"d"`
	Rels    []Rel         `bson:"rels,omitempty" json:"rels,omitempty"`
}

func (i *Item) Validate() error {
	if err := validate.Struct(i); err != nil {
		return err
	}

	return nil
}

// =============================================================================

// create an item out of its type, version and data or die trying
func Create(t string, v int, d ItemData) (Item, error) {

	i := Item{}

	// _initial draft_  Create a mongo id for each new item
	//   we may want to figure out how to make Item.Id
	//   reflect ids found in the data one day
	i.Id = bson.NewObjectId()

	// validate and set type
	if isRegistered(t) == false {
		return i, errors.New("Type not recognized: " + t)
	}
	i.Type = t

	// set default version if zero value
	if v == 0 {
		v = DefaultVersion
	}
	i.Version = v

	// set the data into the item
	i.Data = d

	return i, nil
}
