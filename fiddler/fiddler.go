/*
Package fiddler transform, through a strategy file, data from external source into our local coral schema.

*/
package fiddler

import (
	"database/sql"
	"log"

	configuration "github.com/coralproject/sponge/config"
	"github.com/coralproject/sponge/models"
)

// global variables related to configuration
var config = *configuration.New() // Reads the configuration file

// Transform from external source data into the coral schema
func Transform(modelName string, data *sql.Rows) ([]models.Model, error) {
	var dataCoral []models.Model

	m, err := models.New(modelName)
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
