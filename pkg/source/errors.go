package source

import "fmt"

//* Errors used in this package *//

// When trying to connect to the database with the connection string
type connectError struct {
	connection string
}

func (e connectError) Error() string {
	return fmt.Sprintf("Error when connecting to database with %s.", e.connection)
}

// When trying to query the database with the query string
type queryError struct {
	query string
}

func (e queryError) Error() string {
	return fmt.Sprintf("Error when quering the database with %s.", e.query)
}

// When trying to create a new model... <-- To Do
type modelError struct {
	model string
}

func (e modelError) Error() string {
	return fmt.Sprintf("Error when trying to create a new model %s.", e.model)
}

// When trying to find the connection to dbms
type notFoundError struct {
	dbms string
}

func (e notFoundError) Error() string {
	return fmt.Sprintf("Error when trying to get the dbms %s.", e.dbms)
}
