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
// IT NEEDS TO FIND A WAY TO USE REFLECT/METAPROGRAMMING TO CREATE A "OBJECT" OF TYPE TABLE FOR THE MODELS
func New(table string) (models.Model, error) {

	var err error

	if table == "Comment" {
		return models.Comment{}, err
	} else if table == "Asset" {
		return models.Asset{}, err
	} else if table == "Note" {
		return models.Note{}, err
	}

	err = fmt.Errorf("Error when trying to create a new model with %s. ", table)
	return nil, err
}
