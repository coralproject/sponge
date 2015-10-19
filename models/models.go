package models

import "database/sql"

// This needs to be looking at the strategy and model depending on the data that is being pulled

// Comment is the struct that will hold the comment row
type Comment struct {
	CommentID              int
	AssetID                int
	StatusID               int
	CommentTitle           string
	CommentBody            string
	UserID                 int
	CreateDate             string
	UpdateDate             string
	ApproveDate            sql.NullString
	CommentExcerpt         sql.NullString
	EditorsSelection       string
	RecommendationCount    int
	ReplyCount             int
	IsReply                string
	CommentSequence        sql.NullString
	UserDisplayName        string
	UserReply              string
	UserTitle              string
	UserLocation           string
	ShowCommentExcerpt     string
	HideRegisteredUserName string
	CommentType            sql.NullString
	ParentID               sql.NullString
}
