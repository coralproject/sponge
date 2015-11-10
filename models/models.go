package models

import (
	"database/sql"
	"fmt"
	"log"
)

// To Do: This needs to be looking at the strategy and model depending on the data that is being pulled

// Model is the interface for all the model structs
type Model interface {
	Print()
	ProcessData(*sql.Rows) ([]Model, error)
}

/* FOR COMMENTS */

// Comment is the struct that will hold the comment row
type Comment struct {
	CommentID                int
	AssetID                  int
	StatusID                 int
	CommentTitle             sql.NullString
	CommentBody              sql.NullString
	UserID                   int
	CreateDate               string
	UpdateDate               sql.NullString
	ApproveDate              sql.NullString
	CommentExcerpt           sql.NullString
	EditorsSelection         sql.NullString
	RecommendationCount      sql.NullInt64
	ReplyCount               int
	IsReply                  string
	CommentSequence          sql.NullString
	UserDisplayName          sql.NullString
	UserURL                  sql.NullString // It was UserReply
	UserTitle                sql.NullString
	UserLocation             sql.NullString
	ShowCommentExcerpt       sql.NullString
	HideRegisteredUserName   sql.NullString
	CommentType              sql.NullString
	ParentID                 sql.NullString
	NotifyViaEmailOnApproval sql.NullInt64
}

// Print only print information about the comment
func (c Comment) Print() {
	fmt.Println("Comment: ", c.CommentID, c.CommentBody)
}

// ProcessData get the data from sd
func (c Comment) ProcessData(sd *sql.Rows) ([]Model, error) {

	var comment Comment
	var comments []Model
	for sd.Next() {
		err := sd.Scan(&comment.CommentID, &comment.AssetID, &comment.StatusID, &comment.CommentTitle, &comment.CommentBody, &comment.UserID, &comment.CreateDate,
			&comment.UpdateDate, &comment.ApproveDate, &comment.CommentExcerpt, &comment.EditorsSelection, &comment.RecommendationCount, &comment.ReplyCount, &comment.IsReply,
			&comment.CommentSequence, &comment.UserDisplayName, &comment.UserURL, &comment.UserTitle, &comment.UserLocation, &comment.ShowCommentExcerpt,
			&comment.HideRegisteredUserName, &comment.CommentType, &comment.ParentID, &comment.NotifyViaEmailOnApproval)
		if err != nil {
			log.Fatal(err)
			return nil, err
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
		comments[n] = comment
	}

	return comments, nil
}

/* FOR ASSETS */

// Asset is the struct that will hold the asset row
type Asset struct {
	AssetID    int
	VendorID   sql.NullString
	SourceID   int
	AssetURL   sql.NullString
	CreateDate string
	UpdateDate sql.NullString
}

// Print only print information about the comment
func (a Asset) Print() {
	fmt.Println("Asset: ", a.AssetID, a.AssetURL)
}

// ProcessData get the data from sd
func (a Asset) ProcessData(sd *sql.Rows) ([]Model, error) {

	var asset Asset
	var assets []Model

	for sd.Next() {
		err := sd.Scan(&asset.AssetID, &asset.VendorID, &asset.SourceID, &asset.AssetURL, &asset.CreateDate, &asset.UpdateDate)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		n := len(assets)
		if len(assets) == cap(assets) {
			// Comments is full and we must expand
			// Double the size and add 1
			newAssets := make([]Model, len(assets), 2*len(assets)+1)
			copy(newAssets, assets)
			assets = newAssets
		}
		assets = assets[0 : n+1]
		assets[n] = asset
	}

	return assets, nil
}

/* NOTES */

// Note is the struct for comment notes
type Note struct {
	CommentNoteID   int
	CommentID       int
	UserID          int
	UserDisplayName string
	CreateDate      string
	UpdateDate      sql.NullString
	CommentNote     string
}

// Print only print information about the comment
func (n Note) Print() {
	fmt.Println("Note: ", n.CommentNoteID, n.CommentNote)
}

// ProcessData get the data from sd
func (n Note) ProcessData(sd *sql.Rows) ([]Model, error) {

	var note Note
	var notes []Model

	for sd.Next() {
		err := sd.Scan(&note.CommentNoteID, &note.CommentID, &note.UserID, &note.UserDisplayName, &note.CreateDate, &note.UpdateDate, &note.CommentNote)
		if err != nil {
			log.Fatal(err)
			return nil, err
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
