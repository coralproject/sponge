package fiddler

import (
	"database/sql"
	"fmt"

	"github.com/coralproject/shelf/pkg/srv/comment"
	configuration "github.com/coralproject/sponge/config"
)

//User is embedding the comment package to extend it
type User struct {
	comment.User
}

// // User has information on the user
// type User struct {
// 	ID          string   `json:"id" bson:"_id"`
// 	DisplayName string   `json:"displayName" bson:"displayName"`
// 	Name        string   `json:"name" bson:"name"`
// 	Email       string   `json:"email" bson:"email"`
// 	Raw         []string `json:"raws" bson:"raws"`
// }

// Print only print information about the user
func (u User) Print() {
	fmt.Println("User: ", u.ID, u.UserName)
}

// Transform get the data from sd
func (u User) Transform(sd *sql.Rows, table configuration.Table) ([]Transformer, error) {
	var user User
	var users []Transformer
	// var id, cassetID, statusID, title,
	// 	body, createDate, updateDate,
	// 	approveDate, commentExcerpt, editorSelection, recomendationCount,
	// 	replyCount, isReply, commentSequence, userURL, userTitle,
	// 	userLocation, showCommentExcerpt, hideRegisteredUserName, commentType,
	// 	parentID, notifyViaEmailOnApproval sql.NullString

	for sd.Next() {
		err := sd.Scan(&user.ID, &user.UserName)
		if err != nil {
			return nil, scanError{error: err}
		}

		n := len(users)
		if len(users) == cap(users) {
			// Comments is full and we must expand
			// Double the size and add 1
			newUsers := make([]Transformer, len(users), 2*len(users)+1)
			copy(newUsers, users)
			users = newUsers
		}
		users = users[0 : n+1]
		//user.Raw = strings.Split(raws, ",")
		users[n] = user
	}

	return users, nil
}
