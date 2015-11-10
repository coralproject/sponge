package utils

import "github.com/coralproject/sponge/models"

// Data is a struct that has all the db rows and error field
type Data struct {
	Error error
	Rows  []models.Model // this could be any of the interfaces defined in models package
	Type  string
}
