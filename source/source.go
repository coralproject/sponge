/*
Package source implements a way to get data from external sources.

External possible sources:
* MySQL
* API

*/
package source

import "github.com/coralproject/sponge/utils"

// Source is where the data is coming from (mysql, api)
type Source interface {
	NewSource() *Source
	GetNewData() utils.Data // Data is a struct that has all the db rows and error field
}
