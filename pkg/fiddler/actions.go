package fiddler

import (
	"fmt"
	"log"
	"time"

	"github.com/coralproject/shelf/pkg/srv/comment"
	configuration "github.com/coralproject/sponge/config"
	"github.com/oleiade/reflections"
)

//* ACTIONS *//

// Action denotes an action taken by someone/something on someone/something.
type Action struct {
	comment.Action
}

// 	Type   string    `json:"type" bson:"type"`
// 	UserID string    `json:"user_id" bson:"user_id"`
// 	Value  string    `json:"value" bson:"value"`
// 	Date   time.Time `json:"date" bson:"date"`
// }

// Print only print information about the action
func (a Action) Print() {
	fmt.Println("Action: ", a.UserID, a.Type, a.Value)
}

// Transform get the data from sd
func (a Action) Transform(sd []map[string]interface{}, table configuration.Table) ([]Transformer, error) {

	var action Action
	var actions []Transformer

	// To Do: it needs refactoring as my gut tells me that is quite inefficient
	for _, value := range sd {
		for coralField, f := range table.Fields {

			// convert field f with value value[f] into field coralField
			newValue := transformActionField(f, value[f], coralField)

			err := reflections.SetField(&action, coralField, newValue)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}
		}

		n := len(actions)
		if len(actions) == cap(actions) {
			// actions is full and we must expand
			// Double the size and add 1
			newactions := make([]Transformer, len(actions), 2*len(actions)+1)
			copy(newactions, actions)
			actions = newactions
		}
		actions = actions[0 : n+1]
		//action.Raw = strings.Split(raws, ",")
		actions[n] = action
	}

	return actions, nil
}

//Here we transform the record into what we want (based on the configuration)
// 1. convert types (values are all strings) into the struct
func transformActionField(sourceField string, oldValue interface{}, coralField string) interface{} {

	var newValue interface{}

	// Right now this the simpler thing to do. Needs to look more into reflect to do it
	// dinamically (if it is worth it) or merge this into the comment package
	switch coralField {
	case "Type": //string
		newValue = oldValue
	case "UserID": //string    `json:"parent_id" bson:"parent_d"`
		newValue = oldValue
	case "Value": //string    `json:"asset_id" bson:"asset_id"`
		newValue = oldValue
	case "Date": //time.Time `json:"date_approved" bson:"date_approved"`
		newValue, _ = time.Parse(longForm, oldValue.(string))
		// Actions and Notes are missing
	}

	return newValue
}