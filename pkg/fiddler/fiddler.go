/*
Package fiddler transform, through a strategy file, data from external source into our local coral schema.

*/
package fiddler

import (
	"github.com/coralproject/sponge/pkg/log"
	str "github.com/coralproject/sponge/strategy"
)

// global variables related to strategy
var strategy = str.New() // Reads the strategy file

const longForm = "2015-11-02 12:26:05" // date format. To Do: it needs to be defined in the strategy file for the publisher

// To Do: This needs to be looking at the strategy and model depending on the data that is being pulled

// Transformer is the interface for all the model structs <--- be carefull, this is an interface that Comment, Asset and Note are implementing. Transform is acting on a slice of Model (how that works?)
type Transformer interface {
	Print()
	Transform([]map[string]interface{}, str.Table) ([]Transformer, error)
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

	e := newmodelError{tablename: table}
	log.Error("transform", "New", e, "Factory Model")

	return nil, e
}

// Transform from external source data into the coral schema
func Transform(modelName string, data []map[string]interface{}) ([]Transformer, error) { //data *sql.Rows) ([]Transformer, error) {
	var dataCoral []Transformer

	m, err := New(modelName)
	if err != nil {
		log.Error("transform", "Transform", err, "Transform factory")
		return nil, err
	}

	dataCoral, err = m.Transform(data, strategy.GetTables()[modelName])
	if err != nil {
		log.Error("transform", "Transform", err, "Transform Factory")
		return nil, err
	}

	// Get dataCoral into JSON

	return dataCoral, nil
}
