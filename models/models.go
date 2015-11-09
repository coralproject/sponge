package models

import (
	"database/sql"
	"fmt"
	"log"
)

// This needs to be looking at the strategy and model depending on the data that is being pulled

// Comment is the struct that will hold the comment row
type Comment struct {
	CommentID              int
	AssetID                int
	StatusID               int
	CommentTitle           sql.NullString
	CommentBody            sql.NullString
	UserID                 int
	CreateDate             string
	UpdateDate             sql.NullString
	ApproveDate            sql.NullString
	CommentExcerpt         sql.NullString
	EditorsSelection       sql.NullString
	RecommendationCount    sql.NullInt64
	ReplyCount             int
	IsReply                string
	CommentSequence        sql.NullString
	UserDisplayName        sql.NullString
	UserReply              sql.NullString
	UserTitle              sql.NullString
	UserLocation           sql.NullString
	ShowCommentExcerpt     sql.NullString
	HideRegisteredUserName sql.NullString
	CommentType            sql.NullString
	ParentID               sql.NullString
}

// Model is the interface for all the model structs
type Model interface {
	Print()
	ProcessData(*sql.Rows) ([]Model, error)
}

/* FOR COMMENTS */

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
			&comment.CommentSequence, &comment.UserDisplayName, &comment.UserReply, &comment.UserTitle, &comment.UserLocation, &comment.ShowCommentExcerpt,
			&comment.HideRegisteredUserName, &comment.CommentType, &comment.ParentID)
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
