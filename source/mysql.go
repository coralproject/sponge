/*
Package source implements a way to get data from external MySQL sources.

External possible sources:
* MySQL
* API

*/
package source

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	configuration "github.com/coralproject/sponge/config"

	"github.com/coralproject/sponge/models"
	"github.com/coralproject/sponge/utils"

	_ "github.com/go-sql-driver/mysql" // Check if this can be imported not blank. To Do.
)

// Global configuration variables that holds all the configuration from the config file
var config = *configuration.New() // Reads the configuration file
var credential = config.GetCredential("mysql")

/* Implementing the Sources */

// MySQL is the struct that has the connection string to the external mysql database
type MySQL struct {
	Connection string
	Database   *sql.DB
}

/* Exported Functions */

// NewSource returns a new Mysql struct with the connection string in it
// Method required by source.Interface
func NewSource() *MySQL {
	//connection := mysqlConnection(config)

	// Get MySQL connection string
	return &MySQL{Connection: connection(), Database: nil}
}

// // GetNewData returns the data requested
// // Method required by source.Interface
// func (m *MySQL) GetNewData() utils.Data {
// 	var d utils.Data
// 	return d
// }

// GetNewData returns the data requested for the table
// Method required by source.Interface
// To Do: we need to get the tables from the configuration file and get data per table with limit
func (m *MySQL) GetNewData() utils.Data {

	table := config.GetTables()["comments"] // To Do: it needs to get all the tables and get data for each one. Maybe concurrent for each table?

	// We are sending the data in the utils.Data
	var d utils.Data

	db, err := m.open()
	if err != nil {
		log.Fatal("Error when connection to database with ", m.Connection, err)
		d.Error = err
	}

	m.Database = db

	defer m.close(db)

	queryString := strings.Join([]string{"SELECT * from", table, "limit 10"}, " ")
	//sd, err := m.Database.Query("SELECT commentID, assetID, statusID, commentTitle, commentBody, userID, createDate, updateDate, approveDate, commentExcerpt, editorsSelection, recommendationCount, replyCount, isReply, commentSequence, userDisplayName, userReply, userTitle, userLocation, showCommentExcerpt, hideRegisteredUserName, commentType, parentID from nyt_comments")
	sd, err := m.Database.Query(queryString)

	if err != nil {
		log.Fatal("Error when quering the DB ", err)
		d.Error = err
	}

	var comment models.Comment

	for sd.Next() {
		if sd.Scan(&comment.CommentID, &comment.AssetID, &comment.StatusID, &comment.CommentTitle, &comment.CommentBody, &comment.UserID, &comment.CreateDate, &comment.UpdateDate, &comment.ApproveDate, &comment.CommentExcerpt, &comment.EditorsSelection, &comment.RecommendationCount, &comment.ReplyCount, &comment.IsReply, comment.CommentSequence, &comment.UserDisplayName, &comment.UserReply, &comment.UserTitle, &comment.UserLocation, comment.ShowCommentExcerpt, &comment.HideRegisteredUserName, &comment.CommentType, &comment.ParentID); err != nil {
			log.Fatal(err)
		}
		n := len(d.Comments)
		if len(d.Comments) == cap(d.Comments) {
			// Comments is full and we must expand
			// Double the size and add 1
			newComments := make([]models.Comment, len(d.Comments), 2*len(d.Comments)+1)
			copy(newComments, d.Comments)
			d.Comments = newComments
		}

		d.Comments = d.Comments[0 : n+1]
		d.Comments[n] = comment
	}
	d.Error = nil

	return d
}

/* Not exported functions */

// // Returns the connection string
// func mysqlConnection(conf configuration.Config) string {
//
// 	cred := conf.GetCredential("source")
//
// 	connection := cred.Username + ":" + cred.Password + "@" + "/" + cred.Database
//
// 	return connection
// }

func connection() string {
	fmt.Println(credential)
	fmt.Println("credential")
	return credential.Username + ":" + credential.Password + "@" + "/" + credential.Database
}

// Open gives back a pointer to the DB
func (m *MySQL) open() (*sql.DB, error) {

	var err error
	m.Database, err = sql.Open("mysql", m.Connection)
	if err != nil {
		log.Fatal("Could not connect to MySQL database with ", m.Connection, err)
		return nil, err
	}

	err = m.Database.Ping()
	if err != nil {
		log.Fatal("Could not connect to the database with ", m.Connection, err)
		return nil, err
	}

	return m.Database, nil
}

// Close closes the db
func (m MySQL) close(db *sql.DB) error {
	return db.Close()
}

// Get returns data from the query to the db
func (m MySQL) get(db *sql.DB, query string) *sql.Rows {

	// LOOK INTO config.Strategy to see which is the strategy to follow
	d, err := db.Query(query)
	if err != nil {
		log.Fatal("Error when quering the DB ", err)
	}

	// To Do: it needs to return DATA type
	return d
}
