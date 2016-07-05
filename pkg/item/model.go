package item

import (
	"errors"
	"time"

	"gopkg.in/bluesuncorp/validator.v8"
)

const (
	Collection     = "items"
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

// Item contains the item itself
// this will be high volume so db and json field names are truncated
type Item struct {
	Id      string              `bson:"_id" json:"id"`
	Type    string              `bson:"t" json:"t"` // ItemType.Name
	Version int                 `bson:"v" json:"v"`
	Data    [string]interface{} `bson:"d" json:"d"`
}

func (i *Item) Validate() error {
	if err := validate.Struct(q); err != nil {
		return err
	}

	// set default version if zero value
	if i.Version == 0 {
		i.Version = DefaultVersion
	}

	// validate type
	typeFound = false
	for _, t := range Types {
		if i.Type == t {
			typeFound = true
			break
		}
	}
	if typeFound == fale {
		return errors.New("Type not recognized: " + i.Type)
	}

	return nil
}
