package fiddler

import (
	"database/sql"
	"fmt"

	"github.com/coralproject/shelf/pkg/srv/comment"
	configuration "github.com/coralproject/sponge/config"
)

//* COMMENTS *//

// Comment is embedding the comment package to extend it
type Comment struct {
	comment.Comment
}

// // Comment is the comment struct
// type Comment struct {
// 	ID          string         `json:"id" bson:"_id"`
// 	Body        sql.NullString `json:"body" bson:"body"`
// 	ParentID    sql.NullString `json:"parentId" bson:"parentId"`
// 	AssetID     sql.NullString `json:"assetId" bson:"assetId"`
// 	StatusID    sql.NullString `json:"status" bson:"status"`
// 	CreateDate  time.Time      `json:"createdDate" bson:"createdDate"`
// 	UpdateDate  time.Time      `json:"updatedDate" bson:"updatedDate"`
// 	ApproveDate time.Time      `json:"approvedDate" bson:"approvedDate"`
// 	Raw         []string       `json:"raws" bson:"raws"`
// 	Actions     []Action       `json:"actions" bson:"actions"`
// 	Notes       []Note         `json:"notes" bson:"notes"`
// }

// Print only print information about the comment
func (c Comment) Print() {
	fmt.Println("Comment: ", c.ID, c.Body)
}

// Transform get the data from sd
func (c Comment) Transform(sd *sql.Rows, table configuration.Table) ([]Transformer, error) {
	var comment Comment
	var comments []Transformer

	var createDate sql.NullString
	var updateDate sql.NullString // convert from sql.NullString to string
	var approveDate sql.NullString

	//var raw string
	// var title, userid, editorSelection, recomendationCount, replyCount, isReply, userDisplayName sql.NullString
	// var userURL, userTitle, userLocation, showCommentExcerpt, hideRegisteredUserName, commentType sql.NullString
	// var commentExcerpt, commentSequence, notifyViaEmailOnApproval sql.NullString

	for sd.Next() {
		// To Do: It needs to be able to get the extra fields into raw (we don't really know in the code which fields the comment table has)
		err := sd.Scan(&comment.ID, &comment.Body, &comment.ParentID, &comment.AssetID, &comment.Status,
			&createDate, &updateDate, &approveDate)
		if err != nil {
			return nil, scanError{error: err}
		}

		//comment.CreateDate, _ = time.Parse("2006-01-02", createDate) // To Do: I need to see how to dinamically discover what is the dateTime layout
		//comment.UpdateDate, _ = time.Parse("2006-01-02", updateDate)
		//comment.ApproveDate, _ = time.Parse("2006-01-02", approveDate)

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

	sd.Close()

	return comments, nil
}
