/*
Package source implements a way to get data from external sources.

External possible sources:
* MySQL
* API?

*/
package source

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

// Data is a struct that has all the db rows and error field
type Data struct {
	//Rows     *sql.Rows // Move into appropiate structure because this has to work for API too (no database/sql)
	Comments []Comment
	Error    error
}

// Source is where the data is coming from (mysql, api)
type Source interface {
	NewSource() (*Source, error)
	GetNewData() Data
}

// Possible Errors
// - TimeOutReader
// - WrongConnection
// - DataErrorReader
