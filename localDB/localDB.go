/*
Package localDB implements a way to get data into our local database.

Possible local databases:
* MongoDB

*/
package localDB

import "database/sql"

////// Structures  //////

// This needs to be looking at the strategy and model depending on the data that is being pulled
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

// Data is where we are pushing the data to
type Data struct {
	Comments []Comment // Look for some kind of structure where to put this data into
	error    error
}

// LocalDB is an interface to the different local databases available
type LocalDB interface {
	NewLocalDB() (*LocalDB, error)
	Add(Data) error
}
