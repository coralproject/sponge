package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/coralproject/sponge/config"
)

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
func (n Note) Transform(sd *sql.Rows, table config.Table) ([]Model, error) {
	var note Note
	var notes []Model

	var createDate string
	var updateDate string
	var userID, userDisplayName string

	for sd.Next() {
		err := sd.Scan(&note.ID, &note.CommentID, &userID, &userDisplayName, &createDate, &updateDate, &note.CommentNote)
		if err != nil {
			return nil, scanError{error: err}
		}

		note.CreateDate, _ = time.Parse("2006-01-02", createDate) // To Do: I need to see how to dinamically discover what is the dateTime layout
		note.UpdateDate, _ = time.Parse("2006-01-02", updateDate)

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
