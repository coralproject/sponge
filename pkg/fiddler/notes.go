package fiddler

import (
	"fmt"
	"time"

	str "github.com/coralproject/sponge/strategy"
)

//* NOTES */

// Note is embedding the comment package to extend it
type Note struct {
	fields map[string]interface{}
}

// Print only print information about the comment
func (n Note) Print() {
	fmt.Println("Note: ", n.fields["UserID"], n.fields["Body"])
}

// Transform get the data from sd
func (n Note) Transform(sd []map[string]interface{}, table str.Table) ([]Transformer, error) {

	var note Note
	note.fields = make(map[string]interface{})
	var notes []Transformer

	// To Do: it needs refactoring as my gut tells me that is quite inefficient
	for _, value := range sd {
		for coralField, f := range table.Fields {

			// convert field f with value value[f] into field coralField
			newValue := transformNoteField(f, value[f], coralField)
			if newValue != nil {
				note.fields[coralField] = newValue
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
