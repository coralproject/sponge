# Coral

import "github.com/coralproject/sponge/pkg/coral"

Interact with Pillar endpoints to import data into the Coral system.

## constants

``const (
  	retryTimes int    = 3
  	methodGet  string = "GET"
  	methodPost string = "POST"
  )``

Retrytimes determines how many times to retry if it fails.


## variables

``var (
	endpoints map[string]string // model -> endpoint
)``

Endpoints have all the services where to send data. Right now that is Pillar.

  ``var (
  	uuid string
  	str  strategy.Strategy
  )``

UUID is the universal identifier we are using for logging.
str is the strategy configuration.


## func AddRow

  ``func AddRow(data map[string]interface{}, tableName string) error``

Adds data to the collection "tableName".

## func CreateIndex

  ``func CreateIndex(collection string) error``

Calls the service to create index for collectionn.
