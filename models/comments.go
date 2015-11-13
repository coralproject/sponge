package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

//* ACTIONS *//

// Action has information on all the actions that a user performs on elements like recommends, likes, comments, notes, share, etc.
type Action struct {
	ID         string    `json:"id" bson:"_id"`
	Type       string    `json:"type" bson:"type"`
	UserID     string    `json:"userid" bson:"userid"`
	Value      string    `json:"value" bson:"value"`
	CreateDate time.Time `json:"createdate" bson:"createdate"`
	UpdateDate time.Time `json:"updatedate" bson:"updatedate"`
}

// Print only print information about the action
func (a Action) Print() {
	fmt.Println("Action: ", a.UserID, a.Type, a.Value)
}

// Transform get the data from sd
func (a Action) Transform(sd *sql.Rows) ([]Model, error) {
	var action Action
	var actions []Model

	for sd.Next() {
		err := sd.Scan(&action.ID, &action.Type, &action.UserID, &action.Value, &action.CreateDate, &action.UpdateDate)
		if err != nil {
			return nil, scanError{error: err}
		}

		n := len(actions)
		if len(actions) == cap(actions) {
			// Comments is full and we must expand
			// Double the size and add 1
			newActions := make([]Model, len(actions), 2*len(actions)+1)
			copy(newActions, actions)
			actions = newActions
		}
		actions = actions[0 : n+1]
		actions[n] = action
	}

	return actions, nil
}

//* NOTES */

// Note is a note on a piece of content
type Note struct {
	ID          string `json:"userId" bson:"userId"`
	CommentID   string
	CommentNote string    `json:"body" bson:"body"`
	CreateDate  time.Time `json:"createdate" bson:"createdate"`
	UpdateDate  time.Time `json:"updatedate" bson:"updatedate"`
}

// Print only print information about the comment
func (n Note) Print() {
	fmt.Println("Note: ", n.ID, n.CommentNote)
}

// Transform get the data from sd
func (n Note) Transform(sd *sql.Rows) ([]Model, error) {
	var note Note
	var notes []Model
	for sd.Next() {
		err := sd.Scan(&note.ID, &note.CommentID, &note.CommentNote, &note.CreateDate, &note.UpdateDate)
		if err != nil {
			return nil, scanError{error: err}
		}

		n := len(notes)
		if len(notes) == cap(notes) {
			// Comments is full and we must expand
			// Double the size and add 1
			newNotes := make([]Model, len(notes), 2*len(notes)+1)
			copy(newNotes, notes)
			notes = newNotes
		}
		notes = notes[0 : n+1]
		notes[n] = note
	}

	return notes, nil
}

//* COMMENTS *//

// Comment is the comment struct
type Comment struct {
	ID          string    `json:"id" bson:"_id"`
	Body        string    `json:"body" bson:"body"`
	ParentID    string    `json:"parentId" bson:"parentId"`
	AssetID     string    `json:"assetId" bson:"assetId"`
	StatusID    string    `json:"status" bson:"status"`
	CreateDate  time.Time `json:"createdDate" bson:"createdDate"`
	UpdateDate  time.Time `json:"updatedDate" bson:"updatedDate"`
	ApproveDate time.Time `json:"approvedDate" bson:"approvedDate"`
	Raw         []string  `json:"raws" bson:"raws"`
	Actions     []Action  `json:"actions" bson:"actions"`
	Notes       []Note    `json:"notes" bson:"notes"`
}

// Print only print information about the comment
func (c Comment) Print() {
	fmt.Println("Comment: ", c.ID, c.Body)
}

// Transform get the data from sd
func (c Comment) Transform(sd *sql.Rows) ([]Model, error) {
	var comment Comment
	var comments []Model
	var raws string

	for sd.Next() {
		err := sd.Scan(&comment.ID, &comment.Body, &comment.ParentID, &comment.AssetID, &comment.StatusID, &comment.CreateDate, &comment.UpdateDate,
			&comment.ApproveDate, &raws)
		if err != nil {
			return nil, scanError{error: err}
		}

		n := len(comments)
		if len(comments) == cap(comments) {
			// Comments is full and we must expand
			// Double the size and add 1
			newComments := make([]Model, len(comments), 2*len(comments)+1)
			copy(newComments, comments)
			comments = newComments
		}
		comments = comments[0 : n+1]
		comment.Raw = strings.Split(raws, ",")
		comments[n] = comment
	}

	return comments, nil
}
