/*
Package main

Import source database into mongodb

*/
package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/coralproject/mod-data-import/source"
)

// Comment is the struct that will hold the comment row
type Comment struct {
	commentID              int
	assetID                int
	statusID               int
	commentTitle           string
	commentBody            string
	userID                 int
	createDate             string
	updateDate             string
	approveDate            sql.NullString
	commentExcerpt         sql.NullString
	editorsSelection       string
	recommendationCount    int
	replyCount             int
	isReply                string
	commentSequence        sql.NullString
	userDisplayName        string
	userReply              string
	userTitle              string
	userLocation           string
	showCommentExcerpt     string
	hideRegisteredUserName string
	commentType            sql.NullString
	parentID               sql.NullString
}

func main() {
	// Connects into mysql database and retrieve one row
	var m *source.MySQL
	var err error

	m, err = source.NewSource()
	if err != nil {
		log.Fatal("Error when creating new source ", err)
	}

	db, err := m.Open()
	if err != nil {
		log.Fatal("Error when connection to database. ", err)
	}
	defer m.Close(db)

	var comment Comment
	stmt, err := db.Prepare("select * from nyt_comments where commentID = ?")
	if err != nil {
		log.Fatal(err)
	}
	err = stmt.QueryRow(13471611).Scan(&comment.commentID, &comment.assetID, &comment.statusID, &comment.commentTitle, &comment.commentBody, &comment.userID, &comment.createDate, &comment.updateDate, &comment.approveDate, &comment.commentExcerpt, &comment.editorsSelection, &comment.recommendationCount, &comment.replyCount, &comment.isReply, &comment.commentSequence, &comment.userDisplayName, &comment.userReply, &comment.userTitle, &comment.userLocation, &comment.showCommentExcerpt, &comment.hideRegisteredUserName, &comment.commentType, &comment.parentID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", comment)

	// Get that comment into a MongoDB collection called "Comments"
}
