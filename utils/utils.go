package utils

import "github.com/coralproject/sponge/models"

// Data is a struct that has all the db rows and error field
type Data struct {
	//Rows     *sql.Rows // Move into appropiate structure because this has to work for API too (no database/sql)
	Comments []models.Comment
	Error    error
}
