// Itemtypes describe the properties and relationships of Items
package item

import (
	"fmt"
)

//==============================================================================

// ItemType contains all we need to know in order to handle an Item
type ItemType struct {
	Name string `bson:"name" json:"name"`
}

//==============================================================================

// Types is a slice of all the active ItemTypes
var Types []ItemType

func init() {

	fmt.Printl("Initializing ItemTypes")

	// stub out initial ItemTypes
	Types = append(Types, ItemType{"user"})
	Types = append(Types, ItemType{"comment"})
	Types = append(Types, ItemType{"asset"})

}
