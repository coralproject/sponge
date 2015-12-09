package fiddler

import (
	"fmt"
	"log"
	"time"

	"github.com/coralproject/shelf/pkg/srv/comment"
	configuration "github.com/coralproject/sponge/config"
	"github.com/oleiade/reflections"
)

//* NOTES */

// Note is embedding the comment package to extend it
type Note struct {
	comment.Note
}

// Note denotes a note by a user in the system.
// UserID string    `json:"user_id" bson:"user_id"`
// Body   string    `json:"body" bson:"body"`
// Date   time.Time `json:"date" bson:"date"`

// Print only print information about the comment
func (n Note) Print() {
	fmt.Println("Note: ", n.UserID, n.Body)
}

// Transform get the data from sd
func (n Note) Transform(sd []map[string]interface{}, table configuration.Table) ([]Transformer, error) {

	var note Note
	var notes []Transformer

	// To Do: it needs refactoring as my gut tells me that is quite inefficient
	for _, value := range sd {
		for coralField, f := range table.Fields {

			// convert field f with value value[f] into field coralField
			newValue := transformNoteField(f, value[f], coralField)

			err := reflections.SetField(&note, coralField, newValue)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}
		}

		n := len(notes)
		if len(notes) == cap(notes) {
			// actions is full and we must expand
			// Double the size and add 1
			newnotes := make([]Transformer, len(notes), 2*len(notes)+1)
			copy(newnotes, notes)
			notes = newnotes
		}
		notes = notes[0 : n+1]
		//action.Raw = strings.Split(raws, ",")
		notes[n] = note
	}

	return notes, nil
}

//Here we transform the record into what we want (based on the configuration)
// 1. convert types (values are all strings) into the struct
func transformNoteField(sourceField string, oldValue interface{}, coralField string) interface{} {

	var newValue interface{}

	// Right now this the simpler thing to do. Needs to look more into reflect to do it
	// dinamically (if it is worth it) or merge this into the comment package
	switch coralField {
	case "UserID": //string
		newValue = oldValue
	case "Body": //string    `json:"parent_id" bson:"parent_d"`
		newValue = oldValue
	case "Date": //time.Time `json:"date_approved" bson:"date_approved"`
		newValue, _ = time.Parse(longForm, oldValue.(string))
		// Actions and Notes are missing
	}

	return newValue
}
