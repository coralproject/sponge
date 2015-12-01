package models

import (
	"database/sql"

	"github.com/coralproject/sponge/config"
)

// To Do: This needs to be looking at the strategy and model depending on the data that is being pulled

// Model is the interface for all the model structs <--- be carefull, this is an interface that Comment, Asset and Note are implementing. Transform is acting on a slice of Model (how that works?)
type Model interface {
	Print()
	Transform(*sql.Rows, config.Table) ([]Model, error)
}

// New creates a new model based on a table name. This is a Factory func
// IT NEEDS TO FIND A WAY TO USE REFLECT/METAPROGRAMMING TO CREATE A "OBJECT" OF TYPE TABLE FOR THE MODELS
func New(table string) (Model, error) {

	if table == "Comment" {
		return Comment{}, nil
	} else if table == "Asset" {
		return Asset{}, nil
	} else if table == "Note" {
		return Note{}, nil
	} else if table == "User" {
		return User{}, nil
	}

	return nil, newmodelError{tablename: table}
}
