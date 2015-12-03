package fiddler

import (
	"fmt"

	"github.com/coralproject/shelf/pkg/srv/comment"
	configuration "github.com/coralproject/sponge/config"
	"github.com/goinggo/mapstructure"
)

//* COMMENTS *//

// Comment is embedding the comment package to extend it
type Comment struct {
	comment.Comment
}

// STRUCTS //
// ID           string    `json:"id" bson:"_id"`
// ParentID     string    `json:"parent_id" bson:"parent_d"`
// AssetID      string    `json:"asset_id" bson:"asset_id"`
// Path         string    `json:"path" bson:"path"`
// Body         string    `json:"body" bson:"body"`
// Status       string    `json:"status" bson:"status"`
// DateCreated  time.Time `json:"date_created" bson:"date_created"`
// DateUpdated  time.Time `json:"date_updated" bson:"date_updated"`
// DateApproved time.Time `json:"date_approved" bson:"date_approved"`
// Actions      []Action  `json:"actions" bson:"actions"`
// Notes        []Note    `json:"notes" bson:"notes"`

// type Note struct {
// 	UserID string    `json:"userId" bson:"userId"`
// 	Body   string    `json:"body" bson:"body"`
// 	Date   time.Time `json:"date" bson:"date"`
// }

// Configuration //
// "Comment": {
// 	"name": "crnr_comment",
// 	"fields": {
// 		"ID": "CommentID",
// 		"Body": "CommentBody",
// 		"ParentID": "ParentID",
// 		"AssetID": "AssetID",
// 		"StatusID": "StatusID",
// 		"CreateDate": "CreateDate",
// 		"UpdateDate": "UpdateDate",
// 		"ApproveDate": "ApproveDate"
// 	}
// },

// Print only print information about the comment
func (c Comment) Print() {
	fmt.Println("Comment: ", c.ID, c.Body)
}

// Transform get the data from sd
func (c Comment) Transform(sd []map[string]interface{}, table configuration.Table) ([]Transformer, error) {

	// sd is in the type
	// [map[commentbody:And we Southerners welcome seeing these people leaving! commentid:16571170 createdate:2015-11-04 14:26:31 parentid:16565164
	// statusid:1 updatedate:2015-11-04 14:26:31 approvedate: assetid:3441167] map[createdate:2015-11-04 14:26:38 parentid: statusid:1
	// updatedate:2015-11-04 14:26:38 approvedate: assetid:3441167
	// commentbody:The "right" wing and the church are doing everything they can to promote prejudice and intolerance. The Trans community makes up 0.2% of the population. I bet these people have never even seen let alone met a Trans person yet...they are coming for you and you little dog too!!!<br/>Put up or shut up NFL. Find a new host city.<br/>This is about denying equal rights to gays and lesbians. It has little to do Trans people or bathrooms. It's a red herring.
	// commentid:16571171]]

	var comment *Comment
	var comments []Transformer

	// To Do: it needs refactoring as my gut tells me that is quite inefficient
	for _, value := range sd {
		for coral_field, f := range table.Fields {
		// Really transform the record into what we want
		auxcomment[coral_field] := value[f]
	}
		mapstructure.DecodePath(auxcomment, &comment)
		fmt.Println(value)


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

	// var createDate sql.NullString
	// var updateDate sql.NullString // convert from sql.NullString to string
	// var approveDate sql.NullString

	//var raw string
	// var title, userid, editorSelection, recomendationCount, replyCount, isReply, userDisplayName sql.NullString
	// var userURL, userTitle, userLocation, showCommentExcerpt, hideRegisteredUserName, commentType sql.NullString
	// var commentExcerpt, commentSequence, notifyViaEmailOnApproval sql.NullString

	// for sd.Next() {
	// 	// To Do: It needs to be able to get the extra fields into raw (we don't really know in the code which fields the comment table has)
	// 	err := sd.Scan(&comment.ID, &comment.Body, &comment.ParentID, &comment.AssetID, &comment.Status,
	// 		&createDate, &updateDate, &approveDate)
	// 	if err != nil {
	// 		return nil, scanError{error: err}
	// 	}
	//
	// 	//comment.CreateDate, _ = time.Parse("2006-01-02", createDate) // To Do: I need to see how to dinamically discover what is the dateTime layout
	// 	//comment.UpdateDate, _ = time.Parse("2006-01-02", updateDate)
	// 	//comment.ApproveDate, _ = time.Parse("2006-01-02", approveDate)
	//
	// 	n := len(comments)
	// 	if len(comments) == cap(comments) {
	// 		// Comments is full and we must expand
	// 		// Double the size and add 1
	// 		newComments := make([]Transformer, len(comments), 2*len(comments)+1)
	// 		copy(newComments, comments)
	// 		comments = newComments
	// 	}
	// 	comments = comments[0 : n+1]
	// 	//comment.Raw = strings.Split(raws, ",")
	// 	comments[n] = comment
	// }
	//
	// sd.Close()

	return comments, nil
}
