/*
Package fiddler transform, through a strategy file, data from external source into our local coral schema.

*/
package fiddler

import (
	"database/sql"
	"log"

	configuration "github.com/coralproject/sponge/config"
)

// global variables related to configuration
var config = *configuration.New() // Reads the configuration file

// To Do: This needs to be looking at the strategy and model depending on the data that is being pulled

// Transformer is the interface for all the model structs <--- be carefull, this is an interface that Comment, Asset and Note are implementing. Transform is acting on a slice of Model (how that works?)
type Transformer interface {
	Print()
	Transform(*sql.Rows, configuration.Table) ([]Transformer, error)
}

// New creates a new model based on a table name. This is a Factory func
// IT NEEDS TO FIND A WAY TO USE REFLECT/METAPROGRAMMING TO CREATE A "OBJECT" OF TYPE TABLE FOR THE MODELS
func New(table string) (Transformer, error) {

	if table == "Comment" {
		return Comment{}, nil
	} else if table == "Asset" {
		return Asset{}, nil
	} else if table == "User" {
		return User{}, nil
	}

	return nil, newmodelError{tablename: table}
}

// Transform from external source data into the coral schema
func Transform(modelName string, data *sql.Rows) ([]Transformer, error) {
	var dataCoral []Transformer

	m, err := New(modelName)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	dataCoral, err = m.Transform(data, config.GetTables()[modelName])
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	return dataCoral, nil
}
