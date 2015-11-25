package config

//* Errors used in this package *//
import "fmt"

// When trying to connect to the database with the connection string
type endpointError struct {
	key string
}

func (e endpointError) Error() string {
	return fmt.Sprintf("Error when trying to get endpoint %s.", e.key)
}

// When reading the configuration file
type readingFileError struct {
	filename string
	err      error
}

func (e readingFileError) Error() string {
	return fmt.Sprintf("Error when getting the configuration file %s. Error: %s", e.filename, e.err)
}

// When trying to get the credential for adapter
type getCredentialError struct {
	adapter string
}

func (e getCredentialError) Error() string {
	return fmt.Sprintf("Error when trying to get the credential for %s.", e.adapter)
}

// When trying to parse the configuration file
type parseFileError struct {
	filename string
	err      error
}

func (e parseFileError) Error() string {
	return fmt.Sprintf("Unable to parse JSON in config file %s. Error: %s", e.filename, e.err)
}
