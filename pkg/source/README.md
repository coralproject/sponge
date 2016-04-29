# Source

import "github.com/coralproject/sponge/pkg/source"

Package source has the drivers to connect to the external source and retrieve data. The credentials for the external source are setup in the strategy file.


## Variables

	var strategy str.Strategy

Holds the credentials for the external source as well as all the entities that needs to be extracted.

	var uuid     string

Universally Unique Identifier used for the logs.

	var credential str.Credential

Credential for the external source (database or web service).

## Interface Sourcer

This is the interface that needs to be implemented by any driver to databases that wants to be included.

### func GetData

	`GetData(string, *Options) ([]map[string]interface{}, error)`

Returns all the data (query by options in Options) in the format []map[string]interface{}

### func IsWebService

	`IsWebService() bool`

Returns true if the imlementation of sourcer is a web service.


## func New

	`func New(d string) (Sourcer, error)`

Depending on the parameter it returns a structure with the connection to the external source that implements the interface Sourcer.


## func GetEntities

`func GetEntities() ([]string, error)`

Gets all the entities names from the source

## func GetforeignEntity

	`func GetForeignEntity(name string) string`

Gets the foreign entinty's name for the coral collection.


## Driver mysql

Structure mysql that implements interface Sourcer. It gets data from a Mysql database.

## Driver postgresql

Structure postgresql that implements interface Sourcer. It gets data from a Postgresql database.

## Driver mongodb

Structure to hold a connection to a mongoDB that implements interface Sourcer. It gets data from the mongo database.

## API

Structure the way to get data from a web service that implements interface Sourcer.

## How to add a new source

If you need to add a new external source implements the interface Sourcer over an structure with the connection to the database. 
