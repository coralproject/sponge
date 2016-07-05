// Itemtypes describe the properties and relationships of Items
package item

import (
//	"github.com/ardanlabs/kit/log"
)

//==============================================================================

// ItemType contains all we need to know in order to handle an Item
type Type struct {
	Name string `bson:"name" json:"name"`
}

//==============================================================================

// Types is a slice of all the active Item Types
var Types map[string]Type

func init() {

	// make that types map
	Types = make(map[string]Type)

}

//==============================================================================

// Register a Type
//  if a type of the same name was already registered, it will
//  be overwritten by the new type
func RegisterType(t Type) {

	Types[t.Name] = t

}

//==============================================================================

// Unregister a Type
//  returns true if the type is unregistered, false if not found
func UnregisterType(t Type) bool {

	// how is this done, no internetz here
	//return Types.Delete(t.Name)

	return false

}

//==============================================================================

// isTypeRegistered takes a type name and finds if it's in the Types
func isRegistered(n string) bool {

	for _, t := range Types {
		if n == t.Name {
			return true
		}
	}

	return false
}
