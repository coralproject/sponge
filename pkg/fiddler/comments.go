package fiddler

import (
	"fmt"
	"log"
	"time"

	"github.com/coralproject/shelf/pkg/srv/comment"
	configuration "github.com/coralproject/sponge/config"
	"github.com/oleiade/reflections"
)

//* COMMENTS *//

// Comment is embedding the comment package to extend it
type Comment struct {
	comment.Comment
}

// Print only print information about the comment
func (c Comment) Print() {
	fmt.Println("Comment: ", c.CommentID, c.Body)
}

// Transform get the data from sd
func (c Comment) Transform(sd []map[string]interface{}, table configuration.Table) ([]Transformer, error) {

	var comment Comment
	var comments []Transformer

	// To Do: it needs refactoring as my gut tells me that is quite inefficient
	for _, value := range sd {
		for coralField, f := range table.Fields {

			// convert field f with value value[f] into field coralField
			newValue := transformCommentField(f, value[f], coralField)

			err := reflections.SetField(&comment, coralField, newValue)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}
		}

		n := len(comments)
		if len(comments) == cap(comments) {
			// Comments is full and we must expand
			// Double the size and add 1
			newComments := make([]Transformer, len(comments), 2*len(comments)+1)
			copy(newComments, comments)
			comments = newComments
		}
		comments = comments[0 : n+1]
		//comment.Raw = strings.Split(raws, ",")
		comments[n] = comment
	}

	return comments, nil
}

//Here we transform the record into what we want (based on the configuration)
// 1. convert types (values are all strings) into the struct
func transformCommentField(sourceField string, oldValue interface{}, coralField string) interface{} {

	var newValue interface{}

	// Right now this the simpler thing to do. Needs to look more into reflect to do it
	// dinamically (if it is worth it) or merge this into the comment package
	switch coralField {
	case "CommentID": //string
		newValue = oldValue
	case "ParentID": //string    `json:"parent_id" bson:"parent_d"`
		newValue = oldValue
	case "AssetID": //string    `json:"asset_id" bson:"asset_id"`
		newValue = oldValue
	case "Path": //string    `json:"path" bson:"path"`
		newValue = oldValue
	case "Body": //string    `json:"body" bson:"body"`
		newValue = oldValue
	case "Status": //string    `json:"status" bson:"status"`
		newValue = oldValue
	case "DateApproved": //time.Time `json:"date_approved" bson:"date_approved"`
		newValue, _ = time.Parse(longForm, oldValue.(string))
	case "DateModified": //time.Time `json:"date_modified" bson:"date_modified"`
		newValue, _ = time.Parse(longForm, oldValue.(string))
	case "DateCreated": //  time.Time `json:"date_created" bson:"date_created"`
		newValue, _ = time.Parse(longForm, oldValue.(string))
		// Actions and Notes are missing
	}

	return newValue
}
