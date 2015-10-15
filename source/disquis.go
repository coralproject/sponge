/*
Package source implements a way to get data from external sources.

External possible sources:
* MySQL
* API?

*/
package source

/* Implementing the Sources */

//////////// Disquis API SOURCE ////////////

// DisquisAPI is the struct that has the connection to the Disquis API
type DisquisAPI struct {
	connection string
}

// Reader on the API brings back the data needed from the API
func (a DisquisAPI) Reader() Data {
	var d Data
	return d
}

// Possible Errors
// - TimeOutReader
// - WrongConnection
// - DataErrorReader
