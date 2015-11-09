package utils

import (
	"fmt"

	"github.com/coralproject/sponge/models"
)

// Data is a struct that has all the db rows and error field
type Data struct {
	Error error
	Rows  []models.Model // this could be any of the interfaces defined in models package
	Type  string
}

// New creates a new model based on a table name
func New(table string) models.Model {
	fmt.Println("POOO")
	return models.Comment{}
}
