package fiddler

import "fmt"

//* Errors used in this package models *//

// When trying to scan records
type scanError struct {
	error error
}

func (s scanError) Error() string {
	return fmt.Sprintf("Error when scanning rows: %s.", s.error)
}

// Error when trying to create a new model
type newmodelError struct {
	tablename string
}

func (n newmodelError) Error() string {
	return fmt.Sprintf("Error when trying to create a new model with %s. ", n.tablename)
}
