/*
Package config handles the loading and distribution of configuration related with external sources.
*/
package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/coralproject/core/log"
)

//* CONFIGURATION STRUCTS *//

// Table holds the struct on what is the external source's table name and fields
type Table struct {
	Name   string            `json:"name"`
	Fields map[string]string `json:"fields"`
}

// Strategy explains which tables or data we are getting from the source.
type Strategy struct {
	Typesource string           `json:"typesource"`
	Tables     map[string]Table `json:"tables"`
	Action     []Table
}

// Credential is the interface for APIs or Database Sources
type Credential interface {
	GetAdapter() string
}

// CredentialDatabase has the information to connect to the external databaase source.
type CredentialDatabase struct {
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`    //= '5432'
	Adapter  string `json:"adapter"` //= 'mysql'
	Type     string `json:"type"`    //= 'source' or 'local'
}

// GetAdapter returns the adapter
func (c CredentialDatabase) GetAdapter() string {
	return c.Adapter
}

// CredentialAPI has the information to connect to an external API source.
type CredentialAPI struct {
	Username  string            `json:"username"` // BasicAuth
	Password  string            `json:"password"` // BasicAuth
	Adapter   string            `json:"adapter"`
	Endpoints map[string]string `json:"endpoints"`
}

// GetAdapter returns the adapter
func (c CredentialAPI) GetAdapter() string {
	return c.Adapter
}

// GetEndpoints returns all the endpoints
func (c CredentialAPI) GetEndpoints() map[string]string {
	return c.Endpoints
}

// GetEndpoint gives the endpoint for that modelName
func (c CredentialAPI) GetEndpoint(modelName string) (string, error) {
	endpoints := c.GetEndpoints()
	for k, e := range endpoints {
		if k == modelName {
			return e, nil
		}
	}

	return "", endpointError{key: modelName}
}

// GetAuthenticationEndpoint returns the authentication url
func (c CredentialAPI) GetAuthenticationEndpoint() (string, error) {
	for key, val := range c.Endpoints {
		if key == "authentication" {
			return val, nil
		}
	}

	return "", endpointError{key: "authentication"}
}

// Credentials are all the credentials for external and internal data sources
type Credentials struct {
	Databases []CredentialDatabase
	APIs      []CredentialAPI
}

// Config is a structure with all the information for the specitic strategy (how to get the data, from which source)
type Config struct {
	Name        string
	Strategy    Strategy
	Credentials Credentials // map[string][]Credential // String is "Databases" or "APIs" indicating which kind of credentials are those
}

/* Exported Functions */

// New creates a new config
func New() *Config {

	f := "config/config.json"
	config, err := readConfigFile(f)

	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	return config
}

// GetCredential returns the credentials for connection with the external source
func (conf Config) GetCredential(adapter string) Credential {
	var cred Credential

	creds := conf.Credentials.APIs

	// look at the credentials related to local database (mongodb in our original example)
	for i := 0; i < len(creds); i++ {
		if creds[i].GetAdapter() == adapter {
			cred = creds[i]
			return cred
		}
	}

	creda := conf.Credentials.Databases

	// look at the credentials related to local database (mongodb in our original example)
	for i := 0; i < len(creda); i++ {
		if creda[i].GetAdapter() == adapter {
			cred = creda[i]
			return cred
		}
	}
	log.Fatal(getCredentialError{adapter: adapter}.Error())

	return cred
}

// GetStrategy returns the strategy
// To Do: Needs to manage errors
func (conf Config) GetStrategy() Strategy {
	strategy := conf.Strategy

	// To Do: catch the error on getting credentials and return it
	return strategy
}

// GetTables returns a list of tables to be imported
func (conf Config) GetTables() map[string]Table {
	// To Do: catch the error when no Tables
	return conf.Strategy.Tables
}

// GetTableName returns the external source's table mapped to the coral model
func (conf Config) GetTableName(modelName string) string {
	return conf.Strategy.Tables[modelName].Name
}

// GetTableFields returns the external source's table fields mapped to the coral model
func (conf Config) GetTableFields(modelName string) map[string]string {
	return conf.Strategy.Tables[modelName].Fields
}

/* Not Exported Functions */

// Read the configuration file and load it into the Config
func unmarshal(content []byte) (*Config, error) {

	c := new(Config)

	err := json.Unmarshal(content, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func readConfigFile(f string) (*Config, error) {

	content, err := ioutil.ReadFile(f)
	if err != nil {
		e := readingFileError{filename: f, err: err}
		log.Fatal(e.Error())

		return nil, e
	}

	config, err := unmarshal(content)
	if err != nil {
		e := parseFileError{filename: f, err: err}
		log.Fatal(e.Error())

		return nil, e
	}

	return config, err
}
